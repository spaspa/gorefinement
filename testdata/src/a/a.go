package a

func main() {
	// a: { x: int | x >= 0 }
	a := 3

	// b: { y: int | y <= -100 }
	b := -200

	c := 1

	a = 0  // ok
	a = -1 // want "UNSAFE"
	b = 1  // want "UNSAFE"

	a = b     // want "UNSAFE"
	a = -b    // ok
	a = b * b // ok
	a = a * b // want "UNSAFE"

	a = c // want "UNSAFE"

	println(a)
}
