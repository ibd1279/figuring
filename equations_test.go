package figuring

import (
	"math"
	"testing"
)

func TestConstantPolynomial(t *testing.T) {
	identityTests := []struct {
		eq   Constant
		s    string
		cofs []float64
	}{
		{ConstantA(12), "f(t)=12(t^0)", []float64{12}},
		{ConstantA(-3), "f(t)=-3(t^0)", []float64{-3}},
	}
	for h, test := range identityTests {
		eq := test.eq
		if s := eq.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, eq, s, test.s)
		}

		expectedDegree := 0
		if degree := eq.Degree(); degree != expectedDegree {
			t.Errorf("[%d](%v).Degree() failed. %d != %d",
				h, eq, degree, expectedDegree)
		}

		cofs := eq.Coefficients()
		if len(cofs) != len(test.cofs) {
			t.Fatalf("[%d](%v).Coefficients() length failed. %d != %d",
				h, eq, len(cofs), len(test.cofs))
		}
		for i := 0; i < len(cofs); i++ {
			if !IsEqual(cofs[i], test.cofs[i]) {
				t.Errorf("[%d][%d](%v).Coefficients() failed. %f != %f",
					h, i, eq, cofs[i], test.cofs[i])
			}
		}

		roots := eq.Roots()
		if len(roots) > 0 {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), 0)
		}
	}

	atTests := []struct {
		eq      Constant
		b, m, e float64
	}{
		{ConstantA(12), 12., 12., 12.},
	}
	for h, test := range atTests {
		eq := test.eq
		if b := eq.AtT(-10); !IsEqual(b, test.b) {
			t.Errorf("[%d](%v).AtT(-10) failed. %f != %f",
				h, eq, b, test.b)
		}
		if m := eq.AtT(0.53); !IsEqual(m, test.m) {
			t.Errorf("[%d](%v).AtT(0.53) failed. %f != %f",
				h, eq, m, test.m)
		}
		if e := eq.AtT(10); !IsEqual(e, test.e) {
			t.Errorf("[%d](%v).AtT(10) failed. %f != %f",
				h, eq, e, test.e)
		}
	}
}

func TestLinearPolynomial(t *testing.T) {
	identityTests := []struct {
		eq         Linear
		s          string
		cofs       []float64
		roots      []float64
		derivative Constant
	}{
		{LinearAb(13, 2), "f(t)=13t+2", []float64{13, 2}, []float64{-0.153846153846}, ConstantA(13)},
		{LinearAb(0.4, -2), "f(t)=0.4t-2", []float64{0.4, -2}, []float64{5}, ConstantA(0.4)},
		{LinearAb(0, -1), "f(t)=0t-1", []float64{0, -1}, nil, ConstantA(0)},
		{LinearAb(20, 60), "f(t)=20t+60", []float64{20, 60}, []float64{-3}, ConstantA(20)},
	}
	for h, test := range identityTests {
		eq := test.eq
		if s := eq.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, eq, s, test.s)
		}
		expectedDegree := 1
		if degree := eq.Degree(); degree != expectedDegree {
			t.Errorf("[%d](%v).Degree() failed. %d != %d",
				h, eq, degree, expectedDegree)
		}

		cofs := eq.Coefficients()
		if len(cofs) != len(test.cofs) {
			t.Fatalf("[%d](%v).Coefficients() length failed. %d != %d",
				h, eq, len(cofs), len(test.cofs))
		}
		for i := 0; i < len(cofs); i++ {
			if !IsEqual(cofs[i], test.cofs[i]) {
				t.Errorf("[%d][%d](%v).Coefficients() failed. %f != %f",
					h, i, eq, cofs[i], test.cofs[i])
			}
		}

		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}

		var poly Polynomial = eq
		if _, ok := poly.(Derivable); !ok {
			t.Errorf("[%d](%v).Derivitive() failed. couldn't be converted for %T",
				h, eq, eq)
		} else if deq := eq.FirstDerivative(); !IsEqualEquations(deq, test.derivative) {
			t.Errorf("[%d](%v).Derivitive() failed. %v != %v",
				h, eq, deq, test.derivative)
		}
	}
	atTests := []struct {
		eq      Linear
		b, m, e float64
	}{
		{LinearAb(13, 2), -128, 8.89, 132},
		{LinearAb(0.4, -2), -6, -1.788, 2},
		{LinearAb(0, -1), -1, -1, -1},
		{LinearAb(20, 60), -140, 70.6, 260},
	}
	for h, test := range atTests {
		eq := test.eq
		if b := eq.AtT(-10); !IsEqual(b, test.b) {
			t.Errorf("[%d](%v).AtT(-10) failed. %f != %f",
				h, eq, b, test.b)
		}
		if m := eq.AtT(0.53); !IsEqual(m, test.m) {
			t.Errorf("[%d](%v).AtT(0.53) failed. %f != %f",
				h, eq, m, test.m)
		}
		if e := eq.AtT(10); !IsEqual(e, test.e) {
			t.Errorf("[%d](%v).AtT(10) failed. %f != %f",
				h, eq, e, test.e)
		}
	}
}

