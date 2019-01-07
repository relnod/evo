package entity

// FindOldest returns the moving creature with the highest generation.
func FindOldest(creatures []*Creature) *Creature {
	var oldest *Creature
	for _, c := range creatures {
		if c.Brain == nil {
			continue
		}
		if oldest == nil {
			oldest = c
		}
		if c.Consts.Generation > oldest.Consts.Generation {
			oldest = c
		}
	}
	return oldest
}
