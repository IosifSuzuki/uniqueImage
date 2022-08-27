package prng

type PRNG struct {
	Seed int
}

func (p *PRNG) GenerateNum() int {
	p.Seed = (p.Seed*73129 + 95121) % 100000
	return p.Seed
}
