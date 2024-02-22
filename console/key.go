package main

import (
	"did-sdk/key"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func KeyCMD() *cobra.Command {

	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "ChainMaker DID key command",
		Long:  "ChainMaker DID key command",
	}

	keyCmd.AddCommand(keyGenCMD())
	return keyCmd
}

func keyGenCMD() *cobra.Command {
	var algo, skPath, pkPath string

	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Private key generate",
		Long: strings.TrimSpace(
			`Generate the private key of the specified crypto algorithm.
Supported algorithms: SM2, EC_Secp256k1, EC_NISTP224, EC_NISTP256, EC_NISTP384, EC_NISTP521, RSA2048, RSA3072 .
Example:
$ cmc key gen -al EC_Secp256k1 -pk ./pk.pem -sk ./sk.pem
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(pkPath) == 0 {
				return ParamsEmptyError(ParamsFlagPkPath)
			}

			if len(pkPath) == 0 {
				return ParamsEmptyError(ParamsFlagAlgorithm)
			}

			return keyGen(algo, skPath, pkPath)
		},
	}

	attachFlagString(genCmd, ParamsFlagAlgorithm, &algo)
	attachFlagString(genCmd, ParamsFlagPkPath, &pkPath)
	attachFlagString(genCmd, ParamsFlagSkPath, &skPath)

	return genCmd
}

func keyGen(algo, skPath, pkPath string) error {
	keyInfo, err := key.GenerateKey(algo)
	if err != nil {
		return err
	}

	err = os.WriteFile(skPath, keyInfo.SkPEM, 0777)
	if err != nil {
		return err
	}

	err = os.WriteFile(pkPath, keyInfo.PkPEM, 0777)
	if err != nil {
		return err
	}

	return nil
}
