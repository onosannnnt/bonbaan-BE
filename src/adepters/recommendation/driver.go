package recommendationAdepter

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	recommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	"gorm.io/gorm"
)

type recommendationDriver struct {
    db *gorm.DB
}

func NewRecommendationDriver(db *gorm.DB) recommendationUsecase.RecommendationRepository {
    return &recommendationDriver{
        db: db,
    }
}

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

    return nil
}

func (d *recommendationDriver) SuggestNextServie(userID string, config *model.Pagination) (*[]Entities.Service, int64, error) {
    // Ensure a valid userID is provided.
    if userID == "" {
        return nil, 0, errors.New("userID is required")
    }

    // 1. Get the latest VowRecord for the given user.
    var latestVow Entities.VowRecord
    if err := d.db.
        Where("user_id = ?", userID).
        Order("created_at desc").
        First(&latestVow).Error; err != nil {
        return nil, 0, err
    }
    currentServiceID := latestVow.ServiceID

    // Validate pagination parameters.
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }
    offset := (config.CurrentPage - 1) * config.PageSize

    // 2. Count total services excluding the current service.
    countSQL := `SELECT COUNT(*) FROM services WHERE id <> ?`
    var totalRecords int64
    if err := d.db.Raw(countSQL, currentServiceID).Scan(&totalRecords).Error; err != nil {
        return nil, 0, err
    }

    // 3. Build a subquery for recommendation score.
    subQuery := d.db.Table("recommendations r").
        Select("r.next_service_id, (r.total * 1.0 / ru.total) AS score").
        Joins("JOIN recommendation_utils ru ON ru.current_service_id = r.current_service_id").
        Where("r.current_service_id = ? AND r.next_service_id <> ?", currentServiceID, currentServiceID)

    // 4. Retrieve paginated services using ORM chaining with preloads.
    var services []Entities.Service
    err := d.db.Model(&Entities.Service{}).
        Select("services.*, COALESCE(sub.score, 0) as score").
        Joins("LEFT JOIN (?) as sub ON services.id = sub.next_service_id", subQuery).
        Where("services.id <> ?", currentServiceID).
        Order("COALESCE(sub.score, 0) DESC").
        Limit(config.PageSize).
        Offset(offset).
        Preload("Review_utils").
        Preload("Categories").
        Preload("Packages").
        Preload("Packages.OrderType").
        Preload("Attachments").
        Find(&services).Error
    if err != nil {
        return nil, 0, err
    }

    // Optionally, recalc the Rate field from review_utils if needed.
    for i := range services {
        var reviewUtils Entities.Review_utils
        if err := d.db.Where("service_id = ?", services[i].ID).First(&reviewUtils).Error; err == nil {
            if reviewUtils.TotalReviewer > 0 {
                services[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
            }
        }
    }
    return &services, totalRecords, nil
}



