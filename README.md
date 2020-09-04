# gorefinement

Refinement type checking for golang.

```go
func f() {
  // a: { v: int | v >= 0 }
  a := 3
  
  // b: { v: int | v <= -100 }
  b := 3
  
  a = 0   // ok
  a = -1  // reports "UNSAFE"
  b = 1   // reports "UNSAFE"
  
  a = b   // reports "UNSAFE"
  a = -b  // ok
  a = b*b // ok
  a = a*b // reports "UNSAFE", because assignments like (a, b) = (0, -100) can violate condition
}
```

## requirement
[Z3 SMT Solver](https://github.com/Z3Prover/z3) is required.

Download pre-built binaries from [Z3 repository releases page](https://github.com/Z3Prover/z3/releases),
then place header and library files to appropriate path.

If you are using macOS, install from [HomeBrew](https://formulae.brew.sh/formula/z3) is recommended.

## limitation
- 