func TestQuadraticPolynomial(t *testing.T) {
	identityTests := []struct {
		eq         Quadratic
		s          string
		cofs       []float64
		roots      []float64
		derivative Linear
	}{
		{QuadraticAbc(3, 13, 2), "f(t)=3t^2+13t+2", []float64{3, 13, 2},
			[]float64{-0.159734236868, -4.1735990964654}, LinearAb(6, 13)},
		{QuadraticAbc(0.2, 0.4, -2), "f(t)=0.2t^2+0.4t-2", []float64{0.2, 0.4, -2},
			[]float64{2.3166247903554, -4.3166247903554}, LinearAb(0.4, 0.4)},
		{QuadraticAbc(0, 14, -1), "f(t)=0t^2+14t-1", []float64{0, 14, -1},
			[]float64{0.0714285714286}, LinearAb(0, 14)},
		{QuadraticAbc(-30, 20, 60), "f(t)=-30t^2+20t+60", []float64{-30, 20, 60},
			[]float64{-1.1196329811802, 1.7862996478469}, LinearAb(-60, 20)},
	}
	for h, test := range identityTests {
		eq := test.eq
		if s := eq.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, eq, s, test.s)
		}
		expectedDegree := 2
		if degree := eq.Degree(); degree != expectedDegree {
			t.Errorf("[%d](%v).Degree() failed. %d != %d",
				h, eq, degree, expectedDegree)
		}

		cofs := eq.Coefficients()
		if len(cofs) != len(test.cofs) {
			t.Fatalf("[%d](%v).Coefficients() length failed. %d != %d",
				h, eq, len(cofs), len(test.cofs))
		}
		for i := 0; i < len(cofs); i++ {
			if !IsEqual(cofs[i], test.cofs[i]) {
				t.Errorf("[%d][%d](%v).Coefficients() failed. %f != %f",
					h, i, eq, cofs[i], test.cofs[i])
			}
		}

		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}

		var poly Polynomial = eq
		if _, ok := poly.(Derivable); !ok {
			t.Errorf("[%d](%v).Derivitive() failed. couldn't be converted for %T",
				h, eq, eq)
		} else if deq := eq.FirstDerivative(); !IsEqualEquations(deq, test.derivative) {
			t.Errorf("[%d](%v).Derivitive() failed. %v != %v",
				h, eq, deq, test.derivative)
		}
	}
	atTests := []struct {
		eq      Quadratic
		b, m, e float64
	}{
		{QuadraticAbc(3, 13, 2), 172, 9.7327, 432},
		{QuadraticAbc(0.2, 0.4, -2), 14, -1.73182, 22},
		{QuadraticAbc(0, 14, -1), -141, 6.42, 139},
		{QuadraticAbc(-30, 20, 60), -3140, 62.173, -2740},
	}
	for h, test := range atTests {
		eq := test.eq
		if b := eq.AtT(-10); !IsEqual(b, test.b) {
			t.Errorf("[%d](%v).AtT(-10) failed. %f != %f",
				h, eq, b, test.b)
		}
		if m := eq.AtT(0.53); !IsEqual(m, test.m) {
			t.Errorf("[%d](%v).AtT(0.53) failed. %f != %f",
				h, eq, m, test.m)
		}
		if e := eq.AtT(10); !IsEqual(e, test.e) {
			t.Errorf("[%d](%v).AtT(10) failed. %f != %f",
				h, eq, e, test.e)
		}
	}

	rootTests := []struct {
		a, b, c float64
		roots   []float64
	}{
		{-16, 23, -6, []float64{0.3424501694127, 1.0950498305873}},
		{-4, -9, 36, []float64{-4.3290014044941, 2.0790014044941}},
		{-4, -3, 6, []float64{-1.6558688457449, 0.9058688457449}},
		{-10, 4, 0, []float64{0, 0.4}},
		{3, 1, 3, nil},
		{-7, -1, 7, []float64{-1.0739763462584, 0.9311192034013}},
		{-6, 11, -6, nil},
		{-602.385273, 89.120705, 20.954727, []float64{-0.1266714957325, 0.2746178499921}},
	}
	for h, test := range rootTests {
		eq := QuadraticAbc(test.a, test.b, test.c)
		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}
	}
}

