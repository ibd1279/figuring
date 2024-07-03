package figuring

import (
	"math"
	"testing"
)

func TestLine(t *testing.T) {
	identityTests := []struct {
		a           Line
		s           string
		xi, yi      Length
		horiz, vert bool
	}{
		{LineAbc(2, 0, 5), "2x=-5", -2.5, Length(math.NaN()), false, true},
		{LineAbc(0, 2, 5), "2y=-5", Length(math.NaN()), -2.5, true, false},
		{LineAbc(3, 5, 7), "3x+5y+7=0", -7. / 3., -1.4, false, false},
		{LineFromPt(PtXy(2, 3), PtXy(4, 4)), "1x-2y+4=0", -4, 2, false, false},
		{LineFromVector(PtXy(1, 1), VectorIj(2, 5)), "5x-2y-3=0", 0.6, -1.5, false, false},
		{LineAbc(0, 0, 12), "0x+0y=-12", Length(math.NaN()), Length(math.NaN()), false, false},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		ti, terr := test.xi.OrErr()
		if xi, err := a.XForY(0).OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).XForY(0) failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqual(xi, ti) {
			t.Errorf("[%d](%s).XForY(0) failed. %.8f != %.8f",
				h, a, xi, ti)
		}

		ti, terr = test.yi.OrErr()
		if yi, err := a.YForX(0).OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).YForX(0) failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqual(yi, ti) {
			t.Errorf("[%d](%s).YForX(0) failed. %.8f != %.8f",
				h, a, yi, ti)
		}

		if ish := a.IsHorizontal(); ish != test.horiz {
			t.Errorf("[%d](%s).IsHorizontal() failed. %t != %t",
				h, a, ish, test.horiz)
		}
		if isv := a.IsVertical(); isv != test.vert {
			t.Errorf("[%d](%s).IsVertical() failed. %t != %t",
				h, a, isv, test.vert)
		}
	}

	normalizeTests := []struct {
		a     Line
		bxerr bool
		bx    Line
		byerr bool
		by    Line
		buerr bool
		bu    Line
		v     Vector
		rads  Radians
	}{
		{
			LineAbc(2, 2, 0),
			false, LineAbc(1, 1, 0),
			false, LineAbc(1, 1, 0),
			false, LineAbc(HalfSqrtTwo, HalfSqrtTwo, 0),
			VectorIj(-1, 1).Normalize(),
			3 * math.Pi / 4,
		},
		{
			LineAbc(6, 2, 2),
			false, LineAbc(1, 1./3., 1./3.),
			false, LineAbc(3, 1, 1),
			false, LineAbc(Length(3./math.Sqrt(10)), Length(1./math.Sqrt(10)), Length(1./math.Sqrt(10))),
			VectorIj(-1, 3).Normalize(),
			1.8925468811868438,
		},
		{
			LineAbc(14, -42, 7),
			false, LineAbc(1, -3, 0.5),
			false, LineAbc(-1./3., 1, -5./30.),
			false, LineAbc(Length(1./math.Sqrt(10)), Length(-3./math.Sqrt(10)), Length(0.5/math.Sqrt(10))),
			VectorIj(0.5, 5./30.).Normalize(),
			0.3217505543958438,
		},
	}
	for h, test := range normalizeTests {
		a, terr := test.a.OrErr()
		if b, err := a.NormalizeX().OrErr(); (err != nil) != test.bxerr {
			t.Errorf("[%d](%s).NormalizeX() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualEquations(b, test.bx) {
			t.Errorf("[%d](%s).NormalizeX() failed. %v != %v",
				h, a, b, test.bx)
		}

		if b, err := a.NormalizeY().OrErr(); (err != nil) != test.byerr {
			t.Errorf("[%d](%s).NormalizeY() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualEquations(b, test.by) {
			t.Errorf("[%d](%s).NormalizeY() failed. %v != %v",
				h, a, b, test.by)
		}

		if u, err := a.NormalizeUnit().OrErr(); (err != nil) != test.buerr {
			t.Errorf("[%d](%s).NormalizeUnit() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualEquations(u, test.bu) {
			t.Errorf("[%d](%s).NormalizeUnit() failed. %v != %v",
				h, a, u, test.bu)
		}

		if v, err := a.Vector().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Vector() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualPair(v, test.v) {
			t.Errorf("[%d](%s).Vector() failed. %v != %v",
				h, a, v, test.v)
		}

		if rads, err := a.Angle().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Angle() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqual(rads, test.rads) {
			t.Errorf("[%d](%s).Angle() failed. %f != %f",
				h, a, rads, test.rads)
		}
	}

	errorTests := []struct {
		a     Line
		isErr bool
	}{
		{LineAbc(0, 0, 12), true},
		{LineAbc(3, 4, 12), false},
		{LineAbc(Length(math.NaN()), 4, 12), true},
		{LineAbc(3, Length(math.Inf(1)), 12), true},
		{LineAbc(Length(math.Inf(-1)), 4, 12), true},
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

func TestRay(t *testing.T) {
	identityTests := []struct {
		a     Ray
		s     string
		begin Pt
		vec   Vector
		angle Radians
	}{
		{
			RayFromVector(PtOrig, VectorFromTheta(1./2.)),
			"Ray(Point({0, 0}), Vector(Point({0.877582562, 0.479425539})))",
			PtOrig,
			VectorFromTheta(1. / 2.),
			1. / 2.,
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		tp, terr := test.begin.OrErr()
		if p, err := a.Begin().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Begin() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(p, tp) {
			t.Errorf("[%d](%s).Begin() failed. %v != %v",
				h, a, p, tp)
		}

		tv, terr := test.vec.OrErr()
		if v, err := a.Vector().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Vector() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(v, tv) {
			t.Errorf("[%d](%s).Vector() failed. %v != %v",
				h, a, v, tv)
		}

		if _, err := a.OrErr(); err == nil {
			if angle := a.Angle(); !IsEqual(angle, test.angle) {
				t.Errorf("[%d](%s).Length() failed. %f != %f",
					h, a, angle, test.angle)
			}
		}
	}
}

func TestSegment(t *testing.T) {
	identityTests := []struct {
		a          Segment
		s          string
		begin, end Pt
		length     Length
		angle      Radians
	}{
		{
			SegmentPt(PtXy(5, 0), PtXy(0, 5)),
			"Segment(Point({5, 0}), Point({0, 5}))",
			PtXy(5, 0), PtXy(0, 5),
			7.0710678118655, 3. * math.Pi / 4.,
		}, {
			SegmentPt(PtXy(5, 5), PtXy(0, 5)),
			"Segment(Point({5, 5}), Point({0, 5}))",
			PtXy(5, 5), PtXy(0, 5),
			5, math.Pi,
		}, {
			SegmentPt(PtXy(5, 5), PtXy(Length(math.NaN()), 5)),
			"Segment(Point({5, 5}), Point({NaN, 5}))",
			PtXy(5, 5), PtXy(Length(math.NaN()), 5),
			0, 0,
		}, {
			SegmentFromVector(PtXy(5, 0), VectorIj(-5, 5)),
			"Segment(Point({5, 0}), Point({0, 5}))",
			PtXy(5, 0), PtXy(0, 5),
			7.0710678118655, 3. * math.Pi / 4.,
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		tp, terr := test.begin.OrErr()
		if p, err := a.Begin().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Begin() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(p, tp) {
			t.Errorf("[%d](%s).Begin() failed. %v != %v",
				h, a, p, tp)
		}

		tp, terr = test.end.OrErr()
		if p, err := a.End().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).End() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(p, tp) {
			t.Errorf("[%d](%s).End() failed. %v != %v",
				h, a, p, tp)
		}

		if _, err := a.OrErr(); err == nil {
			if length := a.Length(); !IsEqual(length, test.length) {
				t.Errorf("[%d](%s).Length() failed. %f != %f",
					h, a, length, test.length)
			}
			if angle := a.Angle(); !IsEqual(angle, test.angle) {
				t.Errorf("[%d](%s).Length() failed. %f != %f",
					h, a, angle, test.angle)
			}
		}
	}

	reverseTests := []struct {
		a Segment
		r Segment
	}{
		{SegmentPt(PtXy(0, 5), PtXy(5, 0)), SegmentPt(PtXy(5, 0), PtXy(0, 5))},
		{SegmentPt(PtXy(20, 5), PtXy(5, 2)), SegmentPt(PtXy(5, 2), PtXy(20, 5))},
	}
	for h, test := range reverseTests {
		a := test.a
		r := a.Reverse()
		if IsEqualPts(a, r) {
			t.Errorf("[%d](%s).Reverse() failed (matched source). %v == %v",
				h, a, r, test.r)
		}
		if !IsEqualPts(r, test.r) {
			t.Errorf("[%d](%s).Reverse() failed. %v == %v",
				h, a, r, test.r)
		}
	}

	errorTests := []struct {
		a     Segment
		isErr bool
	}{
		{SegmentPt(PtXy(0, 0), PtXy(0, 0)), false},
		{SegmentPt(PtXy(120, 12), PtXy(455, 30)), false},
		{SegmentPt(PtXy(0, Length(math.NaN())), PtXy(0, 0)), true},
		{SegmentPt(PtXy(Length(math.Inf(-1)), 3), PtXy(3, 3)), true},
		{SegmentPt(PtXy(3, 3), PtXy(3, Length(math.Inf(1)))), true},
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