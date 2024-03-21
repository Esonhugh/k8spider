package main

import (
	"github.com/esonhugh/k8spider/cmd"
	_ "github.com/esonhugh/k8spider/cmd/all"
	_ "github.com/esonhugh/k8spider/cmd/axfr"
	_ "github.com/esonhugh/k8spider/cmd/env"
	_ "github.com/esonhugh/k8spider/cmd/service"
	_ "github.com/esonhugh/k8spider/cmd/subnet"
	_ "github.com/esonhugh/k8spider/cmd/wildcard"
)

func main() {
	cmd.Execute()
}
