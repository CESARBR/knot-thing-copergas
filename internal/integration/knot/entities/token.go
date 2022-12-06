package entities

type Token struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float32 `json:"expires_in"`
	ExpiresUtc  string  `json:"expiresUtc"`
}
