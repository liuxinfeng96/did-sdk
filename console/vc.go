package main

import (
	"did-sdk/key"
	"did-sdk/vc"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
	return vcCmd
}

func vcIssueCmd() *cobra.Command {
	var algo, skPath, pkPath string
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
$ cmc vc issue \
-sk ./testdata/sk.pem \
-pk ./testdata/pk.pem \
-al SM2 \
-ki 1 \
-sub ./testdata/temp.json \
-e 2024-12-30 \
-id 1223355 \
-ti 123131231 \
-t Identity \
-vc ./testdata/vc.json \
-C ./testdata/sdk.yaml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(pkPath) == 0 {
				return ParamsEmptyError(ParamsFlagPkPath)
			}

			if len(algo) == 0 {
				return ParamsEmptyError(ParamsFlagAlgorithm)
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

			keyInfo := &key.KeyInfo{
				PkPEM:     pkPem,
				SkPEM:     skPem,
				Algorithm: algo,
			}

			sub := make(map[string]interface{})

			err = json.Unmarshal(subjectJson, &sub)
			if err != nil {
				return err
			}

			vcBytes, err := vc.IssueVC(keyInfo, keyIndex, sub, c, id, timeUnix, tid, vcType...)
			if err != nil {
				return err
			}

			err = os.WriteFile(vcPath, vcBytes, 0777)
			if err != nil {
				return err
			}

			return nil
		},
	}

	attachFlagString(vcIssueCmd, ParamsFlagPkPath, &pkPath)
	attachFlagString(vcIssueCmd, ParamsFlagSkPath, &skPath)
	attachFlagString(vcIssueCmd, ParamsFlagAlgorithm, &algo)
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
	var skPath, algo, issuer string
	var keyIndex int
	var tempPath, subjectPath, expiration, id, tid, vcPath string
	var timeUnix int64
	var vcType []string

	vcIssueCmd := &cobra.Command{
		Use:   "issue-local",
		Short: "Issue the vc",
		Long: strings.TrimSpace(
			`Issue the vc at local.
Example:
$ cmc vc issue-local \
-sk ./testdata/sk.pem \
-al SM2 \
-ki 1 \
-sub ./testdata/temp.json \
-i did:cm:admin \
-e 2024-12-30 \
-id 1223355 \
-temp ./testdata/template.json \
-t Identity \
-vc ./testdata/vc.json 
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(algo) == 0 {
				return ParamsEmptyError(ParamsFlagAlgorithm)
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

			if len(tid) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateId)
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

			vcBytes, err := vc.IssueVCLocal(skPem, algo, keyIndex, sub, issuer, id, timeUnix, temp, vcType...)
			if err != nil {
				return err
			}

			err = os.WriteFile(vcPath, vcBytes, 0777)
			if err != nil {
				return err
			}

			return nil
		},
	}

	attachFlagString(vcIssueCmd, ParamsFlagSkPath, &skPath)
	attachFlagString(vcIssueCmd, ParamsFlagAlgorithm, &algo)
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
$ cmc vc verify \
-vc ./testdata/vc.json \
-C ./testdata/sdk.yaml
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

			fmt.Printf("the verification result of vc is: [%+v]", ok)

			return nil
		},
	}

	attachFlagString(vcVerifyCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcVerifyCmd, ParamsFlagVcPath, &vcPath)

	return vcVerifyCmd
}
