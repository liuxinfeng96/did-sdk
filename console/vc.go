/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/vc"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"chainmaker.org/chainmaker/did-contract/model"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func VcCMD() *cobra.Command {

	vcCmd := &cobra.Command{
		Use:   "vc",
		Short: "ChainMaker DID vc command",
		Long:  "ChainMaker DID vc command",
	}

	vcCmd.AddCommand(vcIssueCmd())
	vcCmd.AddCommand(vcIssueLocalCmd())
	vcCmd.AddCommand(vcVerifyCmd())
	vcCmd.AddCommand(vcLogCmd())
	return vcCmd
}

func vcIssueCmd() *cobra.Command {
	var skPath, pkPath string
	var keyIndex int
	var subjectPath, expiration, id, tid, vcPath, sdkPath string
	var timeUnix int64
	var vcType []string

	vcIssueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue the vc",
		Long: strings.TrimSpace(
			`Issue the vc on the blockchain.
Example:
$ ./console vc issue \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--subject=./testdata/subject.json \
--expiration=2025-01-25 \
--id=vc001 \
--temp-id=12313213 \
--type=Identity \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(pkPath) == 0 {
				return ParamsEmptyError(ParamsFlagPkPath)
			}

			if len(subjectPath) == 0 {
				return ParamsEmptyError(ParamsFlagSubjectPath)
			}

			if len(id) == 0 {
				return ParamsEmptyError(ParamsFlagId)
			}

			if len(tid) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateId)
			}

			if len(vcPath) == 0 {
				return ParamsEmptyError(ParamsFlagVcPath)
			}

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(expiration) != 0 {
				t, err := time.ParseInLocation("2006-01-02", expiration, time.Local)
				if err != nil {
					return err
				}

				timeUnix = t.Unix()
			} else {
				// 默认1年
				timeUnix = time.Now().Add(time.Hour * 24 * 365).Unix()
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			skPem, err := os.ReadFile(skPath)
			if err != nil {
				return err
			}

			pkPem, err := os.ReadFile(pkPath)
			if err != nil {
				return err
			}

			subjectJson, err := os.ReadFile(subjectPath)
			if err != nil {
				return err
			}

			sub := make(map[string]interface{})

			err = json.Unmarshal(subjectJson, &sub)
			if err != nil {
				return err
			}

			vcBytes, err := vc.IssueVC(skPem, pkPem, keyIndex, sub, c, id, timeUnix, tid, vcType...)
			if err != nil {
				return err
			}

			err = os.WriteFile(vcPath, vcBytes, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcIssueCmd, ParamsFlagPkPath, &pkPath)
	attachFlagString(vcIssueCmd, ParamsFlagSkPath, &skPath)
	attachFlagString(vcIssueCmd, ParamsFlagCMSdkPath, &sdkPath)

	attachFlagString(vcIssueCmd, ParamsFlagSubjectPath, &subjectPath)
	attachFlagString(vcIssueCmd, ParamsFlagExpiration, &expiration)
	attachFlagString(vcIssueCmd, ParamsFlagId, &id)
	attachFlagString(vcIssueCmd, ParamsFlagTemplateId, &tid)
	attachFlagString(vcIssueCmd, ParamsFlagVcPath, &vcPath)

	attachFlagStringSlice(vcIssueCmd, ParamsFlagType, &vcType)

	attachFlagInt(vcIssueCmd, ParamsFlagKeyIndex, &keyIndex)

	return vcIssueCmd
}

func vcIssueLocalCmd() *cobra.Command {
	var skPath, issuer string
	var keyIndex int
	var tempPath, subjectPath, expiration, id, vcPath string
	var timeUnix int64
	var vcType []string

	vcIssueCmd := &cobra.Command{
		Use:   "issue-local",
		Short: "Issue the vc",
		Long: strings.TrimSpace(
			`Issue the vc at local.
Example:
$ ./console vc issue-local \
--sk-path=./testdata/sk.pem \
--subject=./testdata/subject.json \
--issuer=did:cm:admin \
--expiration=2025-01-25 \
--id=vc001 \
--temp-path=./testdata/temp.json \
--type=Identity \
--vc-path=./testdata/vc.json 
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(issuer) == 0 {
				return ParamsEmptyError(ParamsFlagIssuer)
			}

			if len(tempPath) == 0 {
				return ParamsEmptyError(ParamsFlagTemplatePath)
			}

			if len(subjectPath) == 0 {
				return ParamsEmptyError(ParamsFlagSubjectPath)
			}

			if len(id) == 0 {
				return ParamsEmptyError(ParamsFlagId)
			}

			if len(vcPath) == 0 {
				return ParamsEmptyError(ParamsFlagVcPath)
			}

			if len(expiration) != 0 {
				t, err := time.ParseInLocation("2006-01-02", expiration, time.Local)
				if err != nil {
					return err
				}

				timeUnix = t.Unix()
			} else {
				// 默认1年
				timeUnix = time.Now().Add(time.Hour * 24 * 365).Unix()
			}

			skPem, err := os.ReadFile(skPath)
			if err != nil {
				return err
			}

			subjectJson, err := os.ReadFile(subjectPath)
			if err != nil {
				return err
			}

			sub := make(map[string]interface{})

			err = json.Unmarshal(subjectJson, &sub)
			if err != nil {
				return err
			}

			temp, err := os.ReadFile(tempPath)
			if err != nil {
				return err
			}

			var vcTemp model.VcTemplate
			err = json.Unmarshal(temp, &vcTemp)
			if err != nil {
				return err
			}

			vcBytes, err := vc.IssueVCLocal(skPem, keyIndex, sub, issuer, id, timeUnix, []byte(vcTemp.Template), vcType...)
			if err != nil {
				return err
			}

			err = os.WriteFile(vcPath, vcBytes, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcIssueCmd, ParamsFlagSkPath, &skPath)
	attachFlagString(vcIssueCmd, ParamsFlagIssuer, &issuer)

	attachFlagString(vcIssueCmd, ParamsFlagSubjectPath, &subjectPath)
	attachFlagString(vcIssueCmd, ParamsFlagExpiration, &expiration)
	attachFlagString(vcIssueCmd, ParamsFlagId, &id)
	attachFlagString(vcIssueCmd, ParamsFlagTemplatePath, &tempPath)
	attachFlagString(vcIssueCmd, ParamsFlagVcPath, &vcPath)

	attachFlagStringSlice(vcIssueCmd, ParamsFlagType, &vcType)

	attachFlagInt(vcIssueCmd, ParamsFlagKeyIndex, &keyIndex)

	return vcIssueCmd
}

func vcVerifyCmd() *cobra.Command {
	var sdkPath, vcPath string

	vcVerifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the vc",
		Long: strings.TrimSpace(
			`Verify the vc on blockchain.
Example:
$ ./console vc verify \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(vcPath) == 0 {
				return ParamsEmptyError(ParamsFlagVcPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			vcJson, err := os.ReadFile(vcPath)
			if err != nil {
				return err
			}

			ok, err := vc.VerifyVCOnChain(string(vcJson), c)
			if err != nil {
				return err
			}

			fmt.Printf("the verification result of vc is: [%+v]\n", ok)

			return nil
		},
	}

	attachFlagString(vcVerifyCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcVerifyCmd, ParamsFlagVcPath, &vcPath)

	return vcVerifyCmd
}

func vcLogCmd() *cobra.Command {
	var start, count int
	var search, sdkPath string

	vcLogCmd := &cobra.Command{
		Use:   "log",
		Short: "Get vc issue log list",
		Long: strings.TrimSpace(
			`Get the vc issue log list on chain.
Example:
$ ./console vc log \
--start=1 \
--count=10 \
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

			list, err := vc.GetVcIssueLogListFromChain(search, start, count, c)
			if err != nil {
				return err
			}

			for _, v := range list {
				fmt.Printf("%+v\n", v)
			}

			return nil
		},
	}

	attachFlagString(vcLogCmd, ParamsFlagListSearch, &search)
	attachFlagString(vcLogCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagInt(vcLogCmd, ParamsFlagListStart, &start)
	attachFlagInt(vcLogCmd, ParamsFlagListCount, &count)

	return vcLogCmd
}