func TestCubicPolynomial(t *testing.T) {
	identityTests := []struct {
		eq         Cubic
		s          string
		cofs       []float64
		roots      []float64
		derivative Quadratic
	}{
		{CubicAbcd(-1, 3, 13, 2), "f(t)=-1t^3+3t^2+13t+2", []float64{-1, 3, 13, 2},
			[]float64{5.4518160678303, -0.1600748979807, -2.2917411698496},
			QuadraticAbc(-3, 6, 13)},
		{CubicAbcd(5, 0.2, 0.4, -2), "f(t)=5t^3+0.2t^2+0.4t-2", []float64{5, 0.2, 0.4, -2},
			[]float64{0.6882350063453}, QuadraticAbc(15, 0.4, 0.4)},
		{CubicAbcd(0, -2, 14, -1), "f(t)=0t^3-2t^2+14t-1", []float64{0, -2, 14, -1},
			[]float64{0.0721726997995, 6.9278273002005}, QuadraticAbc(0, -4, 14)},
		{CubicAbcd(5, -30, 20, 60), "f(t)=5t^3-30t^2+20t+60", []float64{5, -30, 20, 60},
			[]float64{4.5340701967227, 2.5173040450082, -1.051374241731},
			QuadraticAbc(15, -60, 20)},
	}
	for h, test := range identityTests {
		eq := test.eq
		if s := eq.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, eq, s, test.s)
		}
		expectedDegree := 3
		if degree := eq.Degree(); degree != expectedDegree {
			t.Errorf("[%d](%v).Degree() failed. %d != %d",
				h, eq, degree, expectedDegree)
		}

		cofs := eq.Coefficients()
		if len(cofs) != len(test.cofs) {
			t.Fatalf("[%d](%v).Coefficients() length failed. %d != %d",
				h, eq, len(cofs), len(test.cofs))
		}
		for i := 0; i < len(cofs); i++ {
			if !IsEqual(cofs[i], test.cofs[i]) {
				t.Errorf("[%d][%d](%v).Coefficients() failed. %f != %f",
					h, i, eq, cofs[i], test.cofs[i])
			}
		}

		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}

		var poly Polynomial = eq
		if _, ok := poly.(Derivable); !ok {
			t.Errorf("[%d](%v).Derivitive() failed. couldn't be converted for %T",
				h, eq, eq)
		} else if deq := eq.FirstDerivative(); !IsEqualEquations(deq, test.derivative) {
			t.Errorf("[%d](%v).Derivitive() failed. %v != %v",
				h, eq, deq, test.derivative)
		}
	}
	atTests := []struct {
		eq      Cubic
		b, m, e float64
	}{
		{CubicAbcd(-1, 3, 13, 2), 1172, 9.583823, -568},
		{CubicAbcd(5, 0.2, 0.4, -2), -4986, -0.987435, 5022},
		{CubicAbcd(0, -2, 14, -1), -341, 5.8582, -61},
		{CubicAbcd(5, -30, 20, 60), -8140, 62.917385, 2260},
	}
	for h, test := range atTests {
		eq := test.eq
		if b := eq.AtT(-10); !IsEqual(b, test.b) {
			t.Errorf("[%d](%v).AtT(-10) failed. %f != %f",
				h, eq, b, test.b)
		}
		if m := eq.AtT(0.53); !IsEqual(m, test.m) {
			t.Errorf("[%d](%v).AtT(0.53) failed. %f != %f",
				h, eq, m, test.m)
		}
		if e := eq.AtT(10); !IsEqual(e, test.e) {
			t.Errorf("[%d](%v).AtT(10) failed. %f != %f",
				h, eq, e, test.e)
		}
	}

	rootTests := []struct {
		a, b, c, d float64
		roots      []float64
	}{
		{3, -16, 23, -6, []float64{3, 2, 1. / 3.}},
		{1, -4, -9, 36, []float64{4, 3, -3}},
		{1, -4, -3, 6, []float64{(3. + math.Sqrt(33)) / 2, 1, (3. - math.Sqrt(33)) / 2}},
		{4, -10, 4, 0, []float64{2, 1. / 2., 0}},
		{1, 3, 1, 3, []float64{-3}},
		{1, -7, -1, 7, []float64{7, 1, -1}},
		{1, -6, 11, -6, []float64{3, 2, 1}},
		{531.105540, -602.385273, 89.120705, 20.954727, []float64{0.898616, 0.3581768, -0.1225828}},
	}
	for h, test := range rootTests {
		eq := CubicAbcd(test.a, test.b, test.c, test.d)
		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}
	}
}

