package types_test

import (
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainGetPrettyName(t *testing.T) {
	t.Parallel()

	require.Equal(t, "name", types.Chain{Name: "name"}.GetName())
	require.Equal(t, "Name", types.Chain{Name: "name", PrettyName: "Name"}.GetName())
}

func TestChainsFindByName(t *testing.T) {
	t.Parallel()

	chains := types.Chains{
		{Name: "name"},
	}

	require.NotNil(t, chains.FindByName("name"))
	require.Nil(t, chains.FindByName("name-2"))
}

func TestChainsFindByChainID(t *testing.T) {
	t.Parallel()

	chains := types.Chains{
		{Name: "name", ChainID: "chain-id"},
	}

	require.NotNil(t, chains.FindByChainID("chain-id"))
	require.Nil(t, chains.FindByChainID("chain-id-2"))
}

func TestChainGetWalletLink(t *testing.T) {
	t.Parallel()

	chain1 := types.Chain{
		Name:    "name",
		ChainID: "chain-id",
	}
	link1 := chain1.GetWalletLink("wallet")
	require.Equal(t, "wallet", link1.Value)
	require.Empty(t, link1.Href)
	require.Empty(t, link1.Title)

	chain2 := types.Chain{
		Name:     "name",
		ChainID:  "chain-id",
		Explorer: &types.Explorer{WalletLinkPattern: "test/%s"},
	}

	link2 := chain2.GetWalletLink("wallet")
	require.Equal(t, "wallet", link2.Value)
	require.Equal(t, "test/wallet", link2.Href)
	require.Empty(t, link2.Title)
}

func TestChainGetValidatorLink(t *testing.T) {
	t.Parallel()

	chain1 := types.Chain{
		Name:    "name",
		ChainID: "chain-id",
	}
	link1 := chain1.GetValidatorLink("validator")
	require.Equal(t, "validator", link1.Value)
	require.Empty(t, link1.Href)
	require.Empty(t, link1.Title)

	chain2 := types.Chain{
		Name:     "name",
		ChainID:  "chain-id",
		Explorer: &types.Explorer{ValidatorLinkPattern: "test/%s"},
	}

	link2 := chain2.GetValidatorLink("validator")
	require.Equal(t, "validator", link2.Value)
	require.Equal(t, "test/validator", link2.Href)
	require.Empty(t, link2.Title)
}

func TestChainGetProposalLink(t *testing.T) {
	t.Parallel()

	chain1 := types.Chain{
		Name:    "name",
		ChainID: "chain-id",
	}
	link1 := chain1.GetProposalLink("proposal")
	require.Equal(t, "proposal", link1.Value)
	require.Empty(t, link1.Href)
	require.Empty(t, link1.Title)

	chain2 := types.Chain{
		Name:     "name",
		ChainID:  "chain-id",
		Explorer: &types.Explorer{ProposalLinkPattern: "test/%s"},
	}

	link2 := chain2.GetProposalLink("proposal")
	require.Equal(t, "proposal", link2.Value)
	require.Equal(t, "test/proposal", link2.Href)
	require.Empty(t, link2.Title)
}

func TestChainGetTransactionLink(t *testing.T) {
	t.Parallel()

	chain1 := types.Chain{
		Name:    "name",
		ChainID: "chain-id",
	}
	link1 := chain1.GetTransactionLink("transaction")
	require.Equal(t, "transaction", link1.Value)
	require.Empty(t, link1.Href)
	require.Empty(t, link1.Title)

	chain2 := types.Chain{
		Name:     "name",
		ChainID:  "chain-id",
		Explorer: &types.Explorer{TransactionLinkPattern: "test/%s"},
	}

	link2 := chain2.GetTransactionLink("transaction")
	require.Equal(t, "transaction", link2.Value)
	require.Equal(t, "test/transaction", link2.Href)
	require.Empty(t, link2.Title)
}

func TestChainGetBlockLink(t *testing.T) {
	t.Parallel()

	chain1 := types.Chain{
		Name:    "name",
		ChainID: "chain-id",
	}
	link1 := chain1.GetBlockLink(1337)
	require.Equal(t, "1337", link1.Value)
	require.Empty(t, link1.Href)
	require.Empty(t, link1.Title)

	chain2 := types.Chain{
		Name:     "name",
		ChainID:  "chain-id",
		Explorer: &types.Explorer{BlockLinkPattern: "test/%s"},
	}

	link2 := chain2.GetBlockLink(1337)
	require.Equal(t, "1337", link2.Value)
	require.Equal(t, "test/1337", link2.Href)
	require.Empty(t, link2.Title)
}
