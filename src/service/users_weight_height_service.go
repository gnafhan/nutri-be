package service

import (
	"app/src/model"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersWeightHeightService interface {
	AddWeightHeight(ctx *fiber.Ctx, record *model.UsersWeightHeightHistory) (*model.UsersWeightHeightHistory, error)
	GetWeightHeights(ctx *fiber.Ctx, userID uuid.UUID) ([]model.UsersWeightHeightHistory, error)
	GetWeightHeightByID(ctx *fiber.Ctx, recordID string, userID uuid.UUID) (*model.UsersWeightHeightHistory, error)
	UpdateWeightHeight(ctx *fiber.Ctx, recordID string, record *model.UsersWeightHeightHistory) (*model.UsersWeightHeightHistory, error)
	DeleteWeightHeight(ctx *fiber.Ctx, recordID string, userID uuid.UUID) error
}

type usersWeightHeightService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUsersWeightHeightService(db *gorm.DB) UsersWeightHeightService {
	return &usersWeightHeightService{
		Log: logrus.New(),
		DB:  db,
	}
}

func (s *usersWeightHeightService) AddWeightHeight(ctx *fiber.Ctx, record *model.UsersWeightHeightHistory) (*model.UsersWeightHeightHistory, error) {
	if record.RecordedAt.IsZero() {
		record.RecordedAt = time.Now()
	}

	if err := s.DB.WithContext(ctx.Context()).Create(record).Error; err != nil {
		s.Log.Errorf("Failed to add weight and height record: %+v", err)
		return nil, err
	}

	if err := s.updateUserHeightWeight(ctx, record.UserID); err != nil {
		s.Log.Errorf("Failed to update user's height and weight: %+v", err)
		return nil, err
	}

	return record, nil
}

func (s *usersWeightHeightService) GetWeightHeights(ctx *fiber.Ctx, userID uuid.UUID) ([]model.UsersWeightHeightHistory, error) {
	var records []model.UsersWeightHeightHistory
	if err := s.DB.WithContext(ctx.Context()).
		Where("user_id = ?", userID).
		Order("recorded_at DESC").
		Find(&records).Error; err != nil {
		s.Log.Errorf("Failed to get weight and height records: %+v", err)
		return nil, err
	}

	return records, nil
}

func (s *usersWeightHeightService) GetWeightHeightByID(ctx *fiber.Ctx, recordID string, userID uuid.UUID) (*model.UsersWeightHeightHistory, error) {
	existingRecord := new(model.UsersWeightHeightHistory)
	if err := s.DB.WithContext(ctx.Context()).First(existingRecord, "id = ? AND user_id = ?", recordID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Record not found")
		}
		s.Log.Errorf("Failed to find record: %+v", err)
		return nil, err
	}

	return existingRecord, nil
}

func (s *usersWeightHeightService) UpdateWeightHeight(ctx *fiber.Ctx, recordID string, record *model.UsersWeightHeightHistory) (*model.UsersWeightHeightHistory, error) {
	existingRecord := new(model.UsersWeightHeightHistory)
	if err := s.DB.WithContext(ctx.Context()).First(existingRecord, "id = ? AND user_id = ?", recordID, record.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Record not found")
		}
		s.Log.Errorf("Failed to find record: %+v", err)
		return nil, err
	}

	existingRecord.Weight = record.Weight
	existingRecord.Height = record.Height
	if !record.RecordedAt.IsZero() {
		existingRecord.RecordedAt = record.RecordedAt
	}

	if err := s.DB.WithContext(ctx.Context()).Save(existingRecord).Error; err != nil {
		s.Log.Errorf("Failed to update weight and height record: %+v", err)
		return nil, err
	}

	if err := s.updateUserHeightWeight(ctx, record.UserID); err != nil {
		s.Log.Errorf("Failed to update user's height and weight: %+v", err)
		return nil, err
	}

	return existingRecord, nil
}

func (s *usersWeightHeightService) DeleteWeightHeight(ctx *fiber.Ctx, recordID string, userID uuid.UUID) error {
	existingRecord := new(model.UsersWeightHeightHistory)
	if err := s.DB.WithContext(ctx.Context()).First(existingRecord, "id = ? AND user_id = ?", recordID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Record not found")
		}
		s.Log.Errorf("Failed to find record: %+v", err)
		return err
	}

	if err := s.DB.WithContext(ctx.Context()).Delete(existingRecord).Error; err != nil {
		s.Log.Errorf("Failed to delete weight and height record: %+v", err)
		return err
	}

	if err := s.updateUserHeightWeight(ctx, userID); err != nil {
		s.Log.Errorf("Failed to update user's height and weight: %+v", err)
		return err
	}

	return nil
}

func (s *usersWeightHeightService) updateUserHeightWeight(ctx *fiber.Ctx, userID uuid.UUID) error {
	var latestRecord model.UsersWeightHeightHistory
	if err := s.DB.WithContext(ctx.Context()).
		Where("user_id = ?", userID).
		Order("recorded_at DESC").
		First(&latestRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.DB.WithContext(ctx.Context()).
				Model(&model.User{}).
				Where("id = ?", userID).
				Updates(map[string]interface{}{"height": 0, "weight": 0}).Error; err != nil {
				s.Log.Errorf("Failed to reset user's height and weight: %+v", err)
				return err
			}
			return nil
		}
		s.Log.Errorf("Failed to find latest record: %+v", err)
		return err
	}

	if err := s.DB.WithContext(ctx.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{"height": latestRecord.Height, "weight": latestRecord.Weight}).Error; err != nil {
		s.Log.Errorf("Failed to update user's height and weight: %+v", err)
		return err
	}

	return nil
}
