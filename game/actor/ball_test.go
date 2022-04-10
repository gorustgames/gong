package actor

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCalculateDeflection(t *testing.T) {
	t.Parallel()

	/*
			bat JPG:
			-16 pixel margins up & down
		    -128 pixels bat itself, i.e. 64px(upper part) + 64 px(lower part)

			-------------------------
			|	 			  		|		16 px
			|		   ---   		|
			|		   | |   		|
			|		   | |   		|
			|		   |-|	  		|		80 px
			|		   | |   		|
			|		   | |   		|
			|		   ---	  		|		144 px
			|	 		    		|
			|-----------------------|		160 px

	*/

	// y(ball) position of ball vertically aligned to exact middle of the bat with y(bat)=0
	ballInTheMiddleOfBatY := 0.0 + 16 + 64 - 12

	// ball hits exactly in the middle ->  zero deflection
	def5 := calculateDeflection(ballInTheMiddleOfBatY, 0)
	assert.Equal(t, def5, 0.0)

	// ball hits above the middle ->  negative deflection
	def1 := calculateDeflection(ballInTheMiddleOfBatY-1, 0)
	def2 := calculateDeflection(ballInTheMiddleOfBatY-2, 0)
	assert.Equal(t, def1 < 0, true)
	assert.Equal(t, def2 < 0, true)
	// the more from the middle ball hits the smaller deflection
	assert.Equal(t, def2 < def1, true)

	// ball hits below the middle ->  positive deflection
	def3 := calculateDeflection(ballInTheMiddleOfBatY+1, 0)
	def4 := calculateDeflection(ballInTheMiddleOfBatY+2, 0)
	assert.Equal(t, def3 > 0, true)
	assert.Equal(t, def4 > 0, true)
	// the more from the middle ball hits the bigger deflection
	assert.Equal(t, def3 < def4, true)

	// deflections for hits in the same distance below and above the bat middle are symmetric
	assert.Equal(t, math.Abs(def1) == math.Abs(def3), true)
	assert.Equal(t, math.Abs(def2) == math.Abs(def4), true)
}
