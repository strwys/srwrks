package main

import (
	"github.com/cecepsprd/starworks-test/cmd"
	_ "github.com/cecepsprd/starworks-test/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Starworks server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /v3
func main() {
	cmd.Execute()
}
