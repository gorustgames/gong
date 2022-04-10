package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandInRange(t *testing.T) {
	t.Parallel()
	for i := 1; i < 100; i++ {
		randVal := randInRange(-10, 10)
		assert.Equal(t, randVal >= -10 && randVal <= 10, true)
	}
}

func TestAiDy(t *testing.T) {
	// ball in initial position
	xPosBall := 388.
	yPosBall := 228.0
	// right bat at the very bottom
	xPos := 680.0
	yPos := 320.0

	aiOffset := 0.0
	aiDy(xPosBall, yPosBall, xPos, yPos, aiOffset)
}
