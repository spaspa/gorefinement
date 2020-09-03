package freshname

import "fmt"

var count = 0

func Generate() string {
	count++
	return fmt.Sprintf("_e%d", count)
}
