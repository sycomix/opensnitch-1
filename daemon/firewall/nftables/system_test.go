package nftables

import (
	"testing"

	"github.com/evilsocket/opensnitch/daemon/firewall/nftables/exprs"
)

type sysChainsListT struct {
	family        string
	table         string
	chain         string
	expectedRules int
}

func TestAddSystemRules(t *testing.T) {
	skipIfNotPrivileged(t)

	conn, newNS = OpenSystemConn(t)
	defer CleanupSystemConn(t, newNS)
	nft.conn = conn

	cfg, err := nft.NewSystemFwConfig(nft.preloadConfCallback, nft.reloadConfCallback)
	if err != nil {
		t.Logf("Error creating fw config: %s", err)
	}

	cfg.SetFile("./testdata/test-sysfw-conf.json")
	if err := cfg.LoadDiskConfiguration(false); err != nil {
		t.Errorf("Error loading config from disk: %s", err)
	}

	nft.AddSystemRules(false, false)

	rules, _ := getRulesList(t, conn, exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_FILTER, exprs.NFT_HOOK_INPUT)
	// 3 rules in total, 1 disabled.
	if len(rules) != 1 {
		t.Errorf("test-load-conf.json mangle-output should contain only 3 rules, no -> %d", len(rules))
		for _, r := range rules {
			t.Logf("%+v", r)
		}
	}

	rules, _ = getRulesList(t, conn, exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_OUTPUT)
	// 3 rules in total, 1 disabled.
	if len(rules) != 3 {
		t.Errorf("test-load-conf.json mangle-output should contain only 3 rules, no -> %d", len(rules))
		for _, r := range rules {
			t.Log(r)
		}
	}

	rules, _ = getRulesList(t, conn, exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_FORWARD)
	// 3 rules in total, 1 disabled.
	if len(rules) != 1 {
		t.Errorf("test-load-conf.json mangle-output should contain only 3 rules, no -> %d", len(rules))
		for _, r := range rules {
			t.Log(r)
		}
	}

}

func TestFwConfDisabled(t *testing.T) {
	skipIfNotPrivileged(t)

	conn, newNS = OpenSystemConn(t)
	defer CleanupSystemConn(t, newNS)
	nft.conn = conn

	cfg, err := nft.NewSystemFwConfig(nft.preloadConfCallback, nft.reloadConfCallback)
	if err != nil {
		t.Logf("Error creating fw config: %s", err)
	}

	cfg.SetFile("./testdata/test-sysfw-conf.json")
	if err := cfg.LoadDiskConfiguration(false); err != nil {
		t.Errorf("Error loading config from disk: %s", err)
	}

	nft.AddSystemRules(false, false)

	tests := []sysChainsListT{
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_OUTPUT, 3,
		},
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_FORWARD, 1,
		},
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_FILTER, exprs.NFT_HOOK_INPUT, 1,
		},
	}

	for _, tt := range tests {
		rules, _ := getRulesList(t, conn, tt.family, tt.table, tt.chain)
		if len(rules) != 0 {
			t.Logf("%d rules found, there should be 0", len(rules))
		}
	}
}

func TestDeleteSystemRules(t *testing.T) {
	skipIfNotPrivileged(t)

	conn, newNS = OpenSystemConn(t)
	defer CleanupSystemConn(t, newNS)
	nft.conn = conn

	cfg, err := nft.NewSystemFwConfig(nft.preloadConfCallback, nft.reloadConfCallback)
	if err != nil {
		t.Logf("Error creating fw config: %s", err)
	}

	cfg.SetFile("./testdata/test-sysfw-conf.json")
	if err := cfg.LoadDiskConfiguration(false); err != nil {
		t.Errorf("Error loading config from disk: %s", err)
	}

	nft.AddSystemRules(false, false)

	tests := []sysChainsListT{
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_OUTPUT, 3,
		},
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_MANGLE, exprs.NFT_HOOK_FORWARD, 1,
		},
		{
			exprs.NFT_FAMILY_INET, exprs.NFT_CHAIN_FILTER, exprs.NFT_HOOK_INPUT, 1,
		},
	}
	for _, tt := range tests {
		rules, _ := getRulesList(t, conn, tt.family, tt.table, tt.chain)
		if len(rules) != tt.expectedRules {
			t.Errorf("%d rules found, there should be %d", len(rules), tt.expectedRules)
		}
	}

	t.Run("test-delete-system-rules", func(t *testing.T) {
		nft.DeleteSystemRules(false, false, true)
		for _, tt := range tests {
			rules, _ := getRulesList(t, conn, tt.family, tt.table, tt.chain)
			if len(rules) != 0 {
				t.Errorf("%d rules found, there should be 0", len(rules))
			}

			tbl := nft.getTable(tt.table, tt.family)
			if tbl == nil {
				t.Errorf("table %s-%s should exist", tt.table, tt.family)
			}

			/*chn := nft.getChain(tt.chain, tbl, tt.family)
			if chn == nil {
				if chains, err := conn.ListChains(); err == nil {
					for _, c := range chains {
					}
				}
				t.Errorf("chain %s-%s-%s should exist", tt.family, tt.table, tt.chain)
			}*/
		}

	})
	t.Run("test-delete-system-rules+chains", func(t *testing.T) {
	})
}
