package actor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateDeflection(t *testing.T) {
	t.Parallel()
	var def float64

	/*
			bat JPG:
			-16 pixel margins up & down
		    -128 pixels bat itself, i.e. 64px(upper part) + 64 px(lower part)

				---
				 |			16 px
				---
			    | |
				| |
				|-|			64 px
				| |
				| |
				---			128 px
				 |
				---			160 px

	*/

	// bat hit in upper part (7 pixels from bat) ->  negative deflection
	def = calculateDeflection(0+16+64-1, 0)
	assert.Equal(t, def < 0, true)

	// bat hit in lower part (1 pixel below half of the bat) ->  positive deflection
	def = calculateDeflection(0+16+64+1, 0)
	assert.Equal(t, def > 0, true)

	// bat hit exactly in the middle ->  zero deflection
	def = calculateDeflection(0+16+64, 0)
	assert.Equal(t, def, 0)
}
