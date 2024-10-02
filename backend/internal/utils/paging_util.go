package utils

import (
	"gorm.io/gorm"
)

type PaginationResponse struct {
	TotalPages   int64 `json:"totalPages"`
	TotalItems   int64 `json:"totalItems"`
	CurrentPage  int   `json:"currentPage"`
	ItemsPerPage int   `json:"itemsPerPage"`
}

func Paginate(page int, pageSize int, db *gorm.DB, result interface{}) (PaginationResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	var totalItems int64
	if err := db.Count(&totalItems).Error; err != nil {
		return PaginationResponse{}, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(result).Error; err != nil {
		return PaginationResponse{}, err
	}

	return PaginationResponse{
		TotalPages:   (totalItems + int64(pageSize) - 1) / int64(pageSize),
		TotalItems:   totalItems,
		CurrentPage:  page,
		ItemsPerPage: pageSize,
	}, nil
}
