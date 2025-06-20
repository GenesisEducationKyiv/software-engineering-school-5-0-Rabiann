package models

import "time"

type (
	Subscription struct {
		ID        uint
		Email     string `gorm:"unique" json:"email" form:"email"`
		City      string `json:"city" form:"city"`
		Frequency string `json:"period" form:"period"`
		Confirmed bool
		CreatedAt time.Time
		UpdatedAt time.Time
		Tokens    []Token
	}

	Subscriber struct {
		Recipient string
		Period    string
		City      string
	}
)
