package models

type Weather struct {
	Temperature float64 `json:"temperature" redis:"temperature"`
	Humidity    float64 `json:"humidity" redis:"humidity"`
	Description string  `json:"description" redis:"description"`
}
