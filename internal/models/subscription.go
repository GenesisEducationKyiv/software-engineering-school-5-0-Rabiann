package models

import "time"

type (
	Subscription struct {
		ID        uint
		Email     string `gorm:"unique" form:"email"`
		City      string `form:"city"`
		Frequency string `form:"period"`
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
