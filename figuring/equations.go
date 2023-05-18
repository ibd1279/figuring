package figuring

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Polynomial represents a single indeterminate polynomial equation. Mainly
// used for Bezier curves.
type Polynomial interface {
	// Degree is the polynomial degree. Also known as the largest exponent.
	Degree() int
	// Coefficients is the coefficients of the equation in order from
	// highest degree to least.
	Coefficients() []float64
	// AtT evaluates the polynomial equation for the provided t value.
	AtT(float64) float64
	// Roots returns the roots of the equation.
	Roots() []float64
	// Text returns a string representing the polynomial. The first
	// argument is the value to use for the variable in the formula. the
	// second argument turns on or off the function prefix.
	Text(rune, bool) string
}

// Derivable represents a polynomial that can have a f'(t) generated from it.
type Derivable interface {
	Polynomial
	Derivative() Polynomial
}

// Constant is a polynomial in the form of f(t) = a
type Constant struct {
	a float64
}

func ConstantA(a float64) Constant {
	return Constant{a}
}

func (Constant) Degree() int                { return 0 }
func (co Constant) Coefficients() []float64 { return []float64{co.a} }
func (co Constant) AtT(t float64) float64   { return co.a }
func (co Constant) Roots() []float64        { return nil }
func (co Constant) A() float64              { return co.a }
func (co Constant) String() string          { return co.Text('t', true) }
func (co Constant) Text(unknown rune, addPrefix bool) string {
	a := co.A()
	prefix := ""
	if addPrefix {
		prefix = fmt.Sprintf("f(%c)=", unknown)
	}
	return fmt.Sprintf("%s%s(%c^0)",
		prefix,
		HumanFormat(9, a),
		unknown,
	)
}

// Linear is a polynomial in the form of f(t) = ax + b
type Linear struct {
	ab mgl64.Vec2
}

func LinearAb(a, b float64) Linear {
	return LinearFromVec2(mgl64.Vec2{a, b})
}

func LinearFromVec2(ab mgl64.Vec2) Linear {
	return Linear{
		ab: ab,
	}
}

func (Linear) Degree() int                { return 1 }
func (le Linear) Coefficients() []float64 { return le.ab[:] }
func (le Linear) AtT(t float64) float64   { return le.ab[0]*t + le.ab[1] }
func (le Linear) Roots() []float64 {
	a, b := le.Ab()
	if IsZero(a) {
		return nil
	}
	return []float64{-b / a}
}
func (le Linear) Derivative() Polynomial { return ConstantA(le.ab[0]) }
func (le Linear) Ab() (float64, float64) { return le.ab[0], le.ab[1] }
func (le Linear) String() string         { return le.Text('t', true) }
func (le Linear) Text(unknown rune, addPrefix bool) string {
	a, b := le.Ab()
	ab := '+'
	if b < 0 {
		ab = '-'
		b = -b
	}
	prefix := ""
	if addPrefix {
		prefix = fmt.Sprintf("f(%c)=", unknown)
	}
	return fmt.Sprintf("%s%s%c%c%s",
		prefix,
		HumanFormat(9, a),
		unknown,
		ab,
		HumanFormat(9, b),
	)
}

// Quadratic is a polynomial in the form of f(t) = ax^2 + bx + c
type Quadratic struct {
	abc mgl64.Vec3
}

func QuadraticAbc(a, b, c float64) Quadratic {
	return QuadraticFromVec3(mgl64.Vec3{a, b, c})
}

func QuadraticFromVec3(abc mgl64.Vec3) Quadratic {
	return Quadratic{
		abc: abc,
	}
}

func (Quadratic) Degree() int                 { return 2 }
func (qad Quadratic) Coefficients() []float64 { return qad.abc[:] }
func (qad Quadratic) AtT(t float64) float64 {
	tv := mgl64.Vec3{t * t, t, 1}
	return tv.Dot(qad.abc)
}
func (qad Quadratic) Roots() []float64 {
	a, b, c := qad.Abc()
	if IsZero(a) {
		return LinearAb(b, c).Roots()
	}

	D := b*b - 4*a*c
	if D < 0 {
		return nil
	}
	f := -b / (2 * a)
	if IsZero(D) {
		return []float64{f}
	}
	g := math.Sqrt(D) / (2 * a)
	return []float64{f + g, f - g}
}
func (qad Quadratic) Derivative() Polynomial           { return LinearAb(2*qad.abc[0], qad.abc[1]) }
func (qad Quadratic) Abc() (float64, float64, float64) { return qad.abc[0], qad.abc[1], qad.abc[2] }
func (qad Quadratic) String() string                   { return qad.Text('t', true) }
func (qad Quadratic) Text(unknown rune, addPrefix bool) string {
	a, b, c := qad.Abc()
	ab := '+'
	if b < 0 {
		ab = '-'
		b = -b
	}
	bc := '+'
	if c < 0 {
		bc = '-'
		c = -c
	}
	prefix := ""
	if addPrefix {
		prefix = fmt.Sprintf("f(%c)=", unknown)
	}
	return fmt.Sprintf("%s%s%c^2%c%s%c%c%s",
		prefix,
		HumanFormat(9, a),
		unknown,
		ab,
		HumanFormat(9, b),
		unknown,
		bc,
		HumanFormat(9, c),
	)
}

