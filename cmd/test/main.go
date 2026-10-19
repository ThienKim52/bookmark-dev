package main

import (
	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/google/uuid"
)

func main() {
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	dbClient.AutoMigrate(&model.User{})
	dbClient.Create(&model.User{
		ID: uuid.New().String(),
		Username: "test",
		Password: "test",
	})
}