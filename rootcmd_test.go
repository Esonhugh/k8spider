package main

import (
	"strings"
	"testing"

	"github.com/esonhugh/k8spider/cmd"
)

func TestGetCommandTree(t *testing.T) {
	t.Logf("TestGetCommandTree")
	treeRoot := cmd.RootCmd.Root()
	t.Logf("treeRoot: %v", treeRoot)
	for _, c := range treeRoot.Commands() {
		t.Logf("test: $(BUILD_DIR)/$(MAIN_PROGRAM_NAME) %v --help", strings.ReplaceAll(c.CommandPath(), "k8spider ", ""))
		if c.HasSubCommands() {
			for _, sc := range c.Commands() {
				t.Logf("test: $(BUILD_DIR)/$(MAIN_PROGRAM_NAME) %v --help", strings.ReplaceAll(sc.CommandPath(), "k8spider ", ""))
			}
		}
	}
}
