//go:build tools
// +build tools

package main

import (
	_ "github.com/spf13/cobra/cobra"
	_ "golang.org/x/tools/cmd/goimports"
)
