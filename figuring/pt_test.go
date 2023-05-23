package figuring

import (
	"math"
	"testing"
)

func TestPt(t *testing.T) {
	identityTests := []struct {
		p    Pt
		s    string
		x, y Length
	}{
		{PtXy(10, 10), "Point({10, 10})", 10, 10},
		{PtXy(-12, -32), "Point({-12, -32})", -12, -32},
	}
	for h, test := range identityTests {
		p := test.p
		if s := p.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %s != %s",
				h, p, s, test.s)
		}
		if x, y := p.XY(); !IsEqual(x, test.x) {
			t.Errorf("[%d](%v).XY().X() failed. %f != %f",
				h, p, x, test.x)
		} else if !IsEqual(y, test.y) {
			t.Errorf("[%d](%v).XY().Y() failed. %f != %f",
				h, p, y, test.y)
		}
		if x, y := p.X(), p.Y(); !IsEqual(x, test.x) {
			t.Errorf("[%d](%v).X() failed. %f != %f",
				h, p, x, test.x)
		} else if !IsEqual(y, test.y) {
			t.Errorf("[%d](%v).Y() failed. %f != %f",
				h, p, y, test.y)
		}
	}

	equalTests := []struct {
		a, b  Pt
		equal bool
	}{
		{PtXy(10, 10), PtOrig.Add(VectorIj(10, 10)), true},
		{PtXy(-12, -12), PtOrig.Add(VectorIj(-12, -12)), true},
		{PtXy(-22, -12), PtOrig.Add(VectorIj(-12, -12)), false},
		{PtXy(13, Length(math.NaN())), PtXy(13, Length(math.NaN())), false},
	}
	for h, test := range equalTests {
		eql := IsEqualPair(test.a, test.b)
		if eql != test.equal {
			t.Errorf("[%d]IsEqualPair(%v, %v) failed. %t != %t",
				h, test.a, test.b, eql, test.equal)
		}
	}

	zeroTests := []struct {
		a    Pt
		zero bool
	}{
		{PtXy(10, 0), false},
		{PtXy(0, 0), true},
		{PtXy(Length(math.Nextafter(floatZERO_TOLERANCE, math.Inf(-1))), 0), true},
		{PtXy(0, Length(math.Nextafter(floatZERO_TOLERANCE, -1))), true},
		{PtXy(0, Length(math.Nextafter(floatZERO_TOLERANCE, 1))), false},
	}
	for h, test := range zeroTests {
		zero := IsZeroPair(test.a)
		if zero != test.zero {
			t.Errorf("[%d]IsZeroPair(%v) failed. %t != %t",
				h, test.a, zero, test.zero)
		}
	}

	isErrorTests := []struct {
		a     Pt
		isErr bool
	}{
		{PtOrig, false},
		{PtXy(10, 10), false},
		{PtXy(Length(math.NaN()), 0), true},
		{PtXy(0, Length(math.NaN())), true},
		{PtXy(Length(math.Inf(1)), 0), true},
		{PtXy(Length(math.Inf(-1)), 0), true},
	}

	for h, test := range isErrorTests {
		_, err := test.a.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.a, (err != nil), test.isErr, err)
		}
	}

	sortTests := []struct {
		a   []Pt
		pts []Pt
	}{
		{
			//0
			[]Pt{PtXy(100, 100), PtXy(0, 0), PtXy(50, 50)},
			[]Pt{PtXy(0, 0), PtXy(50, 50), PtXy(100, 100)},
		}, {
			[]Pt{PtXy(100, 100), PtXy(100, 0), PtXy(100, 50)},
			[]Pt{PtXy(100, 0), PtXy(100, 50), PtXy(100, 100)},
		},
	}
	for h, test := range sortTests {
		a := test.a
		pts := SortPts(a)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]SortPts(%v) (length) failed. %v != %v",
				h, a, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]SortPts(%v) failed. %v != %v",
					h, i, a, pts[i], test.pts[i])
			}
		}
	}
}

