package dto

type (
	SubscriptionRequest struct {
		Email     string `gorm:"unique" json:"email" form:"email"`
		City      string `json:"city" form:"city"`
		Frequency string `json:"period" form:"period"`
	}

	Weather struct {
		Temperature float64 `json:"temperature" redis:"temperature"`
		Humidity    float64 `json:"humidity" redis:"humidity"`
		Description string  `json:"description" redis:"description"`
	}

	SubscriptionConfirmationMsg struct {
		Email string `json:"email"`
		Url   string `json:"url"`
	}

	Subscription struct {
		Id     uint
		Email  string
		City   string
		Period string
	}

	Subscriber struct {
		Recipient string `json:"recipient"`
		Period    string `json:"period"`
		City      string `json:"city"`
	}

	ReportMsg struct {
		Subscriber        Subscriber `json:"subscriber"`
		Weather           Weather    `json:"weather"`
		UnsubscriptionUrl string     `json:"unsubscriptionUrl"`
	}
)
