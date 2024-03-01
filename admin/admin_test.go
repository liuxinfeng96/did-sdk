package admin

import (
	"did-sdk/did"
	"did-sdk/testdata"
	"testing"

	"github.com/test-go/testify/require"
)

func TestDidContractAdmin(t *testing.T) {
	creatorC, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	c, err := testdata.GetChainmakerClient(testdata.ConfigPath2)
	require.Nil(t, err)

	err = DeleteAdminForDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), creatorC)
	require.Nil(t, err)

	ok, err := IsAdminOfDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), c)
	require.Nil(t, err)
	require.Equal(t, false, ok)

	err = SetAdminForDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), creatorC)
	require.Nil(t, err)

	ok2, err := IsAdminOfDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), c)
	require.Nil(t, err)
	require.Equal(t, true, ok2)

	err = DeleteAdminForDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), creatorC)
	require.Nil(t, err)

	ok3, err := IsAdminOfDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), c)
	require.Nil(t, err)
	require.Equal(t, false, ok3)
}

func TestPermissionDidContractAdmin(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath2)
	require.Nil(t, err)

	// 添加黑名单测试
	var blackList = []string{"did:cm:test1"}
	err = did.AddDidBlackListToChain(blackList, c)
	require.NotNil(t, err)

	creatorC, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	// 添加管理员
	err = SetAdminForDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), creatorC)
	require.Nil(t, err)

	ok, err := IsAdminOfDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), c)
	require.Nil(t, err)
	require.Equal(t, true, ok)

	// 重新添加黑名单
	err = did.AddDidBlackListToChain(blackList, c)
	require.Nil(t, err)

	err = DeleteAdminForDidContract(c.GetPublicKey().ToStandardKey(), c.GetHashType(), creatorC)
	require.Nil(t, err)
}
