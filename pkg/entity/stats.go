package entity

type Statistics struct {
	All     []uint32 `json:"all"`
	Average float64  `json:"average"`
	Max     float64  `json:"max"`
}

func (s *Statistics) Add(i int) {
	s.All = append(s.All, uint32(i))
	if float64(i) > s.Max {
		s.Max = float64(i)
	}
	s.Average = (s.Average + float64(i)) / 2
}

func NewStatistics() *Statistics {
	return &Statistics{
		All: make([]uint32, 0),
	}
}

type Stats struct {
	Lifetime     *Statistics `json:"lifetime"`
	Interactions *Statistics `json:"interactions"`
	Generation   *Statistics `json:"generation"`
	DeathBy      *Statistics `json:"death_by"`
}
