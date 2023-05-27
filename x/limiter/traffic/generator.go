package traffic

import (
	limiter "github.com/jxo-me/netx/core/limiter/traffic"
)

type limitGenerator struct {
	in  int
	out int
}

func newLimitGenerator(in, out int) *limitGenerator {
	return &limitGenerator{
		in:  in,
		out: out,
	}
}

func (p *limitGenerator) In() limiter.ILimiter {
	if p == nil || p.in <= 0 {
		return nil
	}
	return NewLimiter(p.in)
}

func (p *limitGenerator) Out() limiter.ILimiter {
	if p == nil || p.out <= 0 {
		return nil
	}
	return NewLimiter(p.out)
}
