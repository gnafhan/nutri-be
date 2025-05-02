package response

import (
	"app/src/model"
	"time"
)

type Common struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessWithUser struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

type SuccessWithMeal struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Meal    model.MealHistory `json:"meal_history"`
}
type SuccessWithMealScanDetail struct {
	Status         string                  `json:"status"`
	Message        string                  `json:"message"`
	MealScanDetail model.MealHistoryDetail `json:"meal_history_scan_detail"`
}

type SuccessWithWeightHeight struct {
	Status  string                         `json:"status" example:"success"`
	Message string                         `json:"message" example:"Operation completed successfully"`
	Data    model.UsersWeightHeightHistory `json:"data"`
}

type SuccessWithWeightHeightTarget struct {
	Status  string                        `json:"status" example:"success"`
	Message string                        `json:"message" example:"Operation completed successfully"`
	Data    model.UsersWeightHeightTarget `json:"data"`
}

type SuccessWithWeightHeightList struct {
	Status  string                           `json:"status" example:"success"`
	Message string                           `json:"message" example:"Operation completed successfully"`
	Data    []model.UsersWeightHeightHistory `json:"data"`
}

type SuccessWithWeightHeightTargetList struct {
	Status  string                          `json:"status" example:"success"`
	Message string                          `json:"message" example:"Operation completed successfully"`
	Data    []model.UsersWeightHeightTarget `json:"data"`
}

type SuccessWithTokens struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
	Tokens  Tokens     `json:"tokens"`
}

type SuccessWithPaginate[T any] struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	Results      []T    `json:"results"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TotalPages   int64  `json:"total_pages"`
	TotalResults int64  `json:"total_results"`
}

type ErrorDetails struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type UserStatistics struct {
	Heights  []HeightStat  `json:"heights"`
	Weights  []WeightStat  `json:"weights"`
	Calories []CalorieStat `json:"calories"`
}

type HeightStat struct {
	Height     float64   `json:"height"`
	RecordedAt time.Time `json:"recorded_at"`
}

type WeightStat struct {
	Weight     float64   `json:"weight"`
	RecordedAt time.Time `json:"recorded_at"`
}

type CalorieStat struct {
	Calories   float64   `json:"calories"`
	RecordedAt time.Time `json:"recorded_at"`
}

type SuccessWithArticle struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    model.ArticleResponse `json:"data"`
}

type SuccessWithArticleList struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    []model.ArticleResponse `json:"data"`
}

type SuccessWithArticleCategory struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    model.ArticleCategory `json:"data"`
}

type SuccessWithArticleCategoryList struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    []model.ArticleCategory `json:"data"`
}

type SuccessWithRecipe struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    model.Recipe `json:"data"`
}

type SuccessWithRecipeList struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []model.Recipe `json:"data"`
}

type Recipe struct {
	Day string `json:"day" example:"monday"`
}

type SubscriptionPlansResponse struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []model.SubscriptionPlanResponse `json:"data"`
}

type UserSubscriptionResponse struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    model.UserSubscriptionResponse `json:"data"`
}

type FeatureAccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    FeatureData `json:"data"`
}

type FeatureData struct {
	Feature string `json:"feature"`
	Access  bool   `json:"access"`
}

type CommonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
