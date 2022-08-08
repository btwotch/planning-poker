package main

import (
	"fmt"
	"testing"
)

func TestAverageChoice(t *testing.T) {
	m := newModel()

	tests := []struct {
		chosen bool
		choice uint8
	}{
		{true, 5},
		{true, 3},
		{true, 20},
		{false, 200},
	}

	for i, c := range tests {
		p := player{}
		p.chosen = c.chosen
		p.choice = c.choice

		p.name = fmt.Sprintf("%d", i)
		m.addPlayer(&p)
	}

	t.Logf("avg: %f\n", m.getAverageChoice())
	if uint8(m.getAverageChoice()) != 9 {
		t.Fatalf("Average calculation is wrong")
	}
}