func TestVector(t *testing.T) {
	identityTests := []struct {
		v    Vector
		s    string
		i, j Length
	}{
		{VectorIj(10, 10), "Vector(Point({10, 10}))", 10, 10},
		{VectorIj(-4.4, 3.3), "Vector(Point({-4.4, 3.3}))", -4.4, 3.3},
		{VectorIj(0.22, -0.55), "Vector(Point({0.22, -0.55}))", 0.22, -0.55},
		{VectorZero, "Vector(Point({0, 0}))", 0, 0},
	}
	for h, test := range identityTests {
		v := test.v
		if s := v.String(); s != test.s {
			t.Errorf("[%d](%v).String() failed. %v != %v",
				h, v, s, test.s)
		}

		if i, j := v.Units(); !IsEqual(i, test.i) {
			t.Errorf("[%d](%v).Units().I failed. %v != %v",
				h, v, i, test.i)
		} else if !IsEqual(j, test.j) {
			t.Errorf("[%d](%v).Units().J failed. %v != %v",
				h, v, j, test.j)
		}
	}

	var (
		iarray = []Length{
			100, 98.078528, 92.387953, 83.146961,
			70.710678, 55.557023, 38.268343, 19.509032,
			0, -19.509032, -38.268343, -55.557023,
			-70.710678, -83.146961, -92.387953, -98.078528,
			-100, -98.078528, -92.387953, -83.146961,
			-70.710678, -55.557023, -38.268343, -19.509032,
			0, 19.509032, 38.268343, 55.557023,
			70.710678, 83.146961, 92.387953, 98.078528,
		}
		jarray = []Length{
			0, 19.509032, 38.268343, 55.557023,
			70.710678, 83.146961, 92.387953, 98.078528,
			100, 98.078528, 92.387953, 83.146961,
			70.710678, 55.557023, 38.268343, 19.509032,
			0, -19.509032, -38.268343, -55.557023,
			-70.710678, -83.146961, -92.387953, -98.078528,
			-100, -98.078528, -92.387953, -83.146961,
			-70.710678, -55.557023, -38.268343, -19.509032,
		}

		scaleArray = []Vector{
			VectorIj(iarray[0], jarray[0]),
			VectorIj(iarray[1], jarray[1]),
			VectorIj(iarray[2], jarray[2]),
			VectorIj(iarray[3], jarray[3]),
			VectorIj(iarray[4], jarray[4]),
			VectorIj(iarray[5], jarray[5]),
			VectorIj(iarray[6], jarray[6]),
			VectorIj(iarray[7], jarray[7]),
			VectorIj(iarray[8], jarray[8]),
			VectorIj(iarray[9], jarray[9]),
			VectorIj(iarray[10], jarray[10]),
			VectorIj(iarray[11], jarray[11]),
			VectorIj(iarray[12], jarray[12]),
			VectorIj(iarray[13], jarray[13]),
			VectorIj(iarray[14], jarray[14]),
			VectorIj(iarray[15], jarray[15]),
			VectorIj(iarray[16], jarray[16]),
			VectorIj(iarray[17], jarray[17]),
			VectorIj(iarray[18], jarray[18]),
			VectorIj(iarray[19], jarray[19]),
			VectorIj(iarray[20], jarray[20]),
			VectorIj(iarray[21], jarray[21]),
			VectorIj(iarray[22], jarray[22]),
			VectorIj(iarray[23], jarray[23]),
			VectorIj(iarray[24], jarray[24]),
			VectorIj(iarray[25], jarray[25]),
			VectorIj(iarray[26], jarray[26]),
			VectorIj(iarray[27], jarray[27]),
			VectorIj(iarray[28], jarray[28]),
			VectorIj(iarray[29], jarray[29]),
			VectorIj(iarray[30], jarray[30]),
			VectorIj(iarray[31], jarray[31]),
		}
	)

	increment := math.Pi / 16
	for h, v1 := range scaleArray {
		theta := Radians(increment * float64(h))
		v2 := VectorIj(1, 0).Rotate(theta).Scale(100)
		eql := IsEqualPair(v1, v2)
		if !eql {
			t.Errorf("[%d]IsEqualPair(%v, %v) failed. %t != %t",
				h, v1, v2, eql, true)
		}

		v1a, v2a := v1.Angle(), v2.Angle()
		if !IsEqual(v1a, theta) || !IsEqual(v2a, theta) {
			t.Errorf("[%d](%v).Angle() failed. %v != %v != %v",
				h, v1, v1a, v2a, theta)
		}

		v1m, v2m := v1.Magnitude(), v2.Magnitude()
		if !IsEqual(v1m, 100) || !IsEqual(v2m, 100) {
			t.Errorf("[%d](%v).Magnitude() failed. %v != %v != %v",
				h, v1, v1m, v2m, 100)
		}

		v2 = VectorFromTheta(theta).Scale(100)
		eql = IsEqualPair(v1, v2)
		if !eql {
			t.Errorf("[%d]IsEqualPair(%v, %v) bis failed. %t != %t",
				h, v1, v2, eql, true)
		}

		v2a = v2.Angle()
		if !IsEqual(v1a, theta) || !IsEqual(v2a, theta) {
			t.Errorf("[%d](%v).Angle() bis failed. %v != %v",
				h, v1, v2a, theta)
		}

		v2m = v2.Magnitude()
		if !IsEqual(v2m, 100) {
			t.Errorf("[%d](%v).Magnitude() bis failed. %v != %v",
				h, v1, v2m, 100)
		}
	}

	isErrorTests := []struct {
		a     Vector
		isErr bool
	}{
		{VectorUnit, false},
		{VectorIj(10, 10), false},
		{VectorIj(0, 0), false},
		{VectorIj(Length(math.NaN()), 0), true},
		{VectorIj(0, Length(math.NaN())), true},
		{VectorIj(Length(math.Inf(1)), 0), true},
		{VectorIj(Length(math.Inf(-1)), 0), true},
	}

	for h, test := range isErrorTests {
		_, err := test.a.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.a, (err != nil), test.isErr, err)
		}
	}
}

