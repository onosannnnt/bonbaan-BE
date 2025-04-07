package recommendationAdepter

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	recommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"gorm.io/gorm"
)

type recommendationDriver struct {
	db *gorm.DB
}
// Add this method to the ServiceDriver struct
func (d *recommendationDriver) InitializeFullTextSearchIndex() error {
    // First check if the Thai text search configuration exists
    var thaiConfigExists bool
    err := d.db.Raw("SELECT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'thai')").Scan(&thaiConfigExists).Error
    if err != nil {
        return err
    }
    
    // Create Thai text search configuration if it doesn't exist
    if !thaiConfigExists {
        // Create Thai configuration based on simple
        err = d.db.Exec("CREATE TEXT SEARCH CONFIGURATION thai (COPY = simple)").Error
        if err != nil {
            return err
        }
        
        // Alter the mapping to use simple dictionary for word type
        err = d.db.Exec("ALTER TEXT SEARCH CONFIGURATION thai ALTER MAPPING FOR word WITH simple").Error
        if err != nil {
            return err
        }
    }
    
    // Check if the index already exists
    var indexExists bool
    err = d.db.Raw("SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_services_fts')").Scan(&indexExists).Error
    if err != nil {
        return err
    }

    // Create the index using the Thai configuration if it doesn't exist
    if !indexExists {
        return d.db.Exec("CREATE INDEX idx_services_fts ON services USING gin(to_tsvector('thai', name || ' ' || description || ' ' || address))").Error
    }
	//Ensure pg_trgm is enable
	err = d.db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error
	if err != nil {
		return err
	}

    
    return nil
}
// NewRecommendationDriver initializes the repository.
// Note: The cache is now provided by the utils package.
func NewRecommendationDriver(db *gorm.DB) recommendationUsecase.RecommendationRepository {
    driver := &recommendationDriver{
        db: db,
    }
    
    // Initialize the full-text search index (ignore error for simplicity)
    _ = driver.InitializeFullTextSearchIndex()
    
    return driver
}

// Insert adds or updates a recommendation and its associated util record.
// Consider invalidating related cache entries here if the underlying data changes.
func (d *recommendationDriver) Insert(recParam *Entities.Recommendation) error {
	var rec Entities.Recommendation
	err := d.db.
		Where("current_service_id = ? AND next_service_id = ?", recParam.Current_service_id, recParam.Next_service_id).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rec = Entities.Recommendation{
				Current_service_id: recParam.Current_service_id,
				Next_service_id:    recParam.Next_service_id,
				Total:              1,
			}
			if err := d.db.Create(&rec).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		rec.Total++
		if err := d.db.Save(&rec).Error; err != nil {
			return err
		}
	}

	// Process RecommendationUtil record.
	var recUtil Entities.RecommendationUtil
	err = d.db.
		Where("current_service_id = ?", recParam.Current_service_id).
		First(&recUtil).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			recUtil = Entities.RecommendationUtil{
				Current_service_id: recParam.Current_service_id,
				Total:              1,
			}
			if err := d.db.Create(&recUtil).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		recUtil.Total++
		if err := d.db.Save(&recUtil).Error; err != nil {
			return err
		}
	}

	// Optionally: Invalidate cache entries if needed.
	// For example:
	// utils.Cache.Delete(fmt.Sprintf("suggest_next_service_%s", <userID>))
	// utils.Cache.Delete(fmt.Sprintf("interest_rating_%s", <userID>))
	// utils.Cache.Delete("bestseller_<weekKey>")

	return nil
}

