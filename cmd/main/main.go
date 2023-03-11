package main

import (
	"fmt"

	"github.com/roland-burke/fuel-tracker/internal/config"
	"github.com/roland-burke/fuel-tracker/internal/repository"
	"github.com/roland-burke/fuel-tracker/internal/server"
	"github.com/roland-burke/rollogger"
)

func main() {
	config.Logger = rollogger.Init(rollogger.INFO_LEVEL, true, true)
	port, urlPrefix := config.InitConfig()
	repository.InitDb()
	fmt.Printf("=======================================================\n\n")
	server.StartServer(port, urlPrefix)
}
