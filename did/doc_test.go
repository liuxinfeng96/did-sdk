package did

import (
	"did-sdk/testdata"
	"fmt"
	"testing"

	"github.com/test-go/testify/require"
)

func TestGetDidMethodFromChain(t *testing.T) {

	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	method, err := GetDidMethodFromChain(c)
	require.Nil(t, err)

	fmt.Println(method)
}
