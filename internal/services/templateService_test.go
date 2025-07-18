package services_test

import (
	"os"
	"testing"

	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "template-*.txt")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)

	err = tmpFile.Close()
	require.NoError(t, err, "Failed to close tmp file")

	return tmpFile.Name()
}

func TestNewTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedContent := "TEST"
		filePath := setup(t, expectedContent)

		template, err := services.NewTemplate(filePath)

		assert.NoError(t, err)
		require.NotNil(t, template)
		assert.Equal(t, expectedContent, template.Text)
	})

	t.Run("File Not Found", func(t *testing.T) {
		template, err := services.NewTemplate("not-exists.txt")

		assert.Error(t, err)
		assert.Nil(t, template)
	})
}

func TestNewConfirmationTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		content := "A {} A"
		filePath := setup(t, content)

		confTemplate, err := services.NewConfirmationTemplate(filePath)

		assert.NoError(t, err)
		require.NotNil(t, confTemplate)
		require.NotNil(t, confTemplate.Template)
		assert.Equal(t, content, confTemplate.Template.Text)
	})
	t.Run("File Not Found", func(t *testing.T) {
		confTemplate, err := services.NewConfirmationTemplate("non-existent-file.txt")
		assert.Error(t, err)
		assert.Nil(t, confTemplate)
	})
}

func TestNewWeatherTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		content := "{City} {Temperature}"
		filePath := setup(t, content)

		weatherTemplate, err := services.NewWeatherTemplate(filePath)

		assert.NoError(t, err)
		require.NotNil(t, weatherTemplate)
		require.NotNil(t, weatherTemplate.Template)
	})
	t.Run("File Not Found", func(t *testing.T) {
		weatherTemplate, err := services.NewWeatherTemplate("non-existing-file.txt")
		assert.Error(t, err)
		require.Nil(t, weatherTemplate)
	})
}

func TestBuildConfirmationLetter(t *testing.T) {
	content := "A {}"
	filePath := setup(t, content)

	confirmationTemplate, err := services.NewConfirmationTemplate(filePath)
	assert.NoError(t, err)

	url := "https://example.com/confirm/12345"

	result := confirmationTemplate.BuildConfirmationLetter(url)

	expected := "A https://example.com/confirm/12345"
	assert.Equal(t, result, expected)
}

func TestBuildWeatherLetter(t *testing.T) {
	content := "{City} {Temperature} {Humidity} {UnsubscribeLink} {Description}"
	filePath := setup(t, content)
	weatherTemplate, err := services.NewWeatherTemplate(filePath)

	city := "City"
	temp := "1"
	humid := "1"
	unsubscribe := "123.com"
	description := "abc"

	result := weatherTemplate.BuildWeatherLetter(city, temp, humid, description, unsubscribe)
	expected := "City 1 1 123.com abc"

	assert.NoError(t, err)
	require.Equal(t, result, expected)
}