func TestQuarticPolynomial(t *testing.T) {
	identityTests := []struct {
		eq         Quartic
		s          string
		cofs       []float64
		roots      []float64
		derivative Cubic
	}{
		{QuarticAbcde(-1, 3, 13, 2, 5), "f(t)=-1t^4+3t^3+13t^2+2t+5",
			[]float64{-1, 3, 13, 2, 5},
			[]float64{5.4727058789983, -2.4091319576108},
			CubicAbcd(-4, 9, 26, 2)},
		{QuarticAbcde(5, 0.2, 0.4, -2, -0.7), "f(t)=5t^4+0.2t^3+0.4t^2-2t-0.7",
			[]float64{5, 0.2, 0.4, -2, -0.7},
			[]float64{0.7882165713505, -0.3104815461368},
			CubicAbcd(20, 0.6, 0.8, -2)},
		{QuarticAbcde(0, -2, 14, -1, 28), "f(t)=0t^4-2t^3+14t^2-1t+28",
			[]float64{0, -2, 14, -1, 28},
			[]float64{7.2005794114149},
			CubicAbcd(0, -6, 28, -1)},
		{QuarticAbcde(5, -30, 20, 60, 80), "f(t)=5t^4-30t^3+20t^2+60t+80",
			[]float64{5, -30, 20, 60, 80},
			[]float64{4.0, 3.4798157487551},
			CubicAbcd(20, -90, 40, 60)},
	}
	for h, test := range identityTests {
		eq := test.eq
		if s := eq.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, eq, s, test.s)
		}
		expectedDegree := 3
		if degree := eq.Degree(); degree != expectedDegree {
			t.Errorf("[%d](%v).Degree() failed. %d != %d",
				h, eq, degree, expectedDegree)
		}

		cofs := eq.Coefficients()
		if len(cofs) != len(test.cofs) {
			t.Fatalf("[%d](%v).Coefficients() length failed. %d != %d",
				h, eq, len(cofs), len(test.cofs))
		}
		for i := 0; i < len(cofs); i++ {
			if !IsEqual(cofs[i], test.cofs[i]) {
				t.Errorf("[%d][%d](%v).Coefficients() failed. %f != %f",
					h, i, eq, cofs[i], test.cofs[i])
			}
		}

		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}

		var poly Polynomial = eq
		if _, ok := poly.(Derivable); !ok {
			t.Errorf("[%d](%v).Derivitive() failed. couldn't be converted for %T",
				h, eq, eq)
		} else if deq := eq.FirstDerivative(); !IsEqualEquations(deq, test.derivative) {
			t.Errorf("[%d](%v).Derivitive() failed. %v != %v",
				h, eq, deq, test.derivative)
		}
	}
	atTests := []struct {
		eq      Quartic
		b, m, e float64
	}{
		{QuarticAbcde(-1, 3, 13, 2, 5), -11715, 10.07942619, -5675},
		{QuarticAbcde(5, 0.2, 0.4, -2, -0.7), 49859.3, -1.22334055, 50219.3},
		{QuarticAbcde(0, -2, 14, -1, 28), 3438, 31.104846, -582},
		{QuarticAbcde(5, -30, 20, 60, 80), 81480, 113.34621405, 22680},
	}
	for h, test := range atTests {
		eq := test.eq
		if b := eq.AtT(-10); !IsEqual(b, test.b) {
			t.Errorf("[%d](%v).AtT(-10) failed. %f != %f",
				h, eq, b, test.b)
		}
		if m := eq.AtT(0.53); !IsEqual(m, test.m) {
			t.Errorf("[%d](%v).AtT(0.53) failed. %f != %f",
				h, eq, m, test.m)
		}
		if e := eq.AtT(10); !IsEqual(e, test.e) {
			t.Errorf("[%d](%v).AtT(10) failed. %f != %f",
				h, eq, e, test.e)
		}
	}

	rootTests := []struct {
		a, b, c, d, e float64
		roots         []float64
	}{
		{-1, 3, 13, 2, 5, []float64{5.4727058789983, -2.4091319576108}},
		{2, 0, 16, 0, -5, []float64{0.54878286555, -0.54878286555}},
		{2, 12, 16, 0, 0, []float64{-2, -4, 0}},
		{3, 6, -123, -126, 1080, []float64{5, -6, 3, -4}},
		{-20, 5, 17, -29, 87, []float64{1.4875831103369, -1.6820039265853}},
		{2, 4, 6, 8, 10, []float64{}},
		{3, 6, 0, -126, 1080, []float64{}},
	}
	for h, test := range rootTests {
		eq := QuarticAbcde(test.a, test.b, test.c, test.d, test.e)
		roots := eq.Roots()
		if len(roots) != len(test.roots) {
			t.Fatalf("[%d](%v).Roots() length failed. %d != %d",
				h, eq, len(roots), len(test.roots))
		}
		for i := 0; i < len(roots); i++ {
			if !IsEqual(roots[i], test.roots[i]) {
				t.Errorf("[%d][%d](%v).Roots() failed. %f != %f",
					h, i, eq, roots[i], test.roots[i])
			}
		}
	}
}
