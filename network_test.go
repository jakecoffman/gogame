package gogame

import (
	"testing"
	"reflect"
)

func TestJoin(t *testing.T) {
	join := &Join{ID: 7, You: true}
	bytes := join.Marshal()

	if bytes[0] != JOIN {
		t.Fatal("Join was", bytes[0])
	}

	join2 := &Join{}
	join2.Unmarshal(bytes)
	if !reflect.DeepEqual(join, join2) {
		t.Fatal("Not equal:", join, join2)
	}
}

func TestMove(t *testing.T) {
	move := &Move{ID: 7, AngularVelocity: 7.8, Velocity: 9.9}
	bytes := move.Marshal()

	if bytes[0] != MOVE {
		t.Fatal("Move was", bytes[0])
	}

	move2 := &Move{}
	move2.Unmarshal(bytes)
	if !reflect.DeepEqual(move, move2) {
		t.Fatal("NOt equal:", move, move2)
	}
}