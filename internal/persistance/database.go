package persistance

import (
	"fmt"
	"os"

	"github.com/Rabiann/weather-mailer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ConnectionString struct {
	User     string
	Password string
	DbName   string
	Port     string
	Host     string
}

func NewConnectionString() ConnectionString {
	cs := ConnectionString{}
	cs.User = os.Getenv("POSTGRES_USER")
	cs.Password = os.Getenv("POSTGRES_PASSWORD")
	cs.DbName = os.Getenv("POSTGRES_DB")
	cs.Host = os.Getenv("POSTGRES_HOST")

	return cs
}

func (c ConnectionString) GetConnectionString() string {
	is_prod := os.Getenv("PROD")
	if is_prod == "1" {
		db_url := os.Getenv("PROD_DB_URL")
		if db_url == "" {
			panic("PROD_DB_URL not set")
		}

		return db_url
	}

	return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s", c.Host, c.User, c.Password, c.DbName)
}

func ConnectToDatabase() *gorm.DB {
	cs := NewConnectionString()
	db, err := gorm.Open(postgres.Open(cs.GetConnectionString()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Token{}); err != nil {
		return err
	}

	return nil
}

func SetupInMemoryDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
