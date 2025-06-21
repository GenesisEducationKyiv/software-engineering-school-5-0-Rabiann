package services

import (
	"os"
	"strings"
)

type (
	Template struct {
		Text string
	}

	ConfirmationTemplate struct {
		Template *Template
	}

	WeatherTemplate struct {
		Template *Template
	}
)

func NewTemplate(filepath string) (*Template, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var template Template

	template.Text = string(text)
	return &template, nil
}

func NewConfirmationTemplate(filepath string) (*ConfirmationTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &ConfirmationTemplate{Template: template}, nil
}

func NewWeatherTemplate(filepath string) (*WeatherTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &WeatherTemplate{Template: template}, nil
}

func (ct *ConfirmationTemplate) BuildConfirmationLetter(url string) string {
	return strings.Replace(ct.Template.Text, "{}", url, 3)
}

func (wt *WeatherTemplate) BuildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	let := strings.Replace(wt.Template.Text, "{City}", city, 1)
	let = strings.Replace(let, "{Temperature}", temp, 1)
	let = strings.Replace(let, "{Humidity}", humid, 1)
	let = strings.Replace(let, "{UnsubscribeLink}", unsubscribe, 1)
	let = strings.Replace(let, "{Description}", description, 1)
	return let
}
