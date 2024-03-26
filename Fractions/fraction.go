package fractions

import (
	"fmt"
	"math"
)

type Fraction struct {
	num int64
	den int64
}

func NewFrac(num int, denum int) Fraction {
	out := Fraction{int64(num), int64(denum)}
	out.simplify()
	return out
}
func gcf(a int64, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcf(b, a%b)
}
func (frac *Fraction) simplify() {
	if frac.num == 0 {
		frac.den = 1
		return
	}
	if frac.den < 0 {
		frac.num *= -1
		frac.den *= -1
	}
	gc := gcf(frac.num, frac.den)
	frac.num /= gc
	frac.den /= gc
}
func (f Fraction) ToString() string {
	if f.den == 1 {
		return fmt.Sprintf("%d", f.num)
	}
	return fmt.Sprintf("%d/%d", f.num, f.den)
}
func Mult(f0 Fraction, f1 Fraction) Fraction {
	out := Fraction{(f0.num * f1.num), (f0.den * f1.den)}
	out.simplify()
	return out
}
func Divide(f0 Fraction, f1 Fraction) Fraction {
	f2 := Recip(f1)
	out := Fraction{(f0.num * f2.num), (f0.den * f2.den)}
	out.simplify()
	return out
}
func Sqrt(f0 Fraction) Fraction {
	num := FromFloat(math.Sqrt(float64(f0.num)))
	den := FromFloat(math.Sqrt(float64(f0.num)))
	return Mult(num, Recip(den))
}
func Add(f0 Fraction, f1 Fraction) Fraction {
	out := Fraction{(f0.num*f1.den + f1.num*f0.den), (f0.den * f1.den)}
	out.simplify()
	return out
}
func Scale(f0 Fraction, scalar int64) Fraction {
	out := Fraction{(f0.num * scalar), (f0.den)}
	out.simplify()
	return out
}

// subtracts f1 from f0
func Sub(f0 Fraction, f1 Fraction) Fraction {
	out := Add(f0, Scale(f1, -1))
	return out
}
func Recip(f0 Fraction) Fraction {
	return Fraction{f0.den, f0.num}
}
func Equals(f0 Fraction, f1 Fraction) bool {
	if f0.num == 0 && f1.num == 0 {
		return true
	}
	return f0.num == f1.num && f1.den == f0.den
}
func (frac Fraction) ToFloat() float64 {
	return float64(frac.num) / float64(frac.den)
}
func (frac Fraction) ToComplex() complex128 {
	return complex(frac.ToFloat(), 0)
}
func (frac Fraction) ToInt() int64 {
	return frac.num / frac.den
}
func FromFloat(v float64) Fraction {
	f := v
	count := 1
	for f != math.Floor(f) {
		f *= 10
		count *= 10
	}
	return NewFrac(int(f), count)
}
func FromInt(v int) Fraction {
	return NewFrac(v, 1)
}
func Pow(frac Fraction, power int) Fraction {
	out := FromInt(1)
	pow := power
	frc := frac
	if pow < 0 {
		frc = Recip(frc)
		pow *= -1
	}
	for i := 0; i < power; i++ {
		out = Mult(out, frc)
	}
	return out
}
