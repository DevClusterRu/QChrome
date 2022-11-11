package main

import (
	"QChromium/internal"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(internal.RandStringRunes())
}
