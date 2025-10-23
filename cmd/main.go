package main

import (
	"github.com/TNAHOM/ATS-system-main/initiator"
)

// @title		 	ATS System API
// @version         1.0
// @description     REST API for ATS-system-main (users, auth, job posts)

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	initiator.Initiator()
}
