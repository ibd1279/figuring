package figuring

import (
	"math"
	"testing"
)

func TestRectangle(t *testing.T) {
	identityTests := []struct {
		a        Rectangle
		s        string
		min, max Pt
		w, h     Length
	}{
		{
			//0
			RectanglePt(PtXy(2, -2), PtXy(-2, 2)),
			"Rectangle[ Polygon(Point({-2, -2}), Point({-2, 2}), Point({2, 2}), Point({2, -2})) ]",
			PtXy(-2, -2), PtXy(2, 2),
			4, 4,
		}, {
			//1
			RectangleAppend(RectanglePt(PtXy(2, -2), PtXy(-2, 2)), RectanglePt(PtXy(-1, -1), PtXy(1, 4))),
			"Rectangle[ Polygon(Point({-2, -2}), Point({-2, 4}), Point({2, 4}), Point({2, -2})) ]",
			PtXy(-2, -2), PtXy(2, 4),
			4, 6,
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}
		if min := a.MinPt(); !IsEqualPair(min, test.min) {
			t.Errorf("[%d](%s).MinPt() failed. %v != %v",
				h, a, min, test.min)
		}
		if max := a.MaxPt(); !IsEqualPair(max, test.max) {
			t.Errorf("[%d](%s).MaxPt() failed. %v != %v",
				h, a, max, test.max)
		}
		if width := a.Width(); !IsEqual(width, test.w) {
			t.Errorf("[%d](%s).Width() failed. %f != %f",
				h, a, width, test.w)
		}
		if height := a.Height(); !IsEqual(height, test.h) {
			t.Errorf("[%d](%s).Height() failed. %f != %f",
				h, a, height, test.h)
		}
		if width, height := a.Dims(); !IsEqual(width, test.w) || !IsEqual(height, test.h) {
			t.Errorf("[%d](%s).Dims() failed. (%f, %f) != (%f, %f)",
				h, a, width, height, test.w, test.h)
		}

	}

	errorTests := []struct {
		a     Rectangle
		isErr bool
	}{
		{RectanglePt(PtXy(1, 1), PtXy(5, 5)), false},
		{RectanglePt(PtXy(-1, -1), PtXy(-5, -5)), false},
		{RectanglePt(PtXy(Length(math.NaN()), 1), PtXy(5, 5)), true},
		{RectanglePt(PtXy(1, 1), PtXy(5, Length(math.NaN()))), true},
		{RectanglePt(PtXy(1, Length(math.Inf(1))), PtXy(5, 5)), true},
		{RectanglePt(PtXy(1, 1), PtXy(Length(math.Inf(-1)), 5)), true},
	}
	for h, test := range errorTests {
		a := test.a
		_, err := a.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.a, (err != nil), test.isErr, err)
		}
	}
}

func TestPolygon(t *testing.T) {
	identityTests := []struct {
		a      Polygon
		s      string
		perim  Length
		angles []Radians
	}{
		{
			TriangleEquilateral,
			"Polygon(Point({0, 0}), Point({0.866025404, -0.5}), Point({0.866025404, 0.5}))",
			3,
			[]Radians{math.Pi / 3., math.Pi / 3., math.Pi / 3},
		}, {
			Square,
			"Polygon(Point({0, 0}), Point({1, 0}), Point({1, 1}), Point({0, 1}))",
			4,
			[]Radians{math.Pi / 2., math.Pi / 2., math.Pi / 2., math.Pi / 2.},
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}
		if perim := a.Perimeter(); !IsEqual(perim, test.perim) {
			t.Errorf("[%d](%s).Perimeter() failed. %f != %f",
				h, a, perim, test.perim)
		}
		angles := a.Angles()
		if len(angles) != len(test.angles) {
			t.Fatalf("[%d](%s).Angles() failed. %v != %v",
				h, a, angles, test.angles)
		}
		for i := 0; i < len(angles); i++ {
			if !IsEqual(angles[i], test.angles[i]) {
				t.Errorf("[%d][%d](%s).Angles() failed. %v != %v",
					h, i, a, angles[i], test.angles[i])
			}
		}
	}

	errorTests := []struct {
		a     Polygon
		isErr bool
	}{
		{PolygonPt(PtXy(1, 1), PtXy(5, 5), PtXy(0, 3)), false},
		{PolygonPt(PtXy(-1, -1), PtXy(-5, -5), PtXy(0, -3)), false},
		{PolygonPt(PtXy(Length(math.NaN()), 1), PtXy(5, 5), PtXy(0, 3)), true},
		{PolygonPt(PtXy(1, 1), PtXy(5, Length(math.NaN())), PtXy(0, 3)), true},
		{PolygonPt(PtXy(1, Length(math.Inf(1))), PtXy(5, 5), PtXy(0, 3)), true},
		{PolygonPt(PtXy(1, 1), PtXy(Length(math.Inf(-1)), 5), PtXy(0, 3)), true},
	}
	for h, test := range errorTests {
		a := test.a
		_, err := a.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.a, (err != nil), test.isErr, err)
		}
	}
}
