package proof

import (
	"did-sdk/key"
	"testing"

	"github.com/test-go/testify/require"
)

func TestGenerateProofByKey(t *testing.T) {
	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	msg := []byte("test_data")

	proof, err := GenerateProofByKey(keyInfo.SkPEM, msg, "did:cmid:gongan1234#keys-1", "EC_Secp256k1", "SHA-256")
	require.Nil(t, err)

	println(proof)

	keyInfo, err = key.GenerateKey("SM2")
	require.Nil(t, err)

	msg = []byte("test_data")

	proof, err = GenerateProofByKey(keyInfo.SkPEM, msg, "did:cmid:gongan1234#keys-1", "SM2", "SM3")
	require.Nil(t, err)

	println(proof)

}

func TestVerifyPKProof(t *testing.T) {
	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	msg := []byte("test_data")

	proof, err := GenerateProofByKey(keyInfo.SkPEM, msg, "did:cmid:gongan1234#keys-1", "EC_Secp256k1", "SHA-256")
	require.Nil(t, err)

	ok, err := VerifyPKProof(msg, keyInfo.PkPEM, proof)
	require.Nil(t, err)
	require.Equal(t, true, ok)

	keyInfo, err = key.GenerateKey("SM2")
	require.Nil(t, err)

	msg = []byte("test_data")

	proof, err = GenerateProofByKey(keyInfo.SkPEM, msg, "did:cmid:gongan1234#keys-1", "SM2", "SM3")
	require.Nil(t, err)

	ok, err = VerifyPKProof(msg, keyInfo.PkPEM, proof)
	require.Nil(t, err)
	require.Equal(t, true, ok)
}
