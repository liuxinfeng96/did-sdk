/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"errors"

	bccrypto "github.com/liuxinfeng96/bc-crypto"
)

const (
	// SHA256WithRSA RSA2048+SHA256，RSA3072+SHA256 signature algorithm
	SHA256WithRSA = "SHA256-RSA"
	// ECDSAWithSHA256 EC_Secp256k1、EC_NISTP224、EC_NISTP256+SHA256 signature algorithm
	ECDSAWithSHA256 = "ECDSA-SHA256"
	// ECDSAWithSHA384 EC_NISTP384+SHA384 signature algorithm
	ECDSAWithSHA384 = "ECDSA-SHA384"
	// ECDSAWithSHA512 EC_NISTP521+SHA512 signature algorithm
	ECDSAWithSHA512 = "ECDSA-SHA512"
	// SM2WithSM3 SM2+SM3 signature algorithm
	SM2WithSM3 = "SM2-SM3"
)

// GetHashType infer the hash algorithm from the signature algorithm
// @params signatureAlgo 签名算法的名称，一般从证明结构的`type`字段获取
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
