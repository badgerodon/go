package sdl

import (
	"testing"
)

func TestMain(t *testing.T) {
	err := Init(INIT_EVERYTHING)
	if err != nil {
		t.Fatal(err)
	}
	defer Quit()

	screen, err := SetVideoMode(800, 600, 32, HWSURFACE)
	if err != nil {
		t.Fatal(err)
	}

	hello, err := Load("example.png")
	if err != nil {
		t.Fatal(err)
	}

	err = BlitSurface(hello, nil, screen, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = Flip(screen)
	if err != nil {
		t.Fatal(err)
	}

	Delay(2000)
}
