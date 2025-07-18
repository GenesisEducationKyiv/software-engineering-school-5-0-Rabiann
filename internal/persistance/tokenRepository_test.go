package persistance_test

import (
	"context"
	"testing"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/persistance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTokenDB() *gorm.DB {
	db, err := persistance.SetupInMemoryDb()
	if err != nil {
		panic(err)
	}

	return db
}

func TestCreateToken(t *testing.T) {
	db := setupTokenDB()
	var id uint = 1

	repo := persistance.NewTokenRepository(db)

	token, err := repo.CreateToken(id, context.TODO())
	assert.NoError(t, err)

	tokenModel := models.Token{}

	db.First(&tokenModel)
	require.Equal(t, tokenModel.ID, token)
}

func TestGetSubscriptionOfToken(t *testing.T) {
	db := setupTokenDB()
	var id uint = 1

	repo := persistance.NewTokenRepository(db)

	token, err := repo.CreateToken(id, context.TODO())
	assert.NoError(t, err)

	mailId, err := repo.GetSubscriptionOfToken(token, context.TODO())
	assert.NoError(t, err)

	require.Equal(t, id, mailId)
}

func TestUseToken(t *testing.T) {
	db := setupTokenDB()
	var id uint = 1

	repo := persistance.NewTokenRepository(db)

	token, err := repo.CreateToken(id, context.TODO())
	assert.NoError(t, err)

	err = repo.UseToken(token, context.TODO())
	assert.NoError(t, err)

	id, err = repo.GetSubscriptionOfToken(token, context.TODO())
	assert.NoError(t, err)

	require.Equal(t, uint(0), id)
}