func (d *recommendationDriver) InterestRating(userID string, config *model.Pagination) (*[]Entities.Service, int64, error) {
    // Validate pagination parameters.
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }
    
    // Parse the userID string into uuid.UUID.
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, 0, err
    }
    
    // Count matching services.
    countMatchingSQL := `
        SELECT COUNT(DISTINCT s.id)
        FROM services s
        JOIN services_categories sc ON s.id = sc.service_id
        WHERE sc.category_id IN (
            SELECT category_id FROM interests WHERE user_id = ?
        )
    `
    var matchingCount int64
    if err := d.db.Raw(countMatchingSQL, uid).Scan(&matchingCount).Error; err != nil {
        return nil, 0, err
    }
    
    // Count non-matching services.
    countNonMatchingSQL := `
        SELECT COUNT(DISTINCT s.id)
        FROM services s
        JOIN services_categories sc ON s.id = sc.service_id
        WHERE s.id NOT IN (
            SELECT s.id FROM services s
            JOIN services_categories sc ON s.id = sc.service_id
            WHERE sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?)
        )
    `
    var nonMatchingCount int64
    if err := d.db.Raw(countNonMatchingSQL, uid).Scan(&nonMatchingCount).Error; err != nil {
        return nil, 0, err
    }
    
    totalRecords := matchingCount + nonMatchingCount
    // Combined offset for the entire result set.
    offsetCombined := (config.CurrentPage - 1) * config.PageSize
    if int64(offsetCombined) >= totalRecords {
        // If the offset is beyond available records, return an empty slice.
        return &[]Entities.Service{}, totalRecords, nil
    }
    
    var combinedServices []Entities.Service

    // If the combined offset falls into the matching services.
    if int64(offsetCombined) < matchingCount {
        // Number of matching services available from the offset.
        matchingToFetch := int(math.Min(float64(matchingCount-int64(offsetCombined)), float64(config.PageSize)))
        
        var matchingServices []Entities.Service
        err = d.db.Model(&Entities.Service{}).
            Joins("JOIN services_categories sc ON sc.service_id = services.id").
            Joins("JOIN review_utils ru ON ru.service_id = services.id").
            Where("sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?)", uid).
            Group("services.id, ru.total_rete, ru.total_reviewer").
            Order("(ru.total_rete / NULLIF(ru.total_reviewer, 0)) DESC").
            Limit(matchingToFetch).
            Offset(offsetCombined).
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&matchingServices).Error
        if err != nil {
            return nil, 0, err
        }
        
        // Update the Rate field for matching services.
        for i := range matchingServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", matchingServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    matchingServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }
        
        combinedServices = append(combinedServices, matchingServices...)
        
        // Calculate how many more services are needed to fill the page.
        remaining := config.PageSize - len(matchingServices)
        if remaining > 0 {
            var nonMatchingServices []Entities.Service
            // For non-matching, since matching items were not enough, start at offset 0.
            err = d.db.Model(&Entities.Service{}).
                Joins("JOIN services_categories sc ON sc.service_id = services.id").
                Joins("JOIN review_utils ru ON ru.service_id = services.id").
                Where("services.id NOT IN (SELECT s.id FROM services s JOIN services_categories sc ON s.id = sc.service_id WHERE sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?))", uid).
                Group("services.id, ru.total_rete, ru.total_reviewer").
                Order("(ru.total_rete / NULLIF(ru.total_reviewer, 0)) DESC").
                Limit(remaining).
                Offset(0).
                Preload("Review_utils").
                Preload("Categories").
                Preload("Packages").
                Preload("Packages.OrderType").
                Preload("Attachments").
                Find(&nonMatchingServices).Error
            if err != nil {
                return nil, 0, err
            }
            
            // Update the Rate field for non-matching services.
            for i := range nonMatchingServices {
                var reviewUtils Entities.Review_utils
                if err := d.db.Where("service_id = ?", nonMatchingServices[i].ID).First(&reviewUtils).Error; err == nil {
                    if reviewUtils.TotalReviewer > 0 {
                        nonMatchingServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                    }
                }
            }
            combinedServices = append(combinedServices, nonMatchingServices...)
        }
    } else {
        // The combined offset falls entirely within the non-matching services.
        nonMatchingOffset := offsetCombined - int(matchingCount)
        var nonMatchingServices []Entities.Service
        err = d.db.Model(&Entities.Service{}).
            Joins("JOIN services_categories sc ON sc.service_id = services.id").
            Joins("JOIN review_utils ru ON ru.service_id = services.id").
            Where("services.id NOT IN (SELECT s.id FROM services s JOIN services_categories sc ON s.id = sc.service_id WHERE sc.category_id IN (SELECT category_id FROM interests WHERE user_id = ?))", uid).
            Group("services.id, ru.total_rete, ru.total_reviewer").
            Order("(ru.total_rete / NULLIF(ru.total_reviewer, 0)) DESC").
            Limit(config.PageSize).
            Offset(nonMatchingOffset).
            Preload("Review_utils").
            Preload("Categories").
            Preload("Packages").
            Preload("Packages.OrderType").
            Preload("Attachments").
            Find(&nonMatchingServices).Error
        if err != nil {
            return nil, 0, err
        }
        
        // Update the Rate field for non-matching services.
        for i := range nonMatchingServices {
            var reviewUtils Entities.Review_utils
            if err := d.db.Where("service_id = ?", nonMatchingServices[i].ID).First(&reviewUtils).Error; err == nil {
                if reviewUtils.TotalReviewer > 0 {
                    nonMatchingServices[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
                }
            }
        }
        combinedServices = append(combinedServices, nonMatchingServices...)
    }
    
    return &combinedServices, totalRecords, nil
}

func (d *recommendationDriver) Bestseller(config *model.Pagination) (*[]Entities.Service, int64, error) {
    // Determine the start and end of the current week (assuming week starts on Monday)
    now := time.Now()
    weekday := int(now.Weekday())
    // Adjust so that Monday is the start (Monday = 0)
    daysSinceMonday := (weekday + 6) % 7
    startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
    endOfWeek := startOfWeek.AddDate(0, 0, 7)

    // Aggregate transactions by service_id within the current week.
    type result struct {
        ServiceID uuid.UUID
        Count     int64
    }
    var results []result

    if err := d.db.
        Model(&Entities.Order{}).
        Select("service_id, count(*) as count").
        Where("created_at >= ? AND created_at < ?", startOfWeek, endOfWeek).
        Group("service_id").
        Order("count DESC").
        Find(&results).Error; err != nil {
        return nil, 0, err
    }

    totalRecords := int64(len(results))

    // Apply pagination.
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }
    startIndex := (config.CurrentPage - 1) * config.PageSize
    endIndex := startIndex + config.PageSize
    if startIndex > len(results) {
        return &[]Entities.Service{}, totalRecords, nil
    }
    if endIndex > len(results) {
        endIndex = len(results)
    }
    paginatedResults := results[startIndex:endIndex]

    var serviceIDs []uuid.UUID
    for _, r := range paginatedResults {
        serviceIDs = append(serviceIDs, r.ServiceID)
    }

    // Retrieve services with associations preloaded.
    var services []Entities.Service
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

    // Optionally, recalc the Rate field from review_utils if available.
    for i := range services {
        var reviewUtils Entities.Review_utils
        if err := d.db.Where("service_id = ?", services[i].ID).First(&reviewUtils).Error; err == nil {
            if reviewUtils.TotalReviewer > 0 {
                services[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
            }
        }
    }
    return &services, totalRecords, nil
}
