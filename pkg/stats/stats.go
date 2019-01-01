package stats

type Cummulative struct {
	All     []float64 `json:"all"`
	Average float64   `json:"average"`
	Max     float64   `json:"max"`
}

func (s *Cummulative) Add(i float64) {
	s.All = append(s.All, i)
	if i > s.Max {
		s.Max = i
	}
	s.Average = (s.Average + i) / 2
}

func NewStatistics() *Cummulative {
	return &Cummulative{
		All: make([]float64, 0),
	}
}
