package rate

type RateLimitGenerator interface {
	Limiter() Limiter
}

type rateLimitGenerator struct {
	r float64
}

func NewRateLimitGenerator(r float64) RateLimitGenerator {
	return &rateLimitGenerator{
		r: r,
	}
}

func (p *rateLimitGenerator) Limiter() Limiter {
	if p == nil || p.r <= 0 {
		return nil
	}
	return NewLimiter(p.r, int(p.r)+1)
}

type rateLimitSingleGenerator struct {
	limiter Limiter
}

func NewRateLimitSingleGenerator(r float64) RateLimitGenerator {
	p := &rateLimitSingleGenerator{}
	if r > 0 {
		p.limiter = NewLimiter(r, int(r)+1)
	}

	return p
}

func (p *rateLimitSingleGenerator) Limiter() Limiter {
	return p.limiter
}
