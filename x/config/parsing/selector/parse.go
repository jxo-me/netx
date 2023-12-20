package selector

import (
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/selector"
	"github.com/jxo-me/netx/x/config"
	xs "github.com/jxo-me/netx/x/selector"
)

func ParseChainSelector(cfg *config.SelectorConfig) selector.ISelector[chain.IChainer] {
	if cfg == nil {
		return nil
	}

	var strategy selector.IStrategy[chain.IChainer]
	switch cfg.Strategy {
	case "round", "rr":
		strategy = xs.RoundRobinStrategy[chain.IChainer]()
	case "random", "rand":
		strategy = xs.RandomStrategy[chain.IChainer]()
	case "fifo", "ha":
		strategy = xs.FIFOStrategy[chain.IChainer]()
	case "hash":
		strategy = xs.HashStrategy[chain.IChainer]()
	default:
		strategy = xs.RoundRobinStrategy[chain.IChainer]()
	}
	return xs.NewSelector(
		strategy,
		xs.FailFilter[chain.IChainer](cfg.MaxFails, cfg.FailTimeout),
		xs.BackupFilter[chain.IChainer](),
	)
}

func ParseNodeSelector(cfg *config.SelectorConfig) selector.ISelector[*chain.Node] {
	if cfg == nil {
		return nil
	}

	var strategy selector.IStrategy[*chain.Node]
	switch cfg.Strategy {
	case "round", "rr":
		strategy = xs.RoundRobinStrategy[*chain.Node]()
	case "random", "rand":
		strategy = xs.RandomStrategy[*chain.Node]()
	case "fifo", "ha":
		strategy = xs.FIFOStrategy[*chain.Node]()
	case "hash":
		strategy = xs.HashStrategy[*chain.Node]()
	default:
		strategy = xs.RoundRobinStrategy[*chain.Node]()
	}

	return xs.NewSelector(
		strategy,
		xs.FailFilter[*chain.Node](cfg.MaxFails, cfg.FailTimeout),
		xs.BackupFilter[*chain.Node](),
	)
}

func DefaultNodeSelector() selector.ISelector[*chain.Node] {
	return xs.NewSelector(
		xs.RoundRobinStrategy[*chain.Node](),
		xs.FailFilter[*chain.Node](xs.DefaultMaxFails, xs.DefaultFailTimeout),
		xs.BackupFilter[*chain.Node](),
	)
}

func DefaultChainSelector() selector.ISelector[chain.IChainer] {
	return xs.NewSelector(
		xs.RoundRobinStrategy[chain.IChainer](),
		xs.FailFilter[chain.IChainer](xs.DefaultMaxFails, xs.DefaultFailTimeout),
		xs.BackupFilter[chain.IChainer](),
	)
}
