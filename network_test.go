package gogame

import (
	"testing"
	"reflect"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	handlers := []Handler {
		&Join{ID: 7, You: true},
		&Move{Turn: 7.8, Throttle: 9.9},
		&Location{ID: 6, X: 1.1, Y: 2.2, Vx: 3.3, Vy: 4.4, Angle: 5.5, AngularVelocity: 6.6},
	}

	results := []Handler {
		&Join{},
		&Move{},
		&Location{},
	}

	for i, h := range handlers {
		bytes := h.Marshal()

		result := results[i]
		result.Unmarshal(bytes)
		if !reflect.DeepEqual(h, result) {
			t.Fatal("Not equal", h, result)
		}
	}
}