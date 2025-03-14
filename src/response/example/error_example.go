package example

type Unauthorized struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Please authenticate"`
}

type FailedLogin struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid email or password"`
}

type FailedResetPassword struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Password reset failed"`
}

type FailedVerifyEmail struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Verify email failed"`
}

type FailedVerifyProductToken struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid or already used product token"`
}

type Forbidden struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"You don't have permission to access this resource"`
}

type NotFound struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Not found"`
}

type DuplicateEmail struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Email already taken"`
}
