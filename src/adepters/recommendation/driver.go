package recommendationAdepter

import (
	"errors"
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

    // Use the ServiceID of the latest vow record as the current_service_id.
    currentServiceID := latestVow.ServiceID

    // Validate pagination parameters.
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }
    offset := (config.CurrentPage - 1) * config.PageSize

    // 2. Count total recommended services (exclude cases where next_service_id equals current_service_id).
    countSQL := `
        SELECT COUNT(*) FROM (
            SELECT r.next_service_id
            FROM recommendations r
            JOIN recommendationutils ru ON ru.current_service_id = r.current_service_id
            WHERE r.current_service_id = ? AND r.next_service_id <> ?
        ) AS countSub
    `
    var totalRecords int64
    if err := d.db.Raw(countSQL, currentServiceID, currentServiceID).Scan(&totalRecords).Error; err != nil {
        return nil, 0, err
    }

    // 3. Retrieve paginated recommendations.
    querySQL := `
        SELECT s.*
        FROM services s
        JOIN (
            SELECT r.next_service_id, (r.total * 1.0 / ru.total) AS score
            FROM recommendations r
            JOIN recommendationutils ru ON ru.current_service_id = r.current_service_id
            WHERE r.current_service_id = ? AND r.next_service_id <> ?
        ) AS sub ON s.id = sub.next_service_id
        ORDER BY sub.score DESC
        LIMIT ? OFFSET ?
    `
    var services []Entities.Service
    if err := d.db.Raw(querySQL, currentServiceID, currentServiceID, config.PageSize, offset).Scan(&services).Error; err != nil {
        return nil, 0, err
    }

    return &services, totalRecords, nil
}

func (d *recommendationDriver) InterestRating(userID string, config *model.Pagination) (*[]Entities.Service, int64, error) {
    // Validate pagination parameters.
    if config.CurrentPage < 1 || config.PageSize < 1 {
        return nil, 0, errors.New("invalid pagination parameters")
    }
    offset := (config.CurrentPage - 1) * config.PageSize

    // Parse the userID string into uuid.UUID.
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, 0, err
    }

    // This query retrieves services whose categories match the user's interests.
    // It assumes a join table "service_categories" exists mapping services to categories.
    countSQL := `
        SELECT COUNT(DISTINCT s.id)
        FROM services s
        JOIN services_categories sc ON s.id = sc.service_id
        WHERE sc.category_id IN (
            SELECT category_id FROM interests WHERE user_id = ?
        )
    `
    var totalRecords int64
    if err := d.db.Raw(countSQL, uid).Scan(&totalRecords).Error; err != nil {
        return nil, 0, err
    }

    querySQL := `
        SELECT s.*
        FROM services s
        JOIN services_categories sc ON s.id = sc.service_id
        WHERE sc.category_id IN (
            SELECT category_id FROM interests WHERE user_id = ?
        )
        GROUP BY s.id
        ORDER BY COUNT(sc.category_id) DESC
        LIMIT ? OFFSET ?
    `
    var services []Entities.Service
    if err := d.db.Raw(querySQL, uid, config.PageSize, offset).Scan(&services).Error; err != nil {
        return nil, 0, err
    }

    return &services, totalRecords, nil
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

    var services []Entities.Service
    if err := d.db.
        Where("id in ?", serviceIDs).
        Find(&services).Error; err != nil {
        return nil, 0, err
    }

    return &services, totalRecords, nil
}