package tui

type CostViewModel struct {
	Active     bool
	Provider   string
	TotalUSD   float64
	Projected  float64
}

func NewCostViewModel() CostViewModel {
	return CostViewModel{
		Provider: "aws",
	}
}

func (m *CostViewModel) Toggle() bool {
	m.Active = !m.Active
	return m.Active
}
