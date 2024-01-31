package utils

import (
	"bytes"
	"encoding/json"
	"time"

	bccrypto "github.com/liuxinfeng96/bc-crypto"
)

func HashStringToHashType(h string) bccrypto.Hash {
	switch h {
	case "MD4":
		return bccrypto.MD4
	case "MD5":
		return bccrypto.MD5
	case "SHA-1":
		return bccrypto.SHA1
	case "SHA-224":
		return bccrypto.SHA224
	case "SHA-256":
		return bccrypto.SHA256
	case "SHA-384":
		return bccrypto.SHA384
	case "SHA-512":
		return bccrypto.SHA512
	case "MD5+SHA1":
		return bccrypto.MD5SHA1
	case "RIPEMD-160":
		return bccrypto.RIPEMD160
	case "SHA3-224":
		return bccrypto.SHA3_224
	case "SHA3-256":
		return bccrypto.SHA3_256
	case "SHA3-384":
		return bccrypto.SHA3_384
	case "SHA3-512":
		return bccrypto.SHA3_512
	case "SHA-512/224":
		return bccrypto.SHA512_224
	case "SHA-512/256":
		return bccrypto.SHA512 / 256
	case "BLAKE2s-256":
		return bccrypto.BLAKE2s_256
	case "BLAKE2b-256":
		return bccrypto.BLAKE2b_256
	case "BLAKE2b-384":
		return bccrypto.BLAKE2b_384
	case "BLAKE2b-512":
		return bccrypto.BLAKE2b_512
	case "SM3":
		return bccrypto.SM3
	default:
		return bccrypto.Hash(0)
	}
}

func ISO8601Time(t int64) string {
	unix := time.Unix(t, 0)
	return unix.Format(time.RFC3339)
}

func ISO8601TimeToUnix(t string) (int64, error) {
	ti, err := time.ParseInLocation(time.RFC3339, t, time.Local)
	if err != nil {
		return 0, err
	}

	return ti.Unix(), nil
}

func GetHashTypeByAlgorithm(algo string) string {
	var hash string
	if algo == "SM2" {
		hash = "SM3"
	} else {
		hash = "SHA-256"
	}
	return hash
}

// CompactJson 压缩json字符串，去掉空格换行等
func CompactJson(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, raw)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
