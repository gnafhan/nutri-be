package response

import "app/src/model"

type PaymentResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *model.PaymentResponse `json:"data"`
}
