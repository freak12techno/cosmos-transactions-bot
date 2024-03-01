package alias_manager_test

import (
	"github.com/stretchr/testify/require"
	"main/pkg/alias_manager"
	configTypes "main/pkg/config/types"
	"testing"
)

func TestTypesGetNoSubscription(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{}
	found := aliases.Get("subscription", "chain", "wallet")
	require.Empty(t, found)
}

func TestTypesGetNoChain(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{},
	}
	found := aliases.Get("subscription", "chain", "wallet")
	require.Empty(t, found)
}

func TestTypesGetNoWallets(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{},
		},
	}
	found := aliases.Get("subscription", "chain", "wallet")
	require.Empty(t, found)
}

func TestTypesGetNoWallet(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Aliases: make(map[string]string),
			},
		},
	}
	found := aliases.Get("subscription", "chain", "wallet")
	require.Empty(t, found)
}

func TestTypesGetSuccess(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Aliases: map[string]string{
					"wallet": "alias",
				},
			},
		},
	}
	found := aliases.Get("subscription", "chain", "wallet")
	require.Equal(t, "alias", found)
}

func TestTypesSetAliasesAllPresent(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Aliases: map[string]string{
					"wallet": "alias",
				},
			},
		},
	}
	aliases.Set(
		"subscription",
		&configTypes.Chain{Name: "chain"},
		"wallet2",
		"alias2",
	)
	require.Equal(t, "alias2", aliases.Get("subscription", "chain", "wallet2"))
}

func TestTypesSetAliasesPresentNoWallet(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Aliases: map[string]string{},
			},
		},
	}
	aliases.Set(
		"subscription",
		&configTypes.Chain{Name: "chain"},
		"wallet",
		"alias",
	)
	require.Equal(t, "alias", aliases.Get("subscription", "chain", "wallet"))
}

func TestTypesSetAliasesPresentNoChain(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{},
	}
	aliases.Set(
		"subscription",
		&configTypes.Chain{Name: "chain"},
		"wallet",
		"alias",
	)
	require.Equal(t, "alias", aliases.Get("subscription", "chain", "wallet"))
}

func TestTypesSetAliasesPresentNoSubscription(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{}
	aliases.Set(
		"subscription",
		&configTypes.Chain{Name: "chain"},
		"wallet",
		"alias",
	)
	require.Equal(t, "alias", aliases.Get("subscription", "chain", "wallet"))
}

func TestToTomlAliasesValid(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Chain:   &configTypes.Chain{Name: "chain"},
				Aliases: map[string]string{"wallet": "alias"},
			},
		},
	}

	tomlAliases := aliases.ToTomlAliases()
	require.Len(t, *tomlAliases, 1)

	subscriptionAliases := (*tomlAliases)["subscription"]
	require.Len(t, *subscriptionAliases, 1)

	chainAliases := (*subscriptionAliases)["chain"]
	require.Len(t, *chainAliases, 1)
	require.Equal(t, (*chainAliases)["wallet"], "alias")
}

func TestToAliasesLinksValid(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{
		"subscription": &map[string]alias_manager.ChainAliases{
			"chain": alias_manager.ChainAliases{
				Chain:   &configTypes.Chain{Name: "chain"},
				Aliases: map[string]string{"wallet": "alias"},
			},
		},
	}

	links := aliases.GetAliasesLinks("subscription")
	require.Len(t, links, 1)
	require.Equal(t, links[0].Chain.Name, "chain")
	require.Len(t, links[0].Links, 1)
	require.Equal(t, links[0].Links["wallet"].Value, "wallet")
	require.Equal(t, links[0].Links["wallet"].Title, "alias")
}

func TestToAliasesLinksNoSubscription(t *testing.T) {
	t.Parallel()

	aliases := alias_manager.AllAliases{}

	links := aliases.GetAliasesLinks("subscription")
	require.Len(t, links, 0)
}
