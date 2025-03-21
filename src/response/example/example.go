package example

import (
	"time"
)

type RegisterResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Register successfully"`
	User    User   `json:"user"`
	Tokens  Tokens `json:"tokens"`
}

type LoginResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Login successfully"`
	User    User   `json:"user"`
	Tokens  Tokens `json:"tokens"`
}

type GoogleLoginResponse struct {
	Status  string     `json:"status" example:"success"`
	Message string     `json:"message" example:"Login successfully"`
	User    GoogleUser `json:"user"`
	Tokens  Tokens     `json:"tokens"`
}

type LogoutResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Logout successfully"`
}

type RefreshTokenResponse struct {
	Status string `json:"status" example:"success"`
	Tokens Tokens `json:"tokens"`
}

type ForgotPasswordResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"A password reset link has been sent to your email address."`
}

type ResetPasswordResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Update password successfully"`
}

type SendVerificationEmailResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Please check your email for a link to verify your account"`
}

type VerifyEmailResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Verify email successfully"`
}

type VerifyProductTokenResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Verify product token successfully"`
}

type GetAllUserResponse struct {
	Status       string `json:"status" example:"success"`
	Message      string `json:"message" example:"Get all users successfully"`
	Results      []User `json:"results"`
	Page         int    `json:"page" example:"1"`
	Limit        int    `json:"limit" example:"10"`
	TotalPages   int64  `json:"total_pages" example:"1"`
	TotalResults int64  `json:"total_results" example:"1"`
}

type GetUserResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Get user successfully"`
	User    User   `json:"user"`
}

type GetMealScanDetailResponse struct {
	Status         string           `json:"status" example:"success"`
	Message        string           `json:"message" example:"Get meal's scan detail successfully"`
	MealScanDetail MealScanResponse `json:"meal_scan_detail"`
}

type GetMealResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Get meal successfully"`
	Meal    MealHistory `json:"meal"`
}

type CreateUserResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Create user successfully"`
	User    User   `json:"user"`
}

type UpdateUserResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Update user successfully"`
	User    User   `json:"user"`
}

type DeleteUserResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Delete user successfully"`
}

type AddMealResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Meal added successfully"`
	Meal    MealHistory `json:"meal"`
}

type UpdateMealResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Meal updated successfully"`
	Meal    MealHistory `json:"meal"`
}

type DeleteMealResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Meal deleted successfully"`
}

type AddMealRequest struct {
	Title    string    `json:"title" example:"Nasi Goreng"`
	MealTime time.Time `json:"meal_time" example:"2023-10-10T12:00:00Z"`
	Label    *string   `json:"label,omitempty" example:"Lunch"`
	Calories float64   `json:"calories" example:"500.0"`
	Protein  float64   `json:"protein" example:"20.0"`
	Carbs    float64   `json:"carbs" example:"60.0"`
	Fat      float64   `json:"fat" example:"15.0"`
}

type UpdateMealRequest struct {
	Title    string    `json:"title,omitempty" example:"Nasi Goreng Spesial"`
	MealTime time.Time `json:"meal_time,omitempty" example:"2023-10-10T12:30:00Z"`
	Label    *string   `json:"label,omitempty" example:"Dinner"`
	Calories float64   `json:"calories,omitempty" example:"550.0"`
	Protein  float64   `json:"protein,omitempty" example:"25.0"`
	Carbs    float64   `json:"carbs,omitempty" example:"65.0"`
	Fat      float64   `json:"fat,omitempty" example:"18.0"`
}

type NutrientDetail struct {
	Quantity float64 `json:"quantity" example:"250.5"`
	Unit     string  `json:"unit" example:"kcal"`
}

type Nutrient struct {
	Calories NutrientDetail `json:"calories"`
	Protein  NutrientDetail `json:"protein"`
	Carbs    NutrientDetail `json:"carbs"`
	Fat      NutrientDetail `json:"fat"`
}

