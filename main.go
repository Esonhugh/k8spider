package main

import (
	"github.com/esonhugh/k8spider/cmd"
	_ "github.com/esonhugh/k8spider/cmd/all"
	_ "github.com/esonhugh/k8spider/cmd/service"
	_ "github.com/esonhugh/k8spider/cmd/subnet"
)

func main() {
	cmd.Execute()
}
