package a

func main() {
	// a: { x: int | x >= 0 }
	a := 3

	// b: { x: int | x <= -100 }
	b := -200

	a = b*2 + 2
	a = b + 200
	a = 1 + 200

	println(a)
}
