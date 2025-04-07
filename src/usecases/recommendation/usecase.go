package recommendationUsecase

import (
	"math"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type RecommendationUsecase interface {
    Insert(recommendation *Entities.Recommendation) error
    SuggestNextServies(userID string, config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error)
    InterestRatings(userID string, config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error)
    Bestsellers(config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error)
}

type recommendationService struct {
    recommendationRepository RecommendationRepository
}

func NewRecommendationService(recommendationRepository RecommendationRepository) RecommendationUsecase {
    return &recommendationService{
        recommendationRepository: recommendationRepository,
    }
}

func (r *recommendationService) Insert(recommendation *Entities.Recommendation) error {
    return r.recommendationRepository.Insert(recommendation)
}

func mapServiceToOutput(s Entities.Service) model.ServiceOutput {
	// Convert categories
	categories := make([]model.CategoryOutput, len(s.Categories))
	for i, c := range s.Categories {
		categories[i] = mapCategoryToOutput(c)
	}

	// Convert packages
	packages := make([]model.PackageOutput, len(s.Packages))
	for i, p := range s.Packages {
		packages[i] = mapPackageToOutput(p)
	}

	// Convert attachments
	attachments := make([]model.AttachmentOutput, len(s.Attachments))
	for i, a := range s.Attachments {
		attachments[i] = mapAttachmentToOutput(a)
	}

	return model.ServiceOutput{
		ID:          s.ID.String(),
		Name:        s.Name,
		Description: s.Description,
		Rate:        s.Rate,
		Address:     s.Address,
		Categories:  categories,
		Packages:    packages,
		Attachments: attachments,
		UpdateAt: s.UpdatedAt.String(),
		CreateAt: s.CreatedAt.String(),
	}
}
// mapCategoryToOutput converts an Entities.Category to a model.CategoryOutput.
func mapCategoryToOutput(c Entities.Category) model.CategoryOutput {
	return model.CategoryOutput{
		ID:  c.ID.String(),
		Name: c.Name,
	}
}

// mapPackageToOutput converts an Entities.Package to a model.PackageOutput.
func mapPackageToOutput(p Entities.Package) model.PackageOutput {
	return model.PackageOutput{
		Name:        p.Name,
		Item:        p.Item,
		Price:       p.Price,
		Description: p.Description,
		OrderTypeID: p.OrderTypeID.String(),
	}
}

// mapAttachmentToOutput converts an Entities.Attachment to a model.AttachmentOutput.
func mapAttachmentToOutput(a Entities.Attachment) model.AttachmentOutput {
	return model.AttachmentOutput{
		ID:  a.ID.String(),
		URL: a.URL,
	}
}



func (r *recommendationService) SuggestNextServies(userID string, config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error) {
    
    if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}

	services, totalRecords, err := r.recommendationRepository.SuggestNextServie(userID, config)
    if err != nil {
        return nil, nil, err
    }

    outputs := make([]model.ServiceOutput, 0, len(*services))
    for _, svc := range *services {
        outputs = append(outputs, mapServiceToOutput(svc))
    }

    totalPages := math.Ceil(float64(totalRecords) / float64(config.PageSize))
    pagination := &model.Pagination{
        CurrentPage:  config.CurrentPage,
        PageSize:     config.PageSize,
        TotalRecords: int(totalRecords),
        TotalPages:   int(totalPages),
    }

    return &outputs, pagination, nil
}

func (r *recommendationService) InterestRatings(userID string,config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error) {
    if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}

    services, totalRecords, err := r.recommendationRepository.InterestRating(userID,config)
    if err != nil {
        return nil, nil, err
    }

    outputs := make([]model.ServiceOutput, 0, len(*services))
    for _, svc := range *services {
        outputs = append(outputs, mapServiceToOutput(svc))
    }

    totalPages := math.Ceil(float64(totalRecords) / float64(config.PageSize))
    pagination := &model.Pagination{
        CurrentPage:  config.CurrentPage,
        PageSize:     config.PageSize,
        TotalRecords: int(totalRecords),
        TotalPages:   int(totalPages),
    }

    return &outputs, pagination, nil
}

func (r *recommendationService) Bestsellers(config *model.Pagination) (*[]model.ServiceOutput, *model.Pagination, error) {
    if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
    
    services, totalRecords, err := r.recommendationRepository.Bestseller(config)
    if err != nil {
        return nil, nil, err
    }

    outputs := make([]model.ServiceOutput, 0, len(*services))
    for _, svc := range *services {
        outputs = append(outputs, mapServiceToOutput(svc))
    }

    totalPages := math.Ceil(float64(totalRecords) / float64(config.PageSize))
    pagination := &model.Pagination{
        CurrentPage:  config.CurrentPage,
        PageSize:     config.PageSize,
        TotalRecords: int(totalRecords),
        TotalPages:   int(totalPages),
    }

    return &outputs, pagination, nil
}
