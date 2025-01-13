package main

import (
	"github.com/esonhugh/k8spider/cmd"
	_ "github.com/esonhugh/k8spider/cmd/all"
	_ "github.com/esonhugh/k8spider/cmd/axfr"
	_ "github.com/esonhugh/k8spider/cmd/dnsutils"
	_ "github.com/esonhugh/k8spider/cmd/metrics"
	_ "github.com/esonhugh/k8spider/cmd/neighbor"
	_ "github.com/esonhugh/k8spider/cmd/service"
	_ "github.com/esonhugh/k8spider/cmd/whereisdns"
	_ "github.com/esonhugh/k8spider/cmd/wildcard"
)

func main() {
	cmd.Execute()
}
