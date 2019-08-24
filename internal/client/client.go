package client

type Client struct {
	ID         int    `json:"id" gorm:"PRIMARY_KEY"`
	Email      string `json:"email" gorm:"NOT NULL"`
	HSERuzID   int    `json:"hse_ruz_id" `
	GoogleCode string `json:"google_code" gorm:"UNIQUE;NOT NULL"`
}
