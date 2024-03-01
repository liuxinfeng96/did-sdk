package testdata

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestInstallDidContract(t *testing.T) {

	c, err := GetChainmakerClient(ConfigPath1)
	require.Nil(t, err)

	err = InstallDidContract(c)
	require.Nil(t, err)

}
