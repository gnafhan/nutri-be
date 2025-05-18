package model

import (
	"time"
)

// DailyNutrition represents the nutritional summary for a day
type DailyNutrition struct {
	Date     time.Time `json:"date"`
	Calories float64   `json:"calories"`
	Protein  float64   `json:"protein"`
	Carbs    float64   `json:"carbs"`
	Fat      float64   `json:"fat"`
}

// WeightHeightStatistics represents the weight and height statistics
type WeightHeightStatistics struct {
	CurrentWeight      *float64                   `json:"current_weight"`
	CurrentHeight      *float64                   `json:"current_height"`
	WeightHistory      []UsersWeightHeightHistory `json:"weight_history"`
	LatestWeightTarget *UsersWeightHeightTarget   `json:"latest_weight_target,omitempty"`
}

// HomeStatistics combines multiple statistics for the home page
type HomeStatistics struct {
	DailyNutrition         *DailyNutrition         `json:"daily_nutrition"`
	WeightHeightStatistics *WeightHeightStatistics `json:"weight_height_statistics"`
}
