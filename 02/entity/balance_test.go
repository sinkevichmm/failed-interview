package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceValidate(t *testing.T) {
	type test struct {
		balance int
		value   int
		want    error
	}

	tests := []test{
		{
			balance: 100,
			value:   100,
			want:    nil,
		},

		{
			balance: 1001,
			value:   100,
			want:    nil,
		},
		{
			balance: 100,
			value:   1001,
			want:    ErrNotEnoughbalances,
		},
	}
	for _, tc := range tests {
		b := &Balance{Value: tc.balance}
		err := b.Validate(tc.value)

		assert.Equal(t, err, tc.want)
	}
}
