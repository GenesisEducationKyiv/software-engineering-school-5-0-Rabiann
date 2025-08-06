package persistance

import (
	"fmt"

	localConf "github.com/Rabiann/weather-mailer/pkg/go/user/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/user/internal/models"
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

func NewConnectionString(configuration *localConf.Configuration) ConnectionString {
	cs := ConnectionString{}
	cs.User = configuration.PostgresUser
	cs.Password = configuration.PostgresPassword
	cs.DbName = configuration.PostgresDb
	cs.Host = configuration.PostgresHost

	return cs
}

func (c ConnectionString) GetConnectionString(configuration *localConf.Configuration) string {
	is_prod := configuration.Prod
	if is_prod == "1" {
		return configuration.ProdDbUrl
	}

	return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s", c.Host, c.User, c.Password, c.DbName)
}

func ConnectToDatabase(configuration *localConf.Configuration) *gorm.DB {
	cs := NewConnectionString(configuration)
	db, err := gorm.Open(postgres.Open(cs.GetConnectionString(configuration)), &gorm.Config{})
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
