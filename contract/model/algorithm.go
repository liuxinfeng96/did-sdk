package model

import (
	"errors"

	bccrypto "github.com/liuxinfeng96/bc-crypto"
)

const (
	SHA256WithRSA   = "SHA256-RSA"
	ECDSAWithSHA256 = "ECDSA-SHA256"
	ECDSAWithSHA384 = "ECDSA-SHA384"
	ECDSAWithSHA512 = "ECDSA-SHA512"
	SM2WithSM3      = "SM2-SM3"
)

func GetHashType(signatureAlgo string) (bccrypto.Hash, error) {
	switch signatureAlgo {
	case SHA256WithRSA, ECDSAWithSHA256:
		return bccrypto.SHA256, nil
	case ECDSAWithSHA384:
		return bccrypto.SHA384, nil
	case ECDSAWithSHA512:
		return bccrypto.SHA512, nil
	case SM2WithSM3:
		return bccrypto.SM3, nil
	default:
		return bccrypto.Hash(0), errors.New("unknown signature algorithm")
	}
}
