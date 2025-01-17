package avenger_test

import (
	"testing"

	"rakirahman.me/hello-go/avenger"
)

func TestIsAlive(t *testing.T) {
	avenger := avenger.Avenger{
		RealName: "Steven Grant Rogers",
		HeroName: "Captain America",
		Planet:   "Earth",
	}

	avenger.IsAlive()

	if !avenger.Alive {
		t.Fatalf("Expected avenger to be alive, got %v", avenger.Alive)
	}
}
