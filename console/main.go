package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{
		Use:   "didc",
		Short: "ChainMaker DID CLI",
		Long:  strings.TrimSpace(`Command line interface to support ChainMaker distributed digital identity`),
	}

	mainCmd.AddCommand(KeyCMD())
	mainCmd.AddCommand(DidCMD())
	mainCmd.AddCommand(BlackCMD())

	mainCmd.Execute()
}
