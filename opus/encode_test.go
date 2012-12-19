package opus

import (
	"testing"
	"fmt"
)

func TestEncoder(t *testing.T) {
	enc, err := NewEncoder(44100, 2, VOIP)
	fmt.Println(enc, err)
}
