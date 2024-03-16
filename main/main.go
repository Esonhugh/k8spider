package main

import (
	"github.com/esonhugh/go-cli-template-v2/cmd"
	_ "github.com/esonhugh/go-cli-template-v2/cmd/version"
	"github.com/esonhugh/go-cli-template-v2/internal/config"
	"github.com/esonhugh/go-cli-template-v2/internal/database"
	"github.com/esonhugh/go-cli-template-v2/internal/log"
)

func init() {
	log.Init("info")
	config.Init("config.yaml")
	_ = database.Init("db.sqlite3")
}

func main() {
	cmd.Execute()
}
