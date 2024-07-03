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
	// AtT evaluates the polynomial equation for the provided t value.
	AtT(float64) float64
	// Roots returns the roots (where the equation is equal to zero).
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
func (co Constant) AtT(float64) float64     { return co.a }
func (co Constant) Roots() []float64        { return nil }
func (co Constant) Derivative() Polynomial  { return ConstantA(0) }
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
func (le Linear) Derivative() Polynomial    { return le.FirstDerivative() }
func (le Linear) FirstDerivative() Constant { return ConstantA(le.ab[0]) }
func (le Linear) Ab() (float64, float64)    { return le.ab[0], le.ab[1] }
func (le Linear) String() string            { return le.Text('t', true) }
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
func (qad Quadratic) Derivative() Polynomial           { return qad.FirstDerivative() }
func (qad Quadratic) FirstDerivative() Linear          { return LinearAb(2*qad.abc[0], qad.abc[1]) }
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
func (cub Cubic) Derivative() Polynomial { return cub.FirstDerivative() }
func (cub Cubic) FirstDerivative() Quadratic {
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

// Quartic is a polynomial in the form of f(t) = ax^4 + bx^3 + cx^2 + dx + e
type Quartic struct {
	abcde [5]float64
}

func QuarticAbcde(a, b, c, d, e float64) Quartic {
	return Quartic{
		abcde: [5]float64{a, b, c, d, e},
	}
}

func (Quartic) Degree() int                 { return 3 }
func (qrt Quartic) Coefficients() []float64 { return qrt.abcde[:] }
func (qrt Quartic) AtT(t float64) float64 {
	tv := [5]float64{t * t * t * t, t * t * t, t * t, t, 1}
	var result float64
	for h := 0; h < len(tv); h++ {
		result = math.FMA(tv[h], qrt.abcde[h], result)
	}

	return result
}
func (qrt Quartic) discriminant() (float64, float64, float64, float64, float64) {
	a, b, c, d, e := qrt.Abcde()

	delta := 256.0*(a*a*a)*(e*e*e) - 192.0*(a*a)*b*d*(e*e) - 128.0*(a*a)*(c*c)*(e*e) +
		144.0*(a*a)*c*(d*d)*e - 27.0*(a*a)*(d*d*d*d) + 144.0*a*(b*b)*c*(e*e) -
		6.0*a*(b*b)*(d*d)*e - 80.0*a*b*(c*c)*d*e + 18.0*a*b*c*(d*d*d) +
		16.0*a*(c*c*c*c)*e - 4.0*a*(c*c*c)*(d*d) - 27.0*(b*b*b*b)*(e*e) +
		18.0*(b*b*b)*c*d*e - 4.0*(b*b*b)*(d*d*d) - 4.0*(b*b)*(c*c*c)*e +
		(b*b)*(c*c)*(d*d)

	P := 8.0*a*c - 3.0*(b*b)
	R := b*b*b + 8.0*d*(a*a) - 4.0*a*b*c
	delta0 := c*c - 3.0*b*d + 12.0*a*e
	D := 64.0*(a*a*a)*e - 16.0*(a*a)*(c*c) + 16.0*a*(b*b)*c - 16.0*(a*a)*b*d - 3.0*(b*b*b*b)

	return delta, P, R, delta0, D
}
func (qrt Quartic) depressedRoots(P, R, D float64) []float64 {
	// see https://github.com/vorot/roots/blob/master/src/analytical/quartic_depressed.rs

	a4, a3, a2, a1, a0 := qrt.Abcde()
	undepress := func(roots []float64) []float64 {
		for h := 0; h < len(roots); h++ {
			roots[h] = roots[h] - a3/(4.0*a4)
		}
		return roots
	}
	p := P / (8.0 * a4 * a4)
	q := R / (8.0 * a4 * a4 * a4)
	r := (D + 16.0*(a4*a4)*(12.0*a0*a4-3.0*a1*a3+a2*a2)) / (256.0 * (a4 * a4 * a4 * a4))

	// x^4 + px^2 + qx + r = 0
	if IsZero(r) {
		roots := append([]float64{0}, CubicAbcd(1, 0, p, q).Roots()...)
		return undepress(roots)
	} else if IsZero(q) {
		roots := make([]float64, 0, 4)
		for _, root := range QuadraticAbc(1, p, r).Roots() {
			if IsZero(root) {
				roots = append(roots, 0.0)
			} else if root > 0 {
				x := math.Sqrt(root)
				roots = append(roots, x, -x)
			}
		}
		return undepress(roots)
	}

	b2 := p * 5.0 / 2.0
	b1 := 2.0*(p*p) - r
	halfq := q / 2.0
	b0 := ((p * p * p) - p*r - (halfq * halfq)) / 2.0

	resolvent_roots := CubicAbcd(1, b2, b1, b0).Roots()
	y := resolvent_roots[len(resolvent_roots)-1]

	p2y := p + 2*y
	if p2y > 0 {
		sqrt_p2y := math.Sqrt(p2y)
		q0a := p + y - halfq/sqrt_p2y
		q0b := p + y + halfq/sqrt_p2y

		roots := QuadraticAbc(1.0, sqrt_p2y, q0a).Roots()
		roots = append(roots, QuadraticAbc(1.0, -sqrt_p2y, q0b).Roots()...)
		return undepress(roots)
	}
	return []float64{}
}
func (qrt Quartic) Roots() []float64 {
	// see https://en.wikipedia.org/wiki/Quartic_function#Nature_of_the_roots
	a4, a3, a2, a1, a0 := qrt.Abcde()
	if IsZero(a4) {
		return CubicAbcd(a3, a2, a1, a0).Roots()
	} else if IsZero(a0) {
		roots := make([]float64, 0, 4)
		for _, root := range CubicAbcd(a4, a3, a2, a1).Roots() {
			if !IsZero(root) {
				roots = append(roots, root)
			}
		}
		return append(roots, 0.0)
	} else if IsZero(a1) && IsZero(a3) {
		roots := make([]float64, 0, 4)
		for _, root := range QuadraticAbc(a4, a2, a0).Roots() {
			if IsZero(root) {
				roots = append(roots, 0.0)
			} else if root > 0 {
				x := math.Sqrt(root)
				roots = append(roots, x, -x)
			}
		}
		return roots
	}

	delta, P, R, delta0, D := qrt.discriminant()

	if IsZero(delta) {
		if IsZero(D) && IsZero(delta0) {
			// If ∆ = 0 then
			// if D = 0, then
			// If ∆0 = 0, all four roots are equal to −b/4a
			return []float64{-a3 / 4.0 * a4}
		} else if IsZero(delta0) {
			// If ∆ = 0 then
			// If ∆0 = 0 and D ≠ 0,
			// there are a triple root and a simple root, all real.
			x0 := (-72.0*(a4*a4)*a0 + 10.0*a4*(a2*a2) - 3.0*(a3*a3)*a2) /
				(9.0 * (8.0*(a4*a4)*a1 - 4.0*a4*a3*a2 + a3*a3*a3))
			x1 := -(a3/a4 + 3.0*x0)
			return []float64{x0, x1}
		} else if IsZero(D) && P > 0 && IsZero(R) {
			// If P > 0 and R = 0, there are two complex conjugate double roots.
			return []float64{}
		}
	} else if delta > 0 && (P > 0 || D > 0) {
		return []float64{}
	}

	return qrt.depressedRoots(P, R, D)
}
func (qrt Quartic) Derivative() Polynomial { return qrt.FirstDerivative() }
func (qrt Quartic) FirstDerivative() Cubic {
	a, b, c, d, _ := qrt.Abcde()
	return CubicAbcd(4*a, 3*b, 2*c, d)
}
func (qrt Quartic) Abcde() (float64, float64, float64, float64, float64) {
	return qrt.abcde[0], qrt.abcde[1], qrt.abcde[2], qrt.abcde[3], qrt.abcde[4]
}
func (qrt Quartic) String() string { return qrt.Text('t', true) }
func (qrt Quartic) Text(unknown rune, addPrefix bool) string {
	a, b, c, d, e := qrt.Abcde()
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
	de := '+'
	if e < 0 {
		de = '-'
		e = -e
	}
	prefix := ""
	if addPrefix {
		prefix = fmt.Sprintf("f(%c)=", unknown)
	}
	return fmt.Sprintf("%s%s%c^4%c%s%c^3%c%s%c^2%c%s%c%c%s",
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
		unknown,
		de,
		HumanFormat(9, e),
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
