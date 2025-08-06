package dto

type (
	Subscriber struct {
		Recipient string `json:"recipient"`
		Period    string `json:"period"`
		City      string `json:"city"`
	}

	Weather struct {
		Temperature float32 `json:"temperature"`
		Humidity    float32 `json:"humidity"`
		Description string  `json:"description"`
	}

	SubscriptionConfirmationMsg struct {
		Email string `json:"email"`
		Url   string `json:"url"`
	}

	ReportMsg struct {
		Subscriber        Subscriber `json:"subscriber"`
		Weather           Weather    `json:"weather"`
		UnsubscriptionUrl string     `json:"unsubscriptionUrl"`
	}
)
