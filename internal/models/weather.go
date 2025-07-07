package models

type Weather struct {
	Temperature float64 `json:"temperature" redis:"temperature"`
	Humidity    float64 `json:"humidity" redis:"humidity"`
	Description string  `json:"description" redis:"description"`
}

type WeatherApiResponse struct {
	Current struct {
		Temperature float64 `json:"temp_c"`
		Humidity    float64 `json:"humidity"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type WeatherMapResponse struct {
	Weather struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temperature float64 `json:"temp"`
		Humidity    float64 `json:"humidity"`
	} `json:"main"`
}

type WeatherStackResponse struct {
	Current struct {
		Temperature float64 `json:"temperature"`
		Description string  `json:"weather_description"`
		Humidity    float64 `json:"humidity"`
	} `json:"current"`
}