// Cubic is a polynomial in the form of f(t) = ax^3 + bx^2 + cx + d
type Cubic struct {
	abcd mgl64.Vec4
}

func CubicAbcd(a, b, c, d float64) Cubic {
	return CubicFromVec4(mgl64.Vec4{a, b, c, d})
}

func CubicFromVec4(abcd mgl64.Vec4) Cubic {
	return Cubic{
		abcd: abcd,
	}
}

func (Cubic) Degree() int                 { return 3 }
func (cub Cubic) Coefficients() []float64 { return cub.abcd[:] }
func (cub Cubic) AtT(t float64) float64 {
	tv := mgl64.Vec4{t * t * t, t * t, t, 1}
	return tv.Dot(cub.abcd)
}
func (cub Cubic) Roots() []float64 {
	a, b, c, d := cub.Abcd()
	if IsZero(a) {
		return QuadraticAbc(b, c, d).Roots()
	}
	// depress
	p := (3*a*c - b*b) / (3 * a * a)
	q := (2*b*b*b - 9*a*b*c + 27*a*a*d) / (27 * a * a * a)

	// find depressed roots
	var roots []float64
	if IsZero(p) {
		roots = []float64{-q}
	} else if IsZero(q) {
		if p > 0 {
			return nil
		}
		root := math.Sqrt(-p)
		roots = []float64{0, root, -root}
	}
	D := q*q/4 + p*p*p/27
	if IsZero(D) {
		return []float64{-1.5 * q / p, 3 * q / p}
	} else if D > 0 {
		sd := math.Sqrt(D)
		q2 := -q / 2
		u1 := q2 + sd
		u2 := q2 - sd
		root := math.Cbrt(u1) + math.Cbrt(u2)
		roots = []float64{root}
	} else {
		u := 2 * math.Sqrt(-p/3)
		t := math.Acos(3*q/p/u) / 3
		k := 2 * math.Pi / 3
		roots = []float64{u * math.Cos(t), u * math.Cos(t-k), u * math.Cos(t-2*k)}
	}

	// un-depress
	for i := 0; i < len(roots); i++ {
		roots[i] -= b / (3 * a)
	}

	return roots
}
func (cub Cubic) Derivative() Polynomial {
	a, b, c, _ := cub.Abcd()
	return QuadraticAbc(3*a, 2*b, c)
}
func (cub Cubic) Abcd() (float64, float64, float64, float64) {
	return cub.abcd[0], cub.abcd[1], cub.abcd[2], cub.abcd[3]
}
func (cub Cubic) String() string { return cub.Text('t', true) }
func (cub Cubic) Text(unknown rune, addPrefix bool) string {
	a, b, c, d := cub.Abcd()
	ab := '+'
	if b < 0 {
		ab = '-'
		b = -b
	}
	bc := '+'
	if c < 0 {
		bc = '-'
		c = -c
	}
	cd := '+'
	if d < 0 {
		cd = '-'
		d = -d
	}
	prefix := ""
	if addPrefix {
		prefix = fmt.Sprintf("f(%c)=", unknown)
	}
	return fmt.Sprintf("%s%s%c^3%c%s%c^2%c%s%c%c%s",
		prefix,
		HumanFormat(9, a),
		unknown,
		ab,
		HumanFormat(9, b),
		unknown,
		bc,
		HumanFormat(9, c),
		unknown,
		cd,
		HumanFormat(9, d),
	)
}

func IsEqualEquations[T Coefficienter](a, b T) bool {
	as, bs := a.Coefficients(), b.Coefficients()
	if len(as) != len(bs) {
		return false
	}
	for h := 0; h < len(as); h++ {
		if !IsEqual(as[h], bs[h]) {
			return false
		}
	}
	return true
}
