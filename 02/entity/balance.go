package entity

// Balance data
type Balance struct {
	ID    int
	Value int
}

// Validate validate balance
func (b *Balance) Validate(v int) error {
	if b.Value-v < 0 {
		return ErrNotEnoughbalances
	}

	return nil
}
