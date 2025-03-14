package response

import "time"

type Tokens struct {
	Access  TokenExpires `json:"access"`
	Refresh TokenExpires `json:"refresh"`
}

type TokenExpires struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type RefreshToken struct {
	Status string `json:"status"`
	Tokens Tokens `json:"tokens"`
}