// SuggestNextServie retrieves suggested services for a given user using cached full result set.
func (d *recommendationDriver) SuggestNextServie(userID string, config *model.Pagination) (*[]Entities.Service, int64, error) {
    if userID == "" {
        return nil, 0, errors.New("userID is required")
    }

    searchQuery := config.Search

    // Full Text Search branch if searchQuery is provided.
    if (searchQuery != "") {
        // 1. Get the latest VowRecord for the given user.
        var latestVow Entities.VowRecord
        if err := d.db.
            Where("user_id = ?", userID).
            Order("created_at desc").
            First(&latestVow).Error; err != nil {
            return nil, 0, err
        }
        currentServiceID := latestVow.ServiceID

        // 2. Build the subquery for recommendation score.
        subQuery := d.db.Table("recommendations r").
            Select("r.next_service_id, (r.total * 1.0 / ru.total) AS score").
            Joins("JOIN recommendation_utils ru ON ru.current_service_id = r.current_service_id").
            Where("r.current_service_id = ? AND r.next_service_id <> ?", currentServiceID, currentServiceID)

        // 3. Retrieve the full list with full text search ranking, similarity and combined weight.
        var fullServices []Entities.Service
        err := d.db.Model(&Entities.Service{}).
            Select(`
                services.*,
                COALESCE(sub.score, 0.0) as score,
                ts_rank(
                    setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') ||
                    setweight(to_tsvector('thai', coalesce(services.description, '')), 'B') ||
                    setweight(to_tsvector('thai', coalesce(services.address, '')), 'C'),
                    plainto_tsquery('thai', ?)
                ) as rank,
                similarity(services.name || ' ' || services.description || ' ' || services.address, ?) as sim,
                (COALESCE(sub.score, 0.0) * 0.3 + similarity(services.name || ' ' || services.description || ' ' || services.address, ?) * 0.7) as weight
            `, searchQuery, searchQuery, searchQuery).
            Joins("LEFT JOIN (?) as sub ON services.id = sub.next_service_id", subQuery).
            Where("services.id <> ? AND ( "+
                "(setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') || "+
                "setweight(to_tsvector('thai', coalesce(services.description, '')), 'B') || "+
                "setweight(to_tsvector('thai', coalesce(services.address, '')), 'C')) @@ plainto_tsquery('thai', ?) "+
                "OR similarity(services.name || ' ' || services.description || ' ' || services.address, ?) > 0.0)",
                currentServiceID, searchQuery, searchQuery).
            Order("weight DESC").
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&fullServices).Error
        if err != nil {
            return nil, 0, err
        }

        // Optionally, recalc the Rate field.
        for i := range fullServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", fullServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    fullServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }

        totalRecords := int64(len(fullServices))
        offset := (config.CurrentPage - 1) * config.PageSize
        if offset >= len(fullServices) {
            empty := []Entities.Service{}
            return &empty, totalRecords, nil
        }
        end := offset + config.PageSize
        if end > len(fullServices) {
            end = len(fullServices)
        }
        paginated := fullServices[offset:end]
        return &paginated, totalRecords, nil
    }

    // Original branch using cache when no searchQuery is provided.
    cacheKey := fmt.Sprintf("suggest_next_service_%s", userID)
    var fullServices []Entities.Service
    if cached, found := utils.Cache.Get(cacheKey); found {
        fullServices = cached.([]Entities.Service)
    } else {
        // 1. Get the latest VowRecord for the given user.
        var latestVow Entities.VowRecord
        if err := d.db.
            Where("user_id = ?", userID).
            Order("created_at desc").
            First(&latestVow).Error; err != nil {
            return nil, 0, err
        }
        currentServiceID := latestVow.ServiceID

        // 2. Build the subquery for recommendation score.
        subQuery := d.db.Table("recommendations r").
            Select("r.next_service_id, (r.total * 1.0 / ru.total) AS score").
            Joins("JOIN recommendation_utils ru ON ru.current_service_id = r.current_service_id").
            Where("r.current_service_id = ? AND r.next_service_id <> ?", currentServiceID, currentServiceID)

        // 3. Retrieve the full list of services (without pagination).
        err := d.db.Model(&Entities.Service{}).
            Select("services.*, COALESCE(sub.score, 0.0) as score").
            Joins("LEFT JOIN (?) as sub ON services.id = sub.next_service_id", subQuery).
            Where("services.id <> ?", currentServiceID).
            Order("COALESCE(sub.score, 0.0) DESC NULLS LAST").
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&fullServices).Error
        if err != nil {
            return nil, 0, err
        }

        // Optionally recalc the Rate field.
        for i := range fullServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", fullServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    fullServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }

        // Cache the full result set.
        utils.Cache.Set(cacheKey, fullServices, utils.DefaultExpiration)
    }

    totalRecords := int64(len(fullServices))
    offset := (config.CurrentPage - 1) * config.PageSize
    if offset >= len(fullServices) {
        empty := []Entities.Service{}
        return &empty, totalRecords, nil
    }
    end := offset + config.PageSize
    if end > len(fullServices) {
        end = len(fullServices)
    }
    paginated := fullServices[offset:end]
    return &paginated, totalRecords, nil
}

