package response

import "app/src/model"

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
