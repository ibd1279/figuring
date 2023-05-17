package figuring

import (
	"fmt"
	"math"
	"testing"
)

func TestRadians(t *testing.T) {
	degreesTests := []struct {
		degrees float64
		rads    Radians
	}{
		{180.0, Radians(math.Pi)},
		{90.0, Radians(math.Pi / 2)},
		{30.0, Radians(math.Pi / 6)},
		{720.0, Radians(4 * math.Pi)},
		{-720.0, Radians(-4 * math.Pi)},
	}
	for h, test := range degreesTests {
		rads := RadiansFromDegrees(test.degrees)
		if !IsEqual(rads, test.rads) {
			t.Errorf("[%d] RadiansFromDegrees(%f). %v != %v",
				h, test.degrees, rads, test.rads)
		}
		degrees := rads.Degrees()
		if !IsEqual(degrees, test.degrees) {
			t.Errorf("[%d] (%v).Degrees() failed. %v != %v",
				h, test.degrees, degrees, test.degrees)
		}
	}

	normalizeTests := []struct {
		rads       Radians
		normalized Radians
	}{
		{Radians(math.Pi), RadiansFromFloat(math.Pi)},
		{Radians(2 * math.Pi), Radians(0)},
		{Radians(3 * math.Pi), Radians(math.Pi)},
		{Radians(-3.3 * math.Pi), Radians(0.7 * math.Pi)},
		{Radians(8 * math.Pi), Radians(0)},
		{Radians(-8 * math.Pi), Radians(0)},
		{Radians(3 * math.Pi), RadiansFromFloat(3 * math.Pi)},
		{Radians(-3 * math.Pi), RadiansFromFloat(3 * math.Pi)},
		{Radians(8 * math.Pi), RadiansFromFloat(8 * math.Pi)},
		{Radians(-8 * math.Pi), RadiansFromFloat(8 * math.Pi)},
	}
	for h, test := range normalizeTests {
		rads := test.rads.Normalize()
		if !IsEqual(rads, test.normalized) {
			t.Errorf("[%d] (%f).Normalize() failed. %v != %v",
				h, test.rads, rads, test.normalized)
		}
	}

	zeroTests := []struct {
		rads   Radians
		isZero bool
	}{
		{Radians(0), true},
		{Radians(1e-10), true},
		{Radians(1e-9 - 1e-10), true},
		{Radians(1e-9 + 1e-11), false},
		{Radians(1), false},
		{RadiansFromFloat(2 * math.Pi), true},
	}
	for h, test := range zeroTests {
		result := IsZero(test.rads)
		if result != test.isZero {
			t.Errorf("[%d] (%v).IsZero() failed. %t != %t",
				h, test.rads, result, test.isZero)
		}
	}

	stringTests := []struct {
		rads      Radians
		expected  string
		expected2 string
	}{
		{Radians(0), "θ(0.00000π)", "0.000000"},
		{Radians(math.Pi), "θ(1.00000π)", "3.141593"},
		{Radians(5 * math.Pi / 6), "θ(0.83333π)", "2.617994"},
		{Radians(15 * math.Pi / 6), "θ(2.50000π)", "7.853982"},
	}
	for h, test := range stringTests {
		str := test.rads.String()
		if str != test.expected {
			t.Errorf("[%d] (%v).String() failed. %s != %s",
				h, test.rads, str, test.expected)
		}
		str = fmt.Sprintf("%f", test.rads)
		if str != test.expected2 {
			t.Errorf("[%d] (%v).%%f failed. %s != %s",
				h, test.rads, str, test.expected2)
		}
	}

	limitTests := []struct {
		s   []Radians
		min Radians
		max Radians
	}{
		{[]Radians{-100, -10, -1, 0, 1, 10, 100}, -100, 100},
		{[]Radians{100, 10, 1, 0, -1, -10, -100}, -100, 100},
		{[]Radians{100, -100, 10, -10, 1, -1, 0}, -100, 100},
		{[]Radians{0.001, 0.002, 0.003, 0.004, 0}, 0, 0.004},
		{[]Radians{0, 0.004, 0.003, 0.002, 0.001}, 0, 0.004},
		{[]Radians{}, 0, 0},
	}
	for h, test := range limitTests {
		min := Minimum(test.s...)
		if !IsEqual(min, test.min) {
			t.Errorf("[%d]MinRadians(...) failed. %v != %v",
				h, min, test.min)
		}
		max := Maximum(test.s...)
		if !IsEqual(max, test.max) {
			t.Errorf("[%d]MaxRadians(...) failed. %v != %v",
				h, max, test.max)
		}
	}

	clampTests := []struct {
		rads     Radians
		min, max Radians
		expected Radians
	}{
		{100, -10, 10, 10},
		{-100, -10, 10, -10},
		{-5, -10, 10, -5},
		{5, -10, 10, 5},
		{0.1, 0, 1, 0.1},
		{0.9, 0, 1, 0.9},
		{-0.1, 0, 1, 0},
		{1.1, 0, 1, 1},
	}
	for h, test := range clampTests {
		r := Clamp(test.min, test.rads, test.max)
		if !IsEqual(r, test.expected) {
			t.Errorf("[%d]ClampLength(...) failed. %v != %v",
				h, r, test.expected)
		}
	}

	isErrorTests := []struct {
		rads  Radians
		isErr bool
	}{
		{Radians(0), false},
		{Radians(math.Pi), false},
		{Radians(math.NaN()), true},
		{Radians(math.Inf(1)), true},
		{Radians(math.Inf(-1)), true},
	}

	for h, test := range isErrorTests {
		_, err := test.rads.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.rads, (err != nil), test.isErr, err)
		}
	}
}