// InterestRating retrieves services based on the user's interests using caching.
func (d *recommendationDriver) InterestRating(userID string, config *model.Pagination) (*[]Entities.Service, int64, error) {
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }

    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, 0, err
    }

    searchQuery := config.Search

    // If a search query is provided, bypass the cache and perform full-text search.
    if searchQuery != "" {
        // Combined query for services with full text search ranking, similarity and combined weight.
        var combinedServices []Entities.Service
        err := d.db.Model(&Entities.Service{}).
            Joins("JOIN services_categories sc ON sc.service_id = services.id").
            Joins("JOIN review_utils ru ON ru.service_id = services.id").
            Select(`
                services.*,
                ts_rank(
                    setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') ||
                    setweight(to_tsvector('thai', coalesce(services.description, '')), 'B') ||
                    setweight(to_tsvector('thai', coalesce(services.address, '')), 'C'),
                    plainto_tsquery('thai', ?)
                ) as rank,
                similarity(services.name || ' ' || services.description || ' ' || services.address, ?) as sim,
                (ru.total_rete * 1.0 / NULLIF(ru.total_reviewer,0)) as rating,
                ((ru.total_rete * 1.0 / NULLIF(ru.total_reviewer,0)) * 0.3 + similarity(services.name || ' ' || services.description || ' ' || services.address, ?) * 0.7) as weight
            `, searchQuery, searchQuery, searchQuery).
            Where(`
                (
                    (setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') ||
                    setweight(to_tsvector('thai', coalesce(services.description, '')), 'B') ||
                    setweight(to_tsvector('thai', coalesce(services.address, '')), 'C')
                    ) @@ plainto_tsquery('thai', ?)
                    OR similarity(services.name || ' ' || services.description || ' ' || services.address, ?) > 0.0
                )
            `, searchQuery, searchQuery).
            Group("services.id, ru.total_rete, ru.total_reviewer").
            Order("weight DESC").
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&combinedServices).Error
        if err != nil {
            return nil, 0, err
        }
        totalRecords := int64(len(combinedServices))
        offset := (config.CurrentPage - 1) * config.PageSize
        if offset >= len(combinedServices) {
            empty := []Entities.Service{}
            return &empty, totalRecords, nil
        }
        end := offset + config.PageSize
        if end > len(combinedServices) {
            end = len(combinedServices)
        }
        paginated := combinedServices[offset:end]
        return &paginated, totalRecords, nil
    }

    // Original cached behavior when no search query is provided.
    cacheKey := fmt.Sprintf("interest_rating_%s", userID)
    var combinedServices []Entities.Service
    var totalRecords int64

    if cached, found := utils.Cache.Get(cacheKey); found {
        cachedData := cached.(struct {
            Services     []Entities.Service
            TotalRecords int64
        })
        combinedServices = cachedData.Services
        totalRecords = cachedData.TotalRecords
    } else {
        // 1. Count matching services.
        countMatchingSQL := `
            SELECT COUNT(DISTINCT s.id)
            FROM services s
            JOIN services_categories sc ON s.id = sc.service_id
            WHERE sc.category_id IN (
                SELECT category_id FROM interests WHERE user_id = ?
            )`
        var matchingCount int64
        if err := d.db.Raw(countMatchingSQL, uid).Scan(&matchingCount).Error; err != nil {
            return nil, 0, err
        }

        // 2. Count non-matching services.
        countNonMatchingSQL := `
            SELECT COUNT(DISTINCT s.id)
            FROM services s
            JOIN services_categories sc ON s.id = sc.service_id
            WHERE s.id NOT IN (
                SELECT s.id FROM services s
                JOIN services_categories sc ON s.id = sc.service_id
                WHERE sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?)
            )`
        var nonMatchingCount int64
        if err := d.db.Raw(countNonMatchingSQL, uid).Scan(&nonMatchingCount).Error; err != nil {
            return nil, 0, err
        }

        totalRecords = matchingCount + nonMatchingCount

        // 3. Retrieve matching services.
        var matchingServices []Entities.Service
        err = d.db.Model(&Entities.Service{}).
            Joins("JOIN services_categories sc ON sc.service_id = services.id").
            Joins("JOIN review_utils ru ON ru.service_id = services.id").
            Where("sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?)", uid).
            Group("services.id, ru.total_rete, ru.total_reviewer").
            Order("(ru.total_rete / NULLIF(ru.total_reviewer, 0.0)) DESC  NULLS LAST").
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&matchingServices).Error
        if err != nil {
            return nil, 0, err
        }
        for i := range matchingServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", matchingServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    matchingServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }

        // 4. Retrieve non-matching services.
        var nonMatchingServices []Entities.Service
        err = d.db.Model(&Entities.Service{}).
            Joins("JOIN services_categories sc ON sc.service_id = services.id").
            Joins("JOIN review_utils ru ON ru.service_id = services.id").
            Where("services.id NOT IN (SELECT s.id FROM services s JOIN services_categories sc ON s.id = sc.service_id WHERE sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?))", uid).
            Group("services.id, ru.total_rete, ru.total_reviewer").
            Order("NULLIF(ru.total_rete / NULLIF(ru.total_reviewer, 0.0), 0.0) DESC  NULLS LAST").
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&nonMatchingServices).Error
        if err != nil {
            return nil, 0, err
        }
        for i := range nonMatchingServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", nonMatchingServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    nonMatchingServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }

        // 5. Combine matching and non-matching services.
        combinedServices = append(matchingServices, nonMatchingServices...)

        // Cache the combined result set.
        utils.Cache.Set(cacheKey, struct {
            Services     []Entities.Service
            TotalRecords int64
        }{Services: combinedServices, TotalRecords: totalRecords}, utils.DefaultExpiration)
    }

    // Apply pagination.
    offset := (config.CurrentPage - 1) * config.PageSize
    if offset >= len(combinedServices) {
        empty := []Entities.Service{}
        return &empty, totalRecords, nil
    }
    end := offset + config.PageSize
    if end > len(combinedServices) {
        end = len(combinedServices)
    }
    paginated := combinedServices[offset:end]
    return &paginated, totalRecords, nil
}

