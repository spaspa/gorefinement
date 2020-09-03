package a

func main() {
	// a: { x: int | x >= 0 }
	a := 3

	// maxDiv: (x { x: int | true }, y {v: int |v > 0 } , z { v: int | v > 0 }) -> { r: int | r > x / y && r > y / z }
	maxDiv := func(x int, y, z nat) int {
		if x/y > y/z {
			return x / y
		} else {
			return y / z
		}
	}
	if true {
		oneTwoDiv(a)
		oneTwoDiv(0)	// want "UNSAFE"
	} else {
		maxDiv(5, 1, 2)
		maxDiv(6, 6, 0)	// want "UNSAFE"
	}
}

// oneTwoDiv: (x { v: int | v != 0 }) -> ({ r: int | r == 1 }, { r: int | r == 2 / x })
func oneTwoDiv(x nat) (int, int) {
	return 1, 2 / x
}

func po() {
	if true {
		println("a")
	} else {
		oneTwoDiv(0) // unsafe
	}
}

// type nat = { x: int | x >= 0 }
type nat =  int
