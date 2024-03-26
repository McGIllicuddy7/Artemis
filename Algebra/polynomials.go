package algebra

import (
	fr "artemis/fractions"
	"artemis/utils"
	"errors"
	"math/cmplx"
)

type polycule struct {
	coef fr.Fraction
	pow  int
}
type Polynomial struct {
	data []polycule
}

func (poly *Polynomial) Clone() Polynomial {
	out := Polynomial{make([]polycule, len(poly.data))}
	copy(out.data, poly.data)
	return out
}
func polycule_cmp(a polycule, b polycule) int {
	if a.pow > b.pow {
		return -1
	}
	if a.pow < b.pow {
		return 1
	}
	return 0
}
func polyculeMlt(a polycule, b polycule) polycule {
	return polycule{fr.Mult(a.coef, b.coef), a.pow + b.pow}
}
func slice_contains(data []int, value int) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == value {
			return true
		}
	}
	return false
}
func (poly *Polynomial) compress() {
	utils.SortInplace[polycule](poly.data, polycule_cmp)
	new_data := make([]polycule, 0)
	powers := make([]int, 0)
	for i := 0; i < len(poly.data); i++ {
		if !slice_contains(powers, poly.data[i].pow) && poly.data[i].coef.ToInt() != 0 {
			powers = append(powers, poly.data[i].pow)
		}
	}
	for i := 0; i < len(powers); i++ {
		p := polycule{fr.NewFrac(0, 1), powers[i]}
		for j := 0; j < len(poly.data); j++ {
			if poly.data[j].pow == powers[i] {
				p.coef = fr.Add(p.coef, poly.data[j].coef)
			}
		}
		new_data = append(new_data, p)
	}
	poly.data = new_data
	poly.data = poly.data[:len(powers)]
}
func (poly Polynomial) ZeroCoef() fr.Fraction {
	return poly.data[0].coef
}
func PolynomialAdd(a Polynomial, b Polynomial) Polynomial {
	var out Polynomial
	out.data = make([]polycule, 0)
	out.data = append(out.data, a.data...)
	out.data = append(out.data, b.data...)
	out.compress()
	return out
}
func PolynomialSub(a Polynomial, b Polynomial) Polynomial {
	var out Polynomial
	tmp := PolynomialScale(b, fr.NewFrac(-1, 1))
	out.data = make([]polycule, 0)
	out.data = append(out.data, a.data...)
	out.data = append(out.data, tmp.data...)
	out.compress()
	return out
}
func polyculeMultByPolynomial(a Polynomial, b polycule) Polynomial {
	out := a.Clone()
	for i := 0; i < len(out.data); i++ {
		out.data[i] = polyculeMlt(out.data[i], b)
	}
	out.compress()
	return out
}
func PolynomialMult(a Polynomial, b Polynomial) Polynomial {
	out := Polynomial{make([]polycule, 0)}
	for i := 0; i < len(b.data); i++ {
		out = PolynomialAdd(out, polyculeMultByPolynomial(a, b.data[i]))
	}
	out.compress()
	return out
}
func PolynomialScale(a Polynomial, s fr.Fraction) Polynomial {
	out := a.Clone()
	for i := 0; i < len(a.data); i++ {
		out.data[i].coef = fr.Mult(a.data[i].coef, s)
	}
	return out
}
func PolynomialScaleInplace(a Polynomial, s fr.Fraction) {
	for i := 0; i < len(a.data); i++ {
		a.data[i].coef = fr.Mult(a.data[i].coef, s)
	}
}
func NewPoly(coef fr.Fraction, pow int) Polynomial {
	return Polynomial{[]polycule{{coef, pow}}}
}
func CompPoly(coef fr.Fraction, ceof2 fr.Fraction, pow int) Polynomial {
	return Polynomial{[]polycule{{coef, 0}, {ceof2, pow}}}
}
func (poly *polycule) evaluate(x fr.Fraction) fr.Fraction {
	return fr.Mult(poly.coef, fr.Pow(x, poly.pow))
}
func (poly Polynomial) Evaluate(x fr.Fraction) fr.Fraction {
	out := fr.FromInt(0)
	for i := 0; i < len(poly.data); i++ {
		addr := poly.data[i].evaluate(x)
		out = fr.Add(out, addr)
	}
	return out
}
func (poly *polycule) evaluateComplex(x complex128) complex128 {
	v := poly.coef.ToComplex()
	p := complex(float64(poly.pow), 0)
	y := cmplx.Pow(x, p)
	return v * y
}
func (poly Polynomial) EvaluateComplex(x complex128) complex128 {
	out := complex128(0)
	for i := 0; i < len(poly.data); i++ {
		out += poly.data[i].evaluateComplex(x)
	}
	return out
}
func (poly *polycule) Derivitive() {
	if poly.pow == 0 {
		poly.coef = fr.FromInt(0)
		poly.pow = 0
		return
	}
	if poly.pow == 1 {
		poly.pow = 0
		return
	}
	poly.coef = fr.Scale(poly.coef, int64(poly.pow))
	poly.pow--
}
func (poly *polycule) Integral() error {
	if poly.pow == -1 {
		return errors.New("error unspported function \"natural log\"")
	}
	poly.coef = fr.Scale(poly.coef, fr.Recip(fr.FromInt(poly.pow+1)).ToInt())
	poly.pow++
	return nil
}
func (poly Polynomial) Derivitive() Polynomial {
	out := poly.Clone()
	for i := 0; i < len(poly.data); i++ {
		out.data[i].Derivitive()
	}
	poly.compress()
	return out
}
func (poly Polynomial) Integral() (Polynomial, error) {
	out := poly.Clone()
	for i := 0; i < len(poly.data); i++ {
		err := out.data[i].Integral()
		if err != nil {
			return NewPoly(fr.FromInt(0), 0), err
		}
	}
	return out, error(nil)
}
func (poly Polynomial) MinMaxPowers() (int, int) {
	if (len(poly.data)) < 1 {
		return 0, 0
	}
	lowest := poly.data[0].pow
	highest := poly.data[0].pow
	for i := 0; i < len(poly.data); i++ {
		if poly.data[i].pow > highest {
			highest = poly.data[i].pow
		}
		if poly.data[i].pow < lowest {
			lowest = poly.data[i].pow
		}
	}
	return lowest, highest
}
func (poly Polynomial) GetPowerCoefficient(power int) fr.Fraction {
	for i := 0; i < len(poly.data); i++ {
		if poly.data[i].pow == power {
			return poly.data[i].coef
		}
	}
	return fr.FromInt(0)
}
func (poly Polynomial) FindZero(seed complex128) complex128 {
	der := poly.Derivitive()
	value := seed
	failsafe := 0
restart:
	for i := 0; i < 1000; i++ {
		current := poly.EvaluateComplex(value)
		delta := der.EvaluateComplex(value)
		value -= current / delta
	}
	if failsafe < 10 {
		failsafe++
		value = utils.RandomComplex()
		goto restart
	}
	return value
}
func cmplxContains(slice []complex128, v complex128) bool {
	for i := 0; i < len(slice); i++ {
		if utils.ComplexNearlyEqual(slice[i], v) {
			return true
		}
	}
	return false
}
func (poly Polynomial) FindZeros() []complex128 {
	minp, maxp := poly.MinMaxPowers()
	utils.SortInplace[polycule](poly.data, polycule_cmp)
	if minp >= 0 && maxp == 2 {
		a := poly.GetPowerCoefficient(2).ToComplex()
		b := poly.GetPowerCoefficient(1).ToComplex()
		c := poly.GetPowerCoefficient(0).ToComplex()
		out := make([]complex128, 2)
		out[0] = (-b - cmplx.Sqrt(b*b-4*a*c)) / (2 * a)
		out[1] = (-b + cmplx.Sqrt(b*b-4*a*c)) / (2 * a)
		return out
	}
	zeros := make([]complex128, 0)
	num_zeros := maxp - minp
	for len(zeros) < num_zeros {
		tmp := poly.FindZero(utils.RandomComplex())
		if !cmplxContains(zeros, tmp) {
			zeros = append(zeros, tmp)
		}
	}
	return zeros
}
