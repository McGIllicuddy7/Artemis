package Algebra

import (
	fr "artemis/Fractions"
	"fmt"
	"strconv"
)

func (poly *Polynomial) ToString() string {
	out := ""
	for i := 0; i < len(poly.data); i++ {
		if fr.Equals(poly.data[i].coef, fr.FromInt(0)) {
			continue
		}
		if !fr.Equals(poly.data[i].coef, fr.FromInt(1)) {
			out += poly.data[i].coef.ToString()
		}
		if poly.data[i].pow != 0 {
			out += "x"
			if poly.data[i].pow != 1 {
				out += fmt.Sprintf("^%d", poly.data[i].pow)
			}
		}
		if i < len(poly.data)-1 {
			if poly.data[i+1].coef.ToInt() > 0 {
				out += "+"
			}
		}
	}
	return out
}
func parsePolycule(str *string) polycule {
	var out polycule
	left := ""
	right := ""
	contains_x := false
	minus := false
	if (*str)[0] == '-' {
		minus = true
		(*str) = (*str)[1:]
	}
	for {
		if (*str)[0] == 'x' {
			contains_x = true
			(*str) = (*str)[1:]
			break
		}
		if (*str)[0] == '+' {
			(*str) = (*str)[1:]
			break
		}
		if (*str)[0] == '-' {
			break
		}
		left += string((*str)[0])
		(*str) = (*str)[1:]
		if len(*str) < 1 {
			break
		}
	}
	rightv := 0
	right_minus := false
	if contains_x {
		rightv = 1
		if len(*str) > 1 {
			if (*str)[0] == '^' {
				(*str) = (*str)[1:]
				if len(*str) > 1 {
					if (*str)[0] == '-' {
						right_minus = true
						(*str) = (*str)[1:]
					}
				}
			}
		}
		for {
			if len(*str) < 1 {
				break
			}
			if (*str)[0] == '+' {
				(*str) = (*str)[1:]
				break
			}
			if (*str)[0] == '-' {
				break
			}
			right += string((*str)[0])
			(*str) = (*str)[1:]

		}
	}
	leftv := fr.FromInt(0)
	if len(right) > 0 {
		v, err := strconv.ParseInt(right, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		rightv = int(v)
		if right_minus {
			rightv *= -1
		}
	}
	if len(left) > 0 {
		v, err := strconv.ParseFloat(left, 64)
		if err != nil {
			panic(err.Error())
		}
		leftv = fr.FromFloat(v)
		if minus {
			leftv = fr.Mult(leftv, fr.FromInt(-1))
		}
	}
	out.coef = leftv
	out.pow = rightv
	return out
}
func PolynomialFromString(str string) Polynomial {
	out := make([]polycule, 0)
	for len(str) > 0 {
		out = append(out, parsePolycule(&str))
	}
	outv := Polynomial{out}
	outv.compress()
	return outv
}