func TestVectorTransforms(t *testing.T) {
	skewTests := []struct {
		v        Vector
		i, j     Length
		expected Vector
	}{
		{VectorIj(1, 1), 2, 0, VectorIj(3, 1)},
		{VectorIj(0, 1), 2, 0, VectorIj(2, 1)},
		{VectorIj(0, 0), 2, 0, VectorIj(0, 0)},
		{VectorIj(1, 1), 0, 2, VectorIj(1, 3)},
		{VectorIj(0, 1), 0, 2, VectorIj(0, 1)},
		{VectorIj(0, 0), 0, 2, VectorIj(0, 0)},
		{VectorIj(1, 1), 2, 2, VectorIj(3, 3)},
		{VectorIj(0, 1), 2, 2, VectorIj(2, 1)},
		{VectorIj(0, 0), 2, 2, VectorIj(0, 0)},
	}
	for h, test := range skewTests {
		v := test.v
		r := v.SkewUnits(test.i, test.j)
		if !IsEqualPair(r, test.expected) {
			t.Errorf("[%d](%v).SkewUnits(%f, %f) failed. %v != %v",
				h, v, test.i, test.j, r, test.expected)
		}
	}

	invertTests := []struct {
		v        Vector
		expected Vector
	}{
		{VectorIj(2, 1), VectorIj(-2, -1)},
		{VectorIj(-0.02, 1), VectorIj(0.02, -1)},
	}
	for h, test := range invertTests {
		v := test.v
		v1 := v.Invert()
		v2 := v.Rotate(Radians(math.Pi))
		if !IsEqualPair(v1, test.expected) || !IsEqualPair(v2, test.expected) {
			t.Errorf("[%d](%v).Invert() failed. %v != %v != %v",
				h, v, v1, v2, test.expected)
		}
	}

	normalizeTests := []struct {
		v        Vector
		expected Vector
		isErr    bool
	}{
		{VectorIj(2, 1), VectorIj(0.8944271909, 0.4472135955), false},
		{VectorIj(-0.02, 1), VectorIj(-0.0199960011, 0.9998000599), false},
		{VectorUnit, VectorIj(Length(math.Sqrt(2)/2), Length(math.Sqrt(2)/2)), false},
		{VectorZero, VectorNaN, true},
	}
	for h, test := range normalizeTests {
		v := test.v
		r, err := v.Normalize().OrErr()
		if test.isErr && err == nil {
			t.Errorf("[%d](%v).Normalize() failed. expected error. %v != %v",
				h, v, r, test.expected)
		} else if !test.isErr && !IsEqualPair(r, test.expected) {
			t.Errorf("[%d](%v).Normalize() failed. %v != %v",
				h, v, r, test.expected)
		}
	}

	addTests := []struct {
		v, b     Vector
		expected Vector
	}{
		{VectorIj(2, 1), VectorIj(1, 2), VectorIj(3, 3)},
		{VectorIj(-2, 1), VectorIj(1, 2), VectorIj(-1, 3)},
	}
	for h, test := range addTests {
		v := test.v
		r := v.Add(test.b)
		if !IsEqualPair(r, test.expected) {
			t.Errorf("[%d](%v).Add(%v) failed. %v != %v",
				h, v, test.b, r, test.expected)
		}
	}

	dotTests := []struct {
		v, b     Vector
		expected Length
	}{
		{VectorUnit.Normalize(), VectorIj(Length(math.Sqrt(2)/2), Length(math.Sqrt(2)/2)), 1},
		{VectorUnit.Normalize(), VectorUnit.Normalize().Invert(), -1},
		{VectorUnit.Normalize(), VectorUnit.Normalize().Rotate(math.Pi / 2), 0},
	}
	for h, test := range dotTests {
		v := test.v
		r := v.Dot(test.b)
		if !IsEqual(r, test.expected) {
			t.Errorf("[%d](%v).Dot(%v) failed. %f != %f",
				h, v, test.b, r, test.expected)
		}
	}
}

