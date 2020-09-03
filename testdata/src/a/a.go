package a

func main() {
	// a: { x: int | x >= 0 }
	a := 3

	// b: { y: int | y <= -100 }
	b := -200

	a = 0
	a = -1 // want "UNSAFE"

	b = 1 // want "UNSAFE"

	a = b // want "UNSAFE"

	println(a)
}
