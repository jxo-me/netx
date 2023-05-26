package traffic

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

func (p *limitGenerator) In() Limiter {
	if p == nil || p.in <= 0 {
		return nil
	}
	return NewLimiter(p.in)
}

func (p *limitGenerator) Out() Limiter {
	if p == nil || p.out <= 0 {
		return nil
	}
	return NewLimiter(p.out)
}
