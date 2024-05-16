/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/admin"
	"fmt"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func AdminCMD() *cobra.Command {

	adminCmd := &cobra.Command{
		Use:   "admin",
		Short: "ChainMaker DID admin command",
		Long:  "ChainMaker DID admin command",
	}

	adminCmd.AddCommand(adminAdd())
	adminCmd.AddCommand(adminDelete())
	adminCmd.AddCommand(authAdmin())
	return adminCmd
}

func adminAdd() *cobra.Command {
	var sdkPath, adminSdkPath string

	adminAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add did admin",
		Long: strings.TrimSpace(
			`Add did contract admin to blockchain.
Example:
$ ./console admin add \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(adminSdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagAdminSdkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			adminc, err := cmsdk.NewChainClient(cmsdk.WithConfPath(adminSdkPath))
			if err != nil {
				return err
			}

			adminPk, err := adminc.GetPublicKey().String()
			if err != nil {
				return err
			}

			err = admin.SetAdminForDidContract([]byte(adminPk), c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(adminAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(adminAddCmd, ParamsFlagAdminSdkPath, &adminSdkPath)

	return adminAddCmd
}

func adminDelete() *cobra.Command {
	var sdkPath, adminSdkPath string

	adminDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete did admin",
		Long: strings.TrimSpace(
			`Delete did contract admin on blockchain.
Example:
$ ./console admin delete \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(adminSdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagAdminSdkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			adminc, err := cmsdk.NewChainClient(cmsdk.WithConfPath(adminSdkPath))
			if err != nil {
				return err
			}

			adminPk, err := adminc.GetPublicKey().String()
			if err != nil {
				return err
			}

			err = admin.DeleteAdminForDidContract([]byte(adminPk), c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(adminDeleteCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(adminDeleteCmd, ParamsFlagAdminSdkPath, &adminSdkPath)

	return adminDeleteCmd
}

func authAdmin() *cobra.Command {
	var sdkPath string

	authAdminCmd := &cobra.Command{
		Use:   "auth",
		Short: "Whether to have the admin permission",
		Long: strings.TrimSpace(
			`Whether to have the admin permission in the did contract.
Example:
$ ./console admin auth \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			pk, err := c.GetPublicKey().String()
			if err != nil {
				return err
			}

			ok, err := admin.IsAdminOfDidContract([]byte(pk), c)
			if err != nil {
				return err
			}

			fmt.Printf("Is admin: [%v]\n", ok)

			return nil
		},
	}

	attachFlagString(authAdminCmd, ParamsFlagCMSdkPath, &sdkPath)

	return authAdminCmd
}
