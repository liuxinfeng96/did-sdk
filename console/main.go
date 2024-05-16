/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{
		Use:   "console",
		Short: "ChainMaker DID CLI",
		Long:  strings.TrimSpace(`Command line interface to support ChainMaker distributed digital identity`),
	}

	mainCmd.AddCommand(KeyCMD())
	mainCmd.AddCommand(DidCMD())
	mainCmd.AddCommand(BlackCMD())
	mainCmd.AddCommand(DocCMD())
	mainCmd.AddCommand(IssuerCMD())
	mainCmd.AddCommand(VcRevokeCMD())
	mainCmd.AddCommand(VcTemplateCMD())
	mainCmd.AddCommand(VcTemplateCMD())
	mainCmd.AddCommand(VcCMD())
	mainCmd.AddCommand(VpCMD())
	mainCmd.AddCommand(AdminCMD())

	err := mainCmd.Execute()
	if err != nil {
		panic(err)
	}

}
