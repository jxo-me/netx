package rate

import (
	"github.com/jxo-me/netx/core/limiter/rate"
	limiter "github.com/jxo-me/netx/core/limiter/rate"
)

type RateLimitGenerator interface {
	Limiter() limiter.ILimiter
}

type rateLimitGenerator struct {
	r float64
}

func NewRateLimitGenerator(r float64) RateLimitGenerator {
	return &rateLimitGenerator{
		r: r,
	}
}

func (p *rateLimitGenerator) Limiter() limiter.ILimiter {
	if p == nil || p.r <= 0 {
		return nil
	}
	return NewLimiter(p.r, int(p.r)+1)
}

type rateLimitSingleGenerator struct {
	limiter rate.ILimiter
}

func NewRateLimitSingleGenerator(r float64) RateLimitGenerator {
	p := &rateLimitSingleGenerator{}
	if r > 0 {
		p.limiter = NewLimiter(r, int(r)+1)
	}

	return p
}

func (p *rateLimitSingleGenerator) Limiter() limiter.ILimiter {
	return p.limiter
}