func TestTranformations(t *testing.T) {
	sliceTester := func(h int, s string, a, b []Pt) {
		if len(a) != len(b) {
			t.Fatalf("[%d]%s failed. len(a) != len(b) / %d != %d",
				h, s, len(a), len(b))
		}
		for i := 0; i < len(a); i++ {
			if !IsEqualPair(a[i], b[i]) {
				t.Errorf("[%d][%d]%s failed. %.9f != %.9f / %.9f != %.9f",
					h, i, s, a[i].X(), b[i].X(), a[i].Y(), b[i].Y())
			}
		}
	}

	smallLength := Length(0.5)
	midLength := Length(math.Sqrt(2) / 2)
	largeLength := Length(math.Sqrt(3) / 2)

	rotateTests := []struct {
		theta    Radians
		origin   Pt
		slice    []Pt
		expected []Pt
	}{
		{
			math.Pi / 6.,
			PtOrig,
			[]Pt{PtXy(0, 3), PtXy(4, 0)},
			[]Pt{PtXy(-smallLength*3, largeLength*3), PtXy(largeLength*4, smallLength*4)},
		}, {
			math.Pi / 3.,
			PtOrig,
			[]Pt{PtXy(0, 3), PtXy(4, 0)},
			[]Pt{PtXy(-largeLength*3, smallLength*3), PtXy(smallLength*4, largeLength*4)},
		}, {
			3. * math.Pi / 4.,
			PtOrig,
			[]Pt{PtXy(0, 3), PtXy(4, 0)},
			[]Pt{PtXy(-midLength*3, -midLength*3), PtXy(-midLength*4, midLength*4)},
		},
	}
	for h, test := range rotateTests {
		pts := RotatePts(test.theta,
			test.origin,
			test.slice)
		sliceTester(h, "RotatePts", pts, test.expected)
	}

	translateTests := []struct {
		v        Vector
		slice    []Pt
		expected []Pt
	}{
		{
			VectorIj(10, 10),
			[]Pt{PtXy(25, 25), PtXy(-25, 15), PtXy(0.01, 0)},
			[]Pt{PtXy(35, 35), PtXy(-15, 25), PtXy(10.01, 10)},
		}, {
			VectorIj(-10, 10),
			[]Pt{PtXy(25, 25), PtXy(-25, 15), PtXy(0.01, 0)},
			[]Pt{PtXy(15, 35), PtXy(-35, 25), PtXy(-9.99, 10)},
		},
	}
	for h, test := range translateTests {
		pts := TranslatePts(test.v,
			test.slice)
		sliceTester(h, "TranslatePts", pts, test.expected)
	}

	shearTests := []struct {
		v        Vector
		slice    []Pt
		expected []Pt
	}{
		{
			VectorIj(2, 0),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(3, 1), PtXy(2, 1), PtXy(0, 0)},
		},
		{
			VectorIj(0, 2),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(1, 3), PtXy(0, 1), PtXy(0, 0)},
		},
		{
			VectorIj(2, 2),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(3, 3), PtXy(2, 1), PtXy(0, 0)},
		},
	}
	for h, test := range shearTests {
		pts := ShearPts(test.v,
			test.slice)
		sliceTester(h, "ShearPts", pts, test.expected)
	}

	scaleTests := []struct {
		v        Vector
		slice    []Pt
		expected []Pt
	}{
		{
			VectorIj(5, 0),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(5, 0), PtXy(0, 0), PtXy(0, 0)},
		}, {
			VectorIj(0, 5),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(0, 5), PtXy(0, 5), PtXy(0, 0)},
		}, {
			VectorIj(-10, -10),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(-10, -10), PtXy(0, -10), PtXy(0, 0)},
		}, {
			VectorIj(5, 2),
			[]Pt{PtXy(1, 1), PtXy(0, 1), PtXy(0, 0)},
			[]Pt{PtXy(5, 2), PtXy(0, 2), PtXy(0, 0)},
		},
	}
	for h, test := range scaleTests {
		pts := ScalePts(test.v,
			test.slice)
		sliceTester(h, "ScalePts", pts, test.expected)
	}
}
