package gogame

import (
	"reflect"
	"testing"
	"time"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	handlers := []Handler{
		&Join{ID: 7, You: true},
		&Move{Turn: 7.8, Throttle: 9.9},
		&Location{ID: 6, X: 1.1, Y: 2.2, Vx: 3.3, Vy: 4.4, Angle: 5.5, AngularVelocity: 6.6},
		&Ping{Sent: time.Now()},
	}

	results := []Handler{
		&Join{},
		&Move{},
		&Location{},
		&Ping{},
	}

	for i, h := range handlers {
		bytes, err := h.MarshalBinary()
		if err != nil {
			t.Errorf("Unexpected %v, %#v", err, h)
		}

		result := results[i]
		err = result.UnmarshalBinary(bytes)
		if err != nil {
			t.Errorf("Unexpected %v %#v %#v", err, h, result)
		}
		if !reflect.DeepEqual(h, result) {
			t.Errorf("Not equal %#v %#v", h, result)
		}
	}
}
