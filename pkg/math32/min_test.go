package math32_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/math32"
)

func TestMin(t *testing.T) {
	assert.Equal(t, math32.Min(0.0, 1.0), 0.0)
	assert.Equal(t, math32.Min(1.0, 0.0), 0.0)
}
