package config_test

import (
	"testing"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCamelCaseToUpperCase(t *testing.T) {
	src := "CamelCase"
	res := "CAMEL_CASE"

	require.Equal(t, res, config.FromCamelCaseToUpperCase(src))
}

func TestBuildConfig(t *testing.T) {
	_, err := config.LoadEnvironment()
	require.NoError(t, err)
}
