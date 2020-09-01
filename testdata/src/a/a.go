package a

func main() {
	// maxDiv: (x int, y, z { v: int | v > 0 }) -> { r: int | r > x / y && r > y / z }
	maxDiv := func (x, y, z int) int { if x / y > y / z { return x / y } else { return y / z } }

	maxDiv(5, 1, 2)
	maxDiv(6, 6, 0)	// want "UNSAFE"
}

// oneTwo: () -> ({ r: int | r == 1 }, { r: int | r == 2 })
func oneTwo() (int, int) {
	return 1, 2
}

