package models

type PayloadModel struct {
	ID        string   `json:"id"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"`
	Exp       float64  `json:"exp"`
}
