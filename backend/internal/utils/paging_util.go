package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type PaginationResponse struct {
	TotalPages  int64 `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
	CurrentPage int   `json:"currentPage"`
}

func Paginate(c *gin.Context, db *gorm.DB, result interface{}) (PaginationResponse, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

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
		TotalPages:  (totalItems + int64(pageSize) - 1) / int64(pageSize),
		TotalItems:  totalItems,
		CurrentPage: page,
	}, nil
}
