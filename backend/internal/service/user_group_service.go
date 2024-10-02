package service

import (
	"errors"
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/dto"
	"github.com/stonith404/pocket-id/backend/internal/model"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"gorm.io/gorm"
)

type UserGroupService struct {
	db *gorm.DB
}

func NewUserGroupService(db *gorm.DB) *UserGroupService {
	return &UserGroupService{db: db}
}

func (s *UserGroupService) List(name string, page int, pageSize int) (groups []model.UserGroup, response utils.PaginationResponse, err error) {
	query := s.db.Model(&model.UserGroup{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	response, err = utils.Paginate(page, pageSize, query, &groups)
	return groups, response, err
}

func (s *UserGroupService) Get(id string) (group model.UserGroup, err error) {
	err = s.db.Where("id = ?", id).Preload("Users").First(&group).Error
	return group, err
}

func (s *UserGroupService) Delete(id string) error {
	var group model.UserGroup
	if err := s.db.Where("id = ?", id).First(&group).Error; err != nil {
		return err
	}

	return s.db.Delete(&group).Error
}

func (s *UserGroupService) Create(input dto.UserGroupCreateDto) (group model.UserGroup, err error) {
	group = model.UserGroup{
		FriendlyName: input.FriendlyName,
		Name:         input.Name,
	}

	if err := s.db.Preload("Users").Create(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.UserGroup{}, common.ErrNameAlreadyInUse
		}
		return model.UserGroup{}, err
	}
	return group, nil
}

func (s *UserGroupService) Update(id string, input dto.UserGroupCreateDto) (group model.UserGroup, err error) {
	group, err = s.Get(id)
	if err != nil {
		return model.UserGroup{}, err
	}

	group.Name = input.Name
	group.FriendlyName = input.FriendlyName

	if err := s.db.Preload("Users").Save(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.UserGroup{}, common.ErrNameAlreadyInUse
		}
		return model.UserGroup{}, err
	}
	return group, nil
}

func (s *UserGroupService) UpdateUsers(id string, input dto.UserGroupUpdateUsersDto) (group model.UserGroup, err error) {
	group, err = s.Get(id)
	if err != nil {
		return model.UserGroup{}, err
	}

	// Fetch the users based on UserIDs in input
	var users []model.User
	if len(input.UserIDs) > 0 {
		if err := s.db.Where("id IN (?)", input.UserIDs).Find(&users).Error; err != nil {
			return model.UserGroup{}, err
		}
	}

	// Replace the current users with the new set of users
	if err := s.db.Model(&group).Association("Users").Replace(users); err != nil {
		return model.UserGroup{}, err
	}

	// Save the updated group
	if err := s.db.Save(&group).Error; err != nil {
		return model.UserGroup{}, err
	}

	return group, nil
}

func (s *UserGroupService) GetUserCountOfGroup(id string) (int64, error) {
	var group model.UserGroup
	if err := s.db.Preload("Users").Where("id = ?", id).First(&group).Error; err != nil {
		return 0, err
	}
	return s.db.Model(&group).Association("Users").Count(), nil
}
