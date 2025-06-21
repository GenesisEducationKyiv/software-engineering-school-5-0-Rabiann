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

func setup() *gorm.DB {
	db, err := persistance.SetupInMemoryDb()
	if err != nil {
		panic(err)
	}

	return db
}

func TestGetSubscriptions(t *testing.T) {
	db := setup()
	model := models.Subscription{
		ID:    1,
		Email: "abc@mail.com",
	}

	db.Create(&model)

	repo := persistance.NewSubscriptionRepository(db)

	subs, err := repo.GetSubscriptions(context.TODO(), func() {})
	assert.NoError(t, err)

	require.Equal(t, 1, len(subs))
}

func TestGetSubscriptionById(t *testing.T) {
	db := setup()
	model := models.Subscription{
		ID:    1,
		Email: "abc@mail.com",
	}

	db.Create(&model)

	repo := persistance.NewSubscriptionRepository(db)

	sub, err := repo.GetSubscriptionById(model.ID, context.TODO(), func() {})
	assert.NoError(t, err)

	require.NotNil(t, sub)
	require.Equal(t, model.Email, sub.Email)
}

func TestAddSubscription(t *testing.T) {
	db := setup()
	model := models.Subscription{
		ID:    1,
		Email: "abc@mail.com",
	}

	repo := persistance.NewSubscriptionRepository(db)

	id, err := repo.AddSubscription(model, context.TODO(), func() {})
	assert.NoError(t, err)

	sub, err := repo.GetSubscriptionById(id, context.TODO(), func() {})
	assert.NoError(t, err)

	require.NotNil(t, sub)
	require.Equal(t, model.Email, sub.Email)
}

func TestActivateSubscription(t *testing.T) {
	db := setup()
	model := models.Subscription{
		ID:        1,
		Email:     "abc@mail.com",
		Confirmed: false,
	}

	db.Create(&model)

	repo := persistance.NewSubscriptionRepository(db)

	email, err := repo.ActivateSubscription(model.ID, context.TODO(), func() {})
	assert.NoError(t, err)

	require.NotNil(t, email)
	require.Equal(t, model.Email, email)

	sub, err := repo.GetSubscriptionById(model.ID, context.TODO(), func() {})
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	require.True(t, sub.Confirmed)
}

func TestGetActiveSubscriptions(t *testing.T) {
	db := setup()
	model1 := models.Subscription{
		ID:        1,
		Email:     "abcd@mail.com",
		Confirmed: true,
		Frequency: "daily",
	}

	model2 := models.Subscription{
		ID:        2,
		Email:     "abc@mail.com",
		Confirmed: false,
		Frequency: "daily",
	}

	db.Create(&model1)
	db.Create(&model2)

	repo := persistance.NewSubscriptionRepository(db)

	subs, err := repo.GetActiveSubscriptions("daily", context.TODO(), func() {})
	assert.NoError(t, err)

	require.Equal(t, 1, len(subs))
}

func TestUpdateSubscription(t *testing.T) {
	db := setup()
	model_old := models.Subscription{
		ID:        1,
		Email:     "abc@mail.com",
		Confirmed: false,
	}

	model_new := models.Subscription{
		ID:        1,
		Email:     "def@mail.com",
		Confirmed: true,
	}

	db.Create(&model_old)

	repo := persistance.NewSubscriptionRepository(db)

	err := repo.UpdateSubscription(model_old.ID, model_new, context.TODO(), func() {})
	assert.NoError(t, err)
	sub, err := repo.GetSubscriptionById(model_old.ID, context.TODO(), func() {})
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	require.Equal(t, model_new.Email, sub.Email)
}

func TestDeleteSubscription(t *testing.T) {
	db := setup()
	model := models.Subscription{
		ID:    1,
		Email: "abc@mail.com",
	}

	db.Create(&model)

	repo := persistance.NewSubscriptionRepository(db)

	err := repo.DeleteSubscription(model.ID, context.TODO(), func() {})
	assert.NoError(t, err)

	_, err = repo.GetSubscriptionById(model.ID, context.TODO(), func() {})
	assert.Error(t, err)
}
