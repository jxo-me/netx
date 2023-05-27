package conn

import (
	limiter "github.com/jxo-me/netx/core/limiter/conn"
)

type ConnLimitGenerator interface {
	Limiter() limiter.ILimiter
}

type connLimitGenerator struct {
	n int
}

func NewConnLimitGenerator(n int) ConnLimitGenerator {
	return &connLimitGenerator{
		n: n,
	}
}

func (p *connLimitGenerator) Limiter() limiter.ILimiter {
	if p == nil || p.n <= 0 {
		return nil
	}
	return NewLimiter(p.n)
}

type connLimitSingleGenerator struct {
	limiter limiter.ILimiter
}

func NewConnLimitSingleGenerator(n int) ConnLimitGenerator {
	p := &connLimitSingleGenerator{}
	if n > 0 {
		p.limiter = NewLimiter(n)
	}
	return p
}

func (p *connLimitSingleGenerator) Limiter() limiter.ILimiter {
	return p.limiter
}
