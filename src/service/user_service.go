package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error)
	GetUserByID(c *fiber.Ctx, id string) (*model.User, error)
	GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error)
	CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error)
	UpdatePassOrVerify(c *fiber.Ctx, req *validation.UpdatePassOrVerify, id string) error
	UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string, file *multipart.FileHeader, hasFile bool) (*model.User, error)
	DeleteUser(c *fiber.Ctx, id string) error
	CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error)
}

type userService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserService(db *gorm.DB, validate *validator.Validate) UserService {
	return &userService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *userService) GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error) {
	var users []model.User
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR role LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&users).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search users: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Limit(params.Limit).Offset(offset).Find(&users)
	if result.Error != nil {
		s.Log.Errorf("Failed to get all users: %+v", result.Error)
		return nil, 0, result.Error
	}

	return users, totalResults, result.Error
}

func (s *userService) GetUserByID(c *fiber.Ctx, id string) (*model.User, error) {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).First(user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get user by id: %+v", result.Error)
	}

	return user, result.Error
}

func (s *userService) GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error) {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).Where("email = ?", email).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get user by email: %+v", result.Error)
	}

	return user, result.Error
}

func (s *userService) CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Failed hash password: %+v", err)
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	result := s.DB.WithContext(c.Context()).Create(user)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Email is already in use")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create user: %+v", result.Error)
	}

	return user, result.Error
}

func (s *userService) UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string, file *multipart.FileHeader, hasFile bool) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Email == "" && req.Name == "" && req.Password == "" && !hasFile &&
		req.BirthDate == nil && req.Height == nil && req.Weight == nil &&
		req.Gender == nil && req.ActivityLevel == nil && req.MedicalHistory == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	currentUser, err := s.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}

	tx := s.DB.WithContext(c.Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updateBody := &model.User{}

	if req.Name != "" {
		updateBody.Name = req.Name
	}
	if req.Email != "" {
		updateBody.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		updateBody.Password = hashedPassword
	}
	if req.BirthDate != nil {
		updateBody.BirthDate = req.BirthDate
	}
	if req.Height != nil {
		updateBody.Height = req.Height
	}
	if req.Weight != nil {
		updateBody.Weight = req.Weight
	}
	if req.Gender != nil {
		updateBody.Gender = req.Gender
	}
	if req.ActivityLevel != nil {
		updateBody.ActivityLevel = req.ActivityLevel
	}
	if req.MedicalHistory != nil {
		updateBody.MedicalHistory = req.MedicalHistory
	}

	if hasFile && file != nil {
		if err := validateImageFile(file); err != nil {
			return nil, err
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		uploadPath := "./uploads/profile-pictures/" + filename

		if err := os.MkdirAll("./uploads/profile-pictures", os.ModePerm); err != nil {
			s.Log.Errorf("Failed to create directory: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file")
		}

		if err := c.SaveFile(file, uploadPath); err != nil {
			s.Log.Errorf("Failed to save file: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file")
		}

		profilePictureURL := "/uploads/profile-pictures/" + filename
		updateBody.ProfilePicture = profilePictureURL

		if currentUser.ProfilePicture != "" && !strings.Contains(currentUser.ProfilePicture, "default") {
			oldFilePath := "." + currentUser.ProfilePicture
			if _, err := os.Stat(oldFilePath); err == nil {
				os.Remove(oldFilePath)
			}
		}
	}

	if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(updateBody).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fiber.NewError(fiber.StatusConflict, "Email is already in use")
		}
		s.Log.Errorf("Failed to update user: %+v", err)
		return nil, err
	}

	if req.Weight != nil || req.Height != nil {
		weight := currentUser.Weight
		if req.Weight != nil {
			weight = req.Weight
		}

		height := currentUser.Height
		if req.Height != nil {
			height = req.Height
		}

		userWeightHeight := &model.UsersWeightHeightHistory{
			ID:         uuid.New(),
			UserID:     currentUser.ID,
			Weight:     *weight,
			Height:     *height,
			RecordedAt: time.Now(),
		}

		if err := tx.Create(userWeightHeight).Error; err != nil {
			tx.Rollback()
			s.Log.Errorf("Failed to add weight height to database: %+v", err)
			return nil, err
		}

		if err := tx.Model(&model.User{}).
			Where("id = ?", currentUser.ID).
			Updates(map[string]interface{}{"weight": *weight, "height": *height}).Error; err != nil {
			tx.Rollback()
			s.Log.Errorf("Failed to update user's weight and height: %+v", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, err
	}

	updatedUser, err := s.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func validateImageFile(file *multipart.FileHeader) error {
	if file.Size > 5*1024*1024 {
		return fiber.NewError(fiber.StatusBadRequest, "File too large (max 5MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExtensions[ext] {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file type. Only JPG, JPEG, PNG, GIF, and WEBP are allowed")
	}

	fileHeader, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process file")
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil && err != io.EOF {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file")
	}

	fileHeader.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !validTypes[contentType] {
		return fiber.NewError(fiber.StatusBadRequest, "File content doesn't match a valid image format")
	}

	return nil
}

func (s *userService) UpdatePassOrVerify(c *fiber.Ctx, req *validation.UpdatePassOrVerify, id string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	if req.Password == "" && !req.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Request")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		req.Password = hashedPassword
	}

	updateBody := &model.User{
		Password:      req.Password,
		VerifiedEmail: req.VerifiedEmail,
	}

	result := s.DB.WithContext(c.Context()).Where("id = ?", id).Updates(updateBody)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to update user password or verifiedEmail: %+v", result.Error)
	}

	return result.Error
}

func (s *userService) DeleteUser(c *fiber.Ctx, id string) error {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).Delete(user, "id = ?", id)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to delete user: %+v", result.Error)
	}

	return result.Error
}

func (s *userService) CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userFromDB, err := s.GetUserByEmail(c, req.Email)
	if err != nil {
		if err.Error() == "User not found" {
			user := &model.User{
				Name:           req.Name,
				Email:          req.Email,
				VerifiedEmail:  true,
				ProfilePicture: req.ProfilePicture,
				GoogleIDToken:  req.GoogleIDToken,
			}

			if createErr := s.DB.WithContext(c.Context()).Create(user).Error; createErr != nil {
				s.Log.Errorf("Failed to create user: %+v", createErr)
				return nil, createErr
			}

			return user, nil
		}

		return nil, err
	}

	if updateErr := s.DB.WithContext(c.Context()).Save(userFromDB).Error; updateErr != nil {
		s.Log.Errorf("Failed to update user: %+v", updateErr)
		return nil, updateErr
	}

	return userFromDB, nil
}
