package main

import (

	"github.com/ThienKim52/golang-dev/internal/infrastructure"

)

// @title Health Check API
// @version 1.1
// @description This is a health check API server.
// @Basepath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	engine := infrastructure.CreateAPI()
	if err := engine.Start(); err != nil {
        panic(err)
    }
}
