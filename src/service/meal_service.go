package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"time"

	"app/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MealService interface {
	ScanMeal(c *fiber.Ctx, imageFile *multipart.FileHeader, userID uuid.UUID) (*MealScanResponse, error)
}

type mealService struct {
	Log     *logrus.Logger
	DB      *gorm.DB
	ApiKey  string
	BaseURL string
}

func NewMealService(db *gorm.DB, apiKey, baseURL string) *mealService {
	return &mealService{
		Log:     logrus.New(),
		DB:      db,
		ApiKey:  apiKey,
		BaseURL: baseURL,
	}
}

// Response structure
type NutrientDetail struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type Nutrient struct {
	Calories NutrientDetail `json:"calories"`
	Protein  NutrientDetail `json:"protein"`
	Carbs    NutrientDetail `json:"carbs"`
	Fat      NutrientDetail `json:"fat"`
}

type MealScanResponse struct {
	Foods     [][]string `json:"foods"`
	TotalNutr Nutrient   `json:"total_nutrient"`
}

// ScanMeal handles the image scanning process
func (s *mealService) ScanMeal(c *fiber.Ctx, imageFile *multipart.FileHeader, userID uuid.UUID) (*MealScanResponse, error) {
	file, err := imageFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Step 1: Upload Image to Segmentation API
	imageId, foods, err := s.uploadImageToSegmentationAPI(file, imageFile.Filename)
	if err != nil {
		return nil, err
	}

	// Step 2: Convert imageId to string and get Nutrition Info
	imageIdStr := strconv.Itoa(imageId)
	totalNutr, err := s.getNutritionInfo(imageIdStr)
	if err != nil {
		return nil, err
	}

	// Step 3: Simpan hasil scan ke database (MealHistory & MealHistoryDetail)
	if err := s.saveMealHistory(userID, foods, totalNutr); err != nil {
		return nil, err
	}

	return &MealScanResponse{
		Foods:     foods,
		TotalNutr: totalNutr,
	}, nil
}

// Upload image to segmentation API and extract food names
func (s *mealService) uploadImageToSegmentationAPI(file io.Reader, filename string) (int, [][]string, error) {
	url := fmt.Sprintf("%s/v2/image/segmentation/complete/v1.1?language=eng", s.BaseURL)
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return 0, nil, err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil, errors.New("failed to upload image")
	}

	// Struct for response parsing
	var result struct {
		ImageID          int `json:"imageId"`
		SegmentationData []struct {
			RecognitionResults []struct {
				Name        string  `json:"name"`
				Probability float64 `json:"probability"`
			} `json:"recognition_results"`
		} `json:"segmentation_results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, nil, err
	}

	// Group food names per area (taking top 3 highest probability)
	var groupedFoodNames [][]string
	for _, segment := range result.SegmentationData {
		recognitionResults := segment.RecognitionResults

		// Sort food items by probability (descending)
		sort.Slice(recognitionResults, func(i, j int) bool {
			return recognitionResults[i].Probability > recognitionResults[j].Probability
		})

		// Take up to 3 highest probability foods
		var topFoods []string
		for i := 0; i < len(recognitionResults) && i < 3; i++ {
			topFoods = append(topFoods, recognitionResults[i].Name)
		}

		// Append to grouped list
		groupedFoodNames = append(groupedFoodNames, topFoods)
	}

	return result.ImageID, groupedFoodNames, nil
}

// Fetch nutrition info based on imageId
func (s *mealService) getNutritionInfo(imageId string) (Nutrient, error) {
	url := fmt.Sprintf("%s/v2/nutrition/recipe/nutritionalInfo/v1.1?language=eng", s.BaseURL)
	payload, _ := json.Marshal(map[string]string{"imageId": imageId})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return Nutrient{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Nutrient{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Nutrient{}, errors.New("failed to get nutrition info")
	}

	// Parsing response dengan map[string]interface{}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Nutrient{}, err
	}

	nutritionalInfo := result["nutritional_info"].(map[string]interface{})
	totalNutrients := nutritionalInfo["totalNutrients"].(map[string]interface{})

	extractNutrientDetail := func(nutrientKey string) NutrientDetail {
		if nutrient, exists := totalNutrients[nutrientKey].(map[string]interface{}); exists {
			quantity, _ := nutrient["quantity"].(float64)
			unit, _ := nutrient["unit"].(string)
			return NutrientDetail{Quantity: quantity, Unit: unit}
		}
		return NutrientDetail{Quantity: 0, Unit: ""}
	}

	return Nutrient{
		Calories: extractNutrientDetail("ENERC_KCAL"),
		Protein:  extractNutrientDetail("PROCNT"),
		Carbs:    extractNutrientDetail("CHOCDF"),
		Fat:      extractNutrientDetail("FAT"),
	}, nil
}

// Save meal history and details
func (s *mealService) saveMealHistory(userID uuid.UUID, foods [][]string, totalNutr Nutrient) error {
	mealHistory := model.MealHistory{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Scanned Meal",
		MealTime:  time.Now(),
		Calories:  totalNutr.Calories.Quantity,
		Protein:   totalNutr.Protein.Quantity,
		Carbs:     totalNutr.Carbs.Quantity,
		Fat:       totalNutr.Fat.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.DB.Create(&mealHistory).Error; err != nil {
		return err
	}

	saveMaps := make(map[string]any)
	saveMaps["foods"] = foods
	saveMaps["nutrients"] = totalNutr
	saveJSON, _ := json.Marshal(saveMaps)
	mealDetail := model.MealHistoryDetail{
		ID:            uuid.New(),
		MealHistoryID: mealHistory.ID,
		APIResult:     string(saveJSON),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.DB.Create(&mealDetail).Error
}
