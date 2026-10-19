package fixtures

import (
	"gorm.io/gorm"
	"github.com/ThienKim52/golang-dev/internal/app/model"

)

type UserCommonTestDB struct {
	base
}

func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})
	users := []*model.User{
		{
			ID: "deb745af-1a62-4efa-99a0-f06b274bd990",
			DisplayName: "John Doe",
			Username: "johndoe",
			Password: "password123",
			Email: "johndoe@example.com",
		},
		{
			ID: "deb745af-1a62-4efa-99a0-f06b274bd991",
			DisplayName: "Aimee Than",
			Username: "aimeethan",
			Password: "password123",
			Email: "aimeethan@example.com",
		},
	}
	return db.CreateInBatches(users, 10).Error
}








