package conn

type ConnLimitGenerator interface {
	Limiter() Limiter
}

type connLimitGenerator struct {
	n int
}

func NewConnLimitGenerator(n int) ConnLimitGenerator {
	return &connLimitGenerator{
		n: n,
	}
}

func (p *connLimitGenerator) Limiter() Limiter {
	if p == nil || p.n <= 0 {
		return nil
	}
	return NewLimiter(p.n)
}

type connLimitSingleGenerator struct {
	limiter Limiter
}

func NewConnLimitSingleGenerator(n int) ConnLimitGenerator {
	p := &connLimitSingleGenerator{}
	if n > 0 {
		p.limiter = NewLimiter(n)
	}
	return p
}

func (p *connLimitSingleGenerator) Limiter() Limiter {
	return p.limiter
}
