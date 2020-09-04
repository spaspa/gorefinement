package a

func f() {

	// x: { v: int | v > 0 }
	x := 2

	g(1, 1) // ok
	g(x, 0) // want "UNSAFE"
}

// g: (x { v: int | v >= 0 }, y { v: int | v >= x }) -> { v: int | v >= x * 2 }
func g(x, y int) int {
	return x + y
}