func TestLength(t *testing.T) {
	const (
		FLOAT_TOLERANCE = 0.00000001

		fivenano  Length = 0.005
		fivemicro        = 5.
		fivemilli        = 5000.
		fivecenti        = 50000.
		fivedeci         = 500000.
		fivemeter        = 5000000.
		fivedeka         = 50000000.
		fivehecto        = 500000000.
		fivekilo         = 5000000000.
		fivemega         = 5000000000000.
		fivegiga         = 5000000000000000.
		fivetera         = 5000000000000000000.

		nanohuman  = "0.005µm"
		microhuman = "5.000µm"
		millihuman = "5.000mm"
		centihuman = "50.000mm"
		decihuman  = "500.000mm"
		meterhuman = "5.000m"
		dekahuman  = "50.000m"
		hectohuman = "500.000m"
		kilohuman  = "5.000km"
		megahuman  = "5.000Mm"
		gigahuman  = "5.000Gm"
		terahuman  = "5000.000Gm"
	)
	uomTests := []struct {
		v   Length
		s   string
		uom Length
		f   float64
		f32 float32
		i   int
		i64 int64
	}{
		{fivenano, nanohuman, UOM_METER, 0.000000005, 0.000000005, 0, 0},
		{fivemicro, microhuman, UOM_METER, 0.000005, 0.000005, 0, 0},
		{fivemilli, millihuman, UOM_METER, 0.005, 0.005, 0, 0},
		{fivecenti, centihuman, UOM_METER, 0.05, 0.05, 0, 0},
		{fivedeci, decihuman, UOM_METER, 0.5, 0.5, 0, 0},
		{fivemeter, meterhuman, UOM_METER, 5., 5., 5, 5},
		{fivedeka, dekahuman, UOM_METER, 50., 50., 50, 50},
		{fivehecto, hectohuman, UOM_METER, 500., 500., 500, 500},
		{fivekilo, kilohuman, UOM_METER, 5000., 5000., 5000, 5000},
		{fivemega, megahuman, UOM_METER, 5000000., 5000000., 5000000, 5000000},
		{fivegiga, gigahuman, UOM_KILOMETER, 5000000., 5000000., 5000000, 5000000},
		{fivetera, terahuman, UOM_MICROMETER, 5000000000000000000., 5000000000000000000., 5000000000000000000, 5000000000000000000},
	}
	for h, test := range uomTests {
		uom, _ := test.v.HumanUnitLabel()
		if s := test.v.Text(uom); s != test.s {
			t.Errorf("Length[%d].String() failed. %s != %s",
				h, s, test.s)
		}
		if f := test.v.Float(test.uom); !IsEqual(f, test.f) {
			t.Errorf("Length[%d].Float() failed. %g != %g",
				h, f, test.f)
		}
		if f := test.v.Float32(test.uom); math.Abs(float64(f-test.f32)) > floatTOLERANCE {
			t.Errorf("Length[%d].Float32() failed. %g != %g",
				h, f, test.f32)
		}
		if i := test.v.Int(test.uom); i != test.i {
			t.Errorf("Length[%d].Int() failed. %d != %d",
				h, i, test.i)
		}
		if i := test.v.Int64(test.uom); i != test.i64 {
			t.Errorf("Length[%d].Int64() failed. %d != %d",
				h, i, test.i64)
		}
	}

	lengthUomTests := []struct {
		f        float64
		uom      Length
		expected Length
	}{
		{25, UOM_METER, 25 * UOM_METER},
		{fivetera, UOM_GIGAMETER, fivetera * UOM_GIGAMETER},
		{-fivetera, UOM_GIGAMETER, -fivetera * UOM_GIGAMETER},
	}
	for h, test := range lengthUomTests {
		lngth := LengthUom(test.f, test.uom)
		if !IsEqual(lngth, test.expected) {
			t.Errorf("[%d]LengthUom() failed. %v != %v",
				h, lngth, test.expected)
		}
	}

	parseUomTests := []struct {
		s string
		d Length
		v Length
	}{
		{"micro", UOM_METER, UOM_MICROMETER},
		{"µm", UOM_METER, UOM_MICROMETER},
		{"hecto", UOM_METER, UOM_HECTOMETER},
		{"nano", UOM_METER, UOM_METER},
		{"Mm", UOM_METER, UOM_MEGAMETER},
		{"pm", UOM_METER, UOM_METER},
		{"m", UOM_MILLIMETER, UOM_METER},
	}
	for h, test := range parseUomTests {
		if uom := ParseUnitOfMeasure(test.s, test.d); uom != test.v {
			t.Errorf("[%d]ParseUnitOfMeasure() failed. %v != %v",
				h, uom, test.v)
		}
	}

	// This is to get around the fact that the compiler won't let you do a
	// divide by zero with constants.
	zeroFloat := float64(0.0)
	floatTests := []struct {
		v      Length
		nan    bool
		posinf bool
		neginf bool
		r      Length
		snaps  bool
		orErr  bool
	}{
		{-0.0000000005, false, false, false, 0.0, true, false},              // 0
		{-0.00000005, false, false, false, 0.0, false, false},               // 1
		{0.0, false, false, false, 0.0, true, false},                        // 2
		{0.0000000005, false, false, false, 0.0, true, false},               // 3
		{0.00000005, false, false, false, 0.0, false, false},                // 4
		{0.1, false, false, false, 0.0, false, false},                       // 5
		{10.5, false, false, false, 11.0, false, false},                     // 6
		{2.1e-308, false, false, false, 0.0, true, false},                   // 7
		{1.6e+308, false, false, false, 1.6e+308, false, false},             // 8
		{Length(math.NaN()), true, false, false, 0, false, true},            // 9
		{Length(math.Inf(1)), false, true, false, 0, false, true},           // 10
		{Length(math.Inf(-1)), false, false, true, 0, false, true},          // 11
		{Length(zeroFloat / zeroFloat), true, false, false, 0, false, true}, // 12
		{Length(10 / zeroFloat), false, true, false, 0, false, true},        // 13
		{Length(math.Sqrt(-1)), true, false, false, 0, false, true},         // 14
	}
	for h, test := range floatTests {
		// This is to mostly ensure the tests input matches expectations.
		if math.IsNaN(float64(test.v)) != test.nan {
			t.Errorf("Length[%d].IsNaN() failed. %t != %t",
				h, math.IsNaN(float64(test.v)), test.nan)
		}
		if math.IsInf(float64(test.v), 1) != test.posinf {
			t.Errorf("Length[%d].IsInf(1) failed. %t != %t",
				h, math.IsInf(float64(test.v), 1), test.posinf)
		}
		if math.IsInf(float64(test.v), -1) != test.neginf {
			t.Errorf("Length[%d].IsInf(-1) failed. %t != %t",
				h, math.IsInf(float64(test.v), -1), test.neginf)
		}

		// Testing the non-real testing function.
		if _, err := test.v.OrErr(); test.orErr && err == nil {
			t.Errorf("Length[%d].OrErr failed. Did not return an error.", h)
		} else if !test.orErr && err != nil {
			t.Errorf("Length[%d].OrErr failed. Returned an error.", h)
		}

		if v, err := test.v.OrErr(); err == nil {
			// Test the things that only work on real numbers.
			if r := v.Round(); math.Abs(float64(r-test.r)) > floatTOLERANCE {
				t.Errorf("Length[%d].Round() failed. %g != %g", h, r, test.r)
			}

			if test.snaps && !IsZero(v) {
				t.Errorf("Length[%d].SnapToZero failed. %g != 0", h, v)
			} else if !test.snaps && IsZero(v) {
				t.Errorf("Length[%d].SnapToZero failed. %g == 0", h, v)
			}
		}

	}

	clampTests := []struct {
		lngth    Length
		min, max Length
		expected Length
	}{
		{100, -10, 10, 10},
		{-100, -10, 10, -10},
		{-5, -10, 10, -5},
		{5, -10, 10, 5},
		{0.1, 0, 1, 0.1},
		{0.9, 0, 1, 0.9},
		{-0.1, 0, 1, 0},
		{1.1, 0, 1, 1},
	}
	for h, test := range clampTests {
		r := Clamp(test.min, test.lngth, test.max)
		if !IsEqual(r, test.expected) {
			t.Errorf("[%d]ClampLength(...) failed. %v != %v",
				h, r, test.expected)
		}
	}

	isErrorTests := []struct {
		lngth Length
		isErr bool
	}{
		{Length(0), false},
		{Length(math.Pi), false},
		{Length(math.NaN()), true},
		{Length(math.Inf(1)), true},
		{Length(math.Inf(-1)), true},
	}

	for h, test := range isErrorTests {
		_, err := test.lngth.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.lngth, (err != nil), test.isErr, err)
		}
		if err != nil && test.isErr {
			nan := math.IsNaN(float64(test.lngth))
			if nan && !err.IsNaN() {
				t.Errorf("[%d] FloatingPointError.IsNaN failed. %f",
					h, test.lngth)
			} else if nan && err.Error() != "NaN encountered" {
				t.Errorf("[%d] FloatingPointError.Error() IsNaN failed. %v",
					h, err.Error())
			}
			posinf := math.IsInf(float64(test.lngth), 1)
			if posinf && !err.IsPosInf() {
				t.Errorf("[%d] FloatingPointError.IsPosInf failed. %f",
					h, test.lngth)
			} else if posinf && err.Error() != "Positive Inf encountered" {
				t.Errorf("[%d] FloatingPointError.Error() IsPosInf failed. %v",
					h, err.Error())
			}
			neginf := math.IsInf(float64(test.lngth), -1)
			if neginf && !err.IsNegInf() {
				t.Errorf("[%d] FloatingPointError.IsNegInf failed. %f",
					h, test.lngth)
			} else if neginf && err.Error() != "Negative Inf encountered" {
				t.Errorf("[%d] FloatingPointError.Error() IsNegInf failed. %v",
					h, err.Error())
			}
			inf := math.IsInf(float64(test.lngth), 0)
			if inf && !err.IsInf() {
				t.Errorf("[%d] FloatingPointError.IsInf failed. %f",
					h, test.lngth)
			}
		}
	}
}
