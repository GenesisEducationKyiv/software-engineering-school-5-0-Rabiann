package models

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherResponse struct {
	Current `json:"current"`
}

type Current struct {
	Temperature float64 `json:"temp_c"`
	Humidity    float64 `json:"humidity"`
	Condition   `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
}
