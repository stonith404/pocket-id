package service

import (
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"gorm.io/gorm"
)

// Reserved claims
var reservedClaims = map[string]struct{}{
	"given_name":         {},
	"family_name":        {},
	"name":               {},
	"email":              {},
	"preferred_username": {},
	"groups":             {},
	"sub":                {},
	"iss":                {},
	"aud":                {},
	"exp":                {},
	"iat":                {},
	"auth_time":          {},
	"nonce":              {},
	"acr":                {},
	"amr":                {},
	"azp":                {},
	"nbf":                {},
	"jti":                {},
}

type CustomClaimService struct {
	db *gorm.DB
}

func NewCustomClaimService(db *gorm.DB) *CustomClaimService {
	return &CustomClaimService{db: db}
}

// isReservedClaim checks if a claim key is reserved e.g. email, preferred_username
func isReservedClaim(key string) bool {
	_, ok := reservedClaims[key]
	return ok
}

func (s *CustomClaimService) UpdateUserCustomClaims(userID string, claims []dto.CustomClaimCreateDto) ([]model.CustomClaim, error) {
	var existingClaims []model.CustomClaim
	err := s.db.Where("user_id = ?", userID).Find(&existingClaims).Error
	if err != nil {
		return nil, err
	}

	// Delete claims that are not in the new list
	for _, existingClaim := range existingClaims {
		found := false
		for _, claim := range claims {
			if claim.Key == existingClaim.Key {
				found = true
				break
			}
		}
		if !found {
			err = s.db.Delete(&existingClaim).Error
			if err != nil {
				return nil, err
			}
		}
	}

	// Add or update claims
	for _, claim := range claims {
		if isReservedClaim(claim.Key) {
			return nil, &common.ReservedClaimError{Key: claim.Key}
		}
		// Update the claim if it already exists or create a new one
		err = s.db.Where("user_id = ? AND key = ?", userID, claim.Key).Assign(&model.CustomClaim{
			Key:    claim.Key,
			Value:  claim.Value,
			UserID: userID,
		}).FirstOrCreate(&model.CustomClaim{}).Error
		if err != nil {
			return nil, err
		}
	}

	// Get the updated claims
	var updatedClaims []model.CustomClaim
	err = s.db.Where("user_id = ?", userID).Find(&updatedClaims).Error
	if err != nil {
		return nil, err
	}

	return updatedClaims, nil
}

func (s *CustomClaimService) GetCustomClaimsForUser(userID string) ([]model.CustomClaim, error) {
	var customClaims []model.CustomClaim
	err := s.db.Where("user_id = ?", userID).Find(&customClaims).Error
	return customClaims, err
}

func (s *CustomClaimService) GetSuggestions() ([]string, error) {
	var customClaimsKeys []string

	err := s.db.Model(&model.CustomClaim{}).
		Group("key").
		Order("COUNT(*) DESC").
		Pluck("key", &customClaimsKeys).Error

	return customClaimsKeys, err
}
