package services

import (
	"os"
	"strings"
)

type (
	Template struct {
		text string
	}

	ConfirmationTemplate struct {
		template *Template
	}

	WeatherTemplate struct {
		template *Template
	}
)

func NewTemplate(filepath string) (*Template, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var template Template

	template.text = string(text)
	return &template, nil
}

func NewConfirmationTemplate(filepath string) (*ConfirmationTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &ConfirmationTemplate{template: template}, nil
}

func NewWeatherTemplate(filepath string) (*WeatherTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &WeatherTemplate{template: template}, nil
}

func (ct *ConfirmationTemplate) buildConfirmationLetter(email string) string {
	return strings.Replace(ct.template.text, "{}", email, 3)
}

func (wt *WeatherTemplate) buildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	let := strings.Replace(wt.template.text, "{City}", city, 1)
	let = strings.Replace(let, "{Temperature}", temp, 1)
	let = strings.Replace(let, "{Humidity}", humid, 1)
	let = strings.Replace(let, "{UnsubscribeLink}", unsubscribe, 1)
	let = strings.Replace(let, "{Description}", description, 1)
	return let
}