func (d *recommendationDriver) Bestseller(config *model.Pagination) (*[]Entities.Service, int64, error) {
	// Determine the current week (assuming week starts on Monday).
	now := time.Now()
	weekday := int(now.Weekday())
	daysSinceMonday := (weekday + 6) % 7
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	// Cache key that includes the start date of the week.
	cacheKey := fmt.Sprintf("bestseller_%s", startOfWeek.Format("2006-01-02"))
	var results []struct {
		ServiceID uuid.UUID
		Count     int64
	}
	var services []Entities.Service

	if cached, found := utils.Cache.Get(cacheKey); found {
		fmt.Println("Cache hit for bestseller")
		cachedData := cached.(struct {
			Results  []struct {
				ServiceID uuid.UUID
				Count     int64
			}
			Services []Entities.Service
		})
		results = cachedData.Results
		services = cachedData.Services
	} else {
		// Aggregate transactions by service_id within the current week.
		if err := d.db.
			Model(&Entities.Order{}).
			Select("service_id, count(*) as count").
			Where("created_at >= ? AND created_at < ?", startOfWeek, endOfWeek).
			Group("service_id").
			Order("count DESC NULLS LAST").
			Find(&results).Error; err != nil {
			return nil, 0, err
		}

		// Build a list of service IDs.
		var serviceIDs []uuid.UUID
		for _, r := range results {
			serviceIDs = append(serviceIDs, r.ServiceID)
		}

		if err := d.db.
			Where("id in ?", serviceIDs).
			Preload("Review_utils").
			Preload("Categories").
			Preload("Packages").
			Preload("Packages.OrderType").
			Preload("Attachments").
			Find(&services).Error; err != nil {
			return nil, 0, err
		}

		// Recalculate the Rate field for each service.
		for i := range services {
			var reviewUtils Entities.Review_utils
			if err := d.db.Where("service_id = ?", services[i].ID).First(&reviewUtils).Error; err == nil {
				if reviewUtils.TotalReviewer > 0 {
					services[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
				}
			}
		}

		// Reorder services based on the sorted results from the aggregation.
		serviceMap := make(map[uuid.UUID]Entities.Service)
		for _, service := range services {
			serviceMap[service.ID] = service
		}
		var sortedServices []Entities.Service
		for _, r := range results {
			if service, ok := serviceMap[r.ServiceID]; ok {
				sortedServices = append(sortedServices, service)
			}
		}
		services = sortedServices

		// Cache the aggregated and sorted results.
		utils.Cache.Set(cacheKey, struct {
			Results  []struct {
				ServiceID uuid.UUID
				Count     int64
			}
			Services []Entities.Service
		}{Results: results, Services: services}, utils.DefaultExpiration)
	}

	// You might want to return the total count of orders as well.
	// For example, if needed:
	var totalCount int64
	for _, r := range results {
		totalCount += r.Count
	}
	return &services, totalCount, nil
}

