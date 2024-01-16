package key

import (
	"fmt"
	"testing"

	"github.com/test-go/testify/require"
)

func TestGenerateKey(t *testing.T) {
	for _, v := range SupportAlgorithm {
		keyInfo, err := GenerateKey(v)
		require.Nil(t, err)

		fmt.Printf("%+v\n", keyInfo)
	}
}

func TestIsSupportAlgorithm(t *testing.T) {
	for _, v := range SupportAlgorithm {
		ok := IsSupportAlgorithm(v)
		require.Equal(t, true, ok)
	}
}