type MealScanResponse struct {
	Status   string   `json:"status" example:"success"`
	Message  string   `json:"message" example:"Meal scanned successfully"`
	Foods    []string `json:"foods" example:"chicken,rice,salad"`
	Nutrient Nutrient `json:"nutrient"`
}

type GetAllMealsResponse struct {
	Status       string        `json:"status" example:"success"`
	Message      string        `json:"message" example:"Get all meals successfully"`
	Results      []MealHistory `json:"results"`
	Page         int           `json:"page" example:"1"`
	Limit        int           `json:"limit" example:"10"`
	TotalPages   int64         `json:"total_pages" example:"5"`
	TotalResults int64         `json:"total_results" example:"50"`
}

type MealHistory struct {
	ID       string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID   string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title    string    `json:"title" example:"Scanned Meal"`
	MealTime time.Time `json:"meal_time" example:"2023-10-01T12:00:00Z"`
	Label    *string   `json:"label,omitempty" example:"Lunch"`
	Calories float64   `json:"calories" example:"250.5"`
	Protein  float64   `json:"protein" example:"30.2"`
	Carbs    float64   `json:"carbs" example:"45.3"`
	Fat      float64   `json:"fat" example:"10.1"`
}

type UsersWeightHeightHistory struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Weight     float64   `json:"weight" example:"50.5"`
	Height     float64   `json:"height" example:"180.8"`
	RecordedAt time.Time `json:"recorded_at" example:"2023-10-01T12:00:00Z"`
}

type AddWeightHeightResponse struct {
	Status  string                   `json:"status" example:"success"`
	Message string                   `json:"message" example:"Weight and height record added successfully"`
	Data    UsersWeightHeightHistory `json:"data"`
}

type GetAllWeightHeightResponse struct {
	Status  string                     `json:"status" example:"success"`
	Message string                     `json:"message" example:"Weight and height records fetched successfully"`
	Data    []UsersWeightHeightHistory `json:"data"`
}

type GetWeightHeightResponse struct {
	Status  string                   `json:"status" example:"success"`
	Message string                   `json:"message" example:"Weight and height records fetched successfully"`
	Data    UsersWeightHeightHistory `json:"data"`
}

type UpdateWeightHeightResponse struct {
	Status  string                   `json:"status" example:"success"`
	Message string                   `json:"message" example:"Weight and height record updated successfully"`
	Data    UsersWeightHeightHistory `json:"data"`
}

type AddWeightHeightRequest struct {
	Weight     float64    `json:"weight,omitempty" example:"70.5"`
	Height     float64    `json:"height,omitempty" example:"175.0"`
	RecordedAt *time.Time `json:"recorded_at,omitempty" example:"2023-10-10T12:00:00Z"`
}

type UpdateWeightHeightRequest struct {
	Weight     float64    `json:"weight,omitempty" example:"70.5"`
	Height     float64    `json:"height,omitempty" example:"175.0"`
	RecordedAt *time.Time `json:"recorded_at,omitempty" example:"2023-10-10T12:00:00Z"`
}

type DeleteWeightHeightResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Weight and height record deleted successfully"`
}

type UserStatisticsResponse struct {
	Body struct {
		Heights  []HeightStat  `json:"heights"`
		Weights  []WeightStat  `json:"weights"`
		Calories []CalorieStat `json:"calories"`
	}
}

type HeightStat struct {
	Height     float64   `json:"height" example:"170"`
	RecordedAt time.Time `json:"recorded_at" example:"2023-10-10T08:00:00Z"`
}

type WeightStat struct {
	Weight     float64   `json:"weight" example:"50"`
	RecordedAt time.Time `json:"recorded_at" example:"2023-10-10T08:00:00Z"`
}

type CalorieStat struct {
	Calories   float64   `json:"calories" example:"500"`
	RecordedAt time.Time `json:"recorded_at" example:"2023-10-10T08:00:00Z"`
}

type AddMealScanDetailRequest struct {
	MealHistoryId string `json:"meal_history_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ApiResult     string `json:"api_result" example:"{...}"`
}

type AddMealScanDetailResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Meal scan detail added successfully"`
}
