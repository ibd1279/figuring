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
		{LineAbc(2, 0, 5), "2x=5", 2.5, Length(math.NaN()), false, true},
		{LineAbc(0, 2, 5), "2y=5", Length(math.NaN()), 2.5, true, false},
		{LineAbc(3, 5, 7), "3x+5y=7", 7. / 3., 1.4, false, false},
		{LineFromPt(PtXy(2, 3), PtXy(4, 4)), "1x-2y=-4", -4, 2, false, false},
		{LineFromVector(PtXy(1, 1), VectorIj(2, 5)), "5x-2y=3", 0.6, -1.5, false, false},
		{LineAbc(0, 0, 12), "0x+0y=12", Length(math.NaN()), Length(math.NaN()), false, false},
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
		v     Vector
		rads  Radians
	}{
		{
			LineAbc(2, 2, 0),
			false, LineAbc(1, 1, 0),
			false, LineAbc(1, 1, 0),
			VectorIj(-1, 1).Normalize(),
			3 * math.Pi / 4,
		},
		{
			LineAbc(6, 2, 2),
			false, LineAbc(1, 1./3., 1./3.),
			false, LineAbc(3, 1, 1),
			VectorIj(-1, 3).Normalize(),
			1.8925468811868438,
		},
		{
			LineAbc(14, -42, 7),
			false, LineAbc(1, -3, 0.5),
			false, LineAbc(-1./3., 1, -5./30.),
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

func TestParamCurve(t *testing.T) {
	linearTests := []struct {
		p1, p2   Pt
		s        string
		xc, yc   []float64
		p33, p67 Pt
	}{
		{
			PtXy(0, 10), PtXy(20, 15), "Curve(20t+0, 5t+10, t, 0, 1)",
			[]float64{20, 0}, []float64{5, 10},
			PtXy(6.6, 11.65), PtXy(13.4, 13.35),
		}, {
			PtXy(-10, -10), PtXy(-30, -5), "Curve(-20t-10, 5t-10, t, 0, 1)",
			[]float64{-20, -10}, []float64{5, -10},
			PtXy(-16.6, -8.35), PtXy(-23.4, -6.65),
		},
	}
	for h, test := range linearTests {
		a := ParamLinear(test.p1, test.p2)
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		if p := a.PtAtT(0); !IsEqualPair(p, test.p1) {
			t.Errorf("[%d](%s).PtAtT(0) failed. %v != %v",
				h, a, p, test.p1)
		}
		if p := a.PtAtT(0.33); !IsEqualPair(p, test.p33) {
			t.Errorf("[%d](%s).PtAtT(0.33) failed. %v != %v",
				h, a, p, test.p33)
		}
		if p := a.PtAtT(0.67); !IsEqualPair(p, test.p67) {
			t.Errorf("[%d](%s).PtAtT(0.67) failed. %v != %v",
				h, a, p, test.p67)
		}
		if p := a.PtAtT(2); !IsEqualPair(p, test.p2) {
			t.Errorf("[%d](%s).PtAtT(2) failed. %v != %v",
				h, a, p, test.p2)
		}
	}

	quadraticTests := []struct {
		p1, p2, p3    Pt
		s             string
		xc, yc        []float64
		p33, p50, p67 Pt
		t33, t67      Vector
	}{
		{
			PtXy(70, 250), PtXy(20, 110), PtXy(220, 60),
			"Curve(250t^2-100t+70, 90t^2-280t+250, t, 0, 1)",
			[]float64{250, -100, 70}, []float64{90, -280, 250},
			PtXy(64.225, 167.401), PtXy(82.5, 132.5), PtXy(115.225, 102.801),
			VectorIj(65, -220.6).Normalize(), VectorIj(235, -159.4).Normalize(),
		}, {
			PtXy(95, 235), PtXy(46, 25), PtXy(217, 217),
			"Curve(220t^2-98t+95, 402t^2-420t+235, t, 0, 1)",
			[]float64{220, -98, 95}, []float64{402, -420, 235},
			PtXy(86.618, 140.1778), PtXy(101, 125.5), PtXy(128.098, 134.0578),
			VectorIj(47.2, -154.68).Normalize(), VectorIj(196.8, 118.68).Normalize(),
		},
	}
	for h, test := range quadraticTests {
		a := ParamQuadratic(test.p1, test.p2, test.p3)
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		if p := a.PtAtT(0); !IsEqualPair(p, test.p1) {
			t.Errorf("[%d](%s).PtAtT(0) failed. %v != %v",
				h, a, p, test.p1)
		}
		if p := a.PtAtT(0.33); !IsEqualPair(p, test.p33) {
			t.Errorf("[%d](%s).PtAtT(0.33) failed. %v != %v",
				h, a, p, test.p33)
		}
		if p := a.PtAtT(0.50); !IsEqualPair(p, test.p50) {
			t.Errorf("[%d](%s).PtAtT(0.50) failed. %v != %v",
				h, a, p, test.p50)
		}
		if p := a.PtAtT(0.67); !IsEqualPair(p, test.p67) {
			t.Errorf("[%d](%s).PtAtT(0.67) failed. %v != %v",
				h, a, p, test.p67)
		}
		if p := a.PtAtT(2); !IsEqualPair(p, test.p3) {
			t.Errorf("[%d](%s).PtAtT(2) failed. %v != %v",
				h, a, p, test.p3)
		}
		if tt, _ := a.TangentAtT(0.33); !IsEqualPair(tt, test.t33) {
			t.Errorf("[%d](%s).TangentAtT(0.33) failed. %v != %v",
				h, a, tt, test.t33)
		}
		if tt, _ := a.TangentAtT(0.67); !IsEqualPair(tt, test.t67) {
			t.Errorf("[%d](%s).TangentAtT(0.67) failed. %v != %v",
				h, a, tt, test.t67)
		}
	}

	cubicTests := []struct {
		p1, p2, p3, p4 Pt
		s              string
		xc, yc         []float64
		p33, p50, p67  Pt
		t33, t67       Vector
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			"Curve(-85t^3+120t^2+0t+10, -35t^3-75t^2+90t+10, t, 0, 1)",
			[]float64{-85, 120, 0, 10}, []float64{-35, -75, 90, 10},
			PtXy(20.013355, 30.274705), PtXy(29.375, 31.875), PtXy(38.303145, 26.105795),
			VectorIj(51.4305, 29.0655).Normalize(), VectorIj(46.3305, -57.6345).Normalize(),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			"Curve(-740t^3+870t^2+330t-10, -240t^3-1080t^2+1230t-10, t, 0, 1)",
			[]float64{-740, 870, 33, -10}, []float64{-240, 1080, 1230, -10},
			PtXy(167.04962, 269.66312), PtXy(280, 305), PtXy(379.07838, 257.10488),
			VectorIj(662.442, 438.792).Normalize(), VectorIj(499.242, -540.408).Normalize(),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			"Curve(-3.25t^3-6t^2+3.9t-0.1, -1.2t^3-11.4t^2+12.6t-0.1, t, 0, 1)",
			[]float64{-3.25, -6, 3.9, 0.1}, []float64{-1.2, -11.4, 12.6, -0.1},
			PtXy(0.41680475, 2.7734156), PtXy(-0.05625, 3.2), PtXy(-1.15787975, 2.8636244),
			VectorIj(-1.121775, 4.68396).Normalize(), VectorIj(-8.516775, -4.29204).Normalize(),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			"Curve(-202t^3+345t^2-42t+51, 37t^3-396t^2+396t+113, t, 0, 1)",
			[]float64{-202, 345, -42, 51}, []float64{37, -396, 396, 113},
			PtXy(67.451226, 201.885269), PtXy(91, 216.625), PtXy(116.976374, 211.683831),
			VectorIj(119.7066, 146.7279).Normalize(), VectorIj(148.2666, -84.8121).Normalize(),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			"Curve(-455t^3+810t^2-255t+110, -300t^3+60t^2+120t+150, t, 0, 1)",
			[]float64{-455, 810, -255, 110}, []float64{-300, 60, 120, 150},
			PtXy(97.707665, 185.3529), PtXy(128.125, 187.5), PtXy(165.911835, 167.1051),
			VectorIj(130.9515, 61.59).Normalize(), VectorIj(217.6515, -203.61).Normalize(),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			"Curve(-466t^3+1251t^2-921t+396, -479t^3+267t^2+258t+34, t, 0, 1)",
			[]float64{-466, 1251, -821, 396}, []float64{-479, 267, 258, 34},
			PtXy(211.557258, 131.002477), PtXy(190, 169.875), PtXy(200.348342, 182.650823),
			VectorIj(-247.5822, 277.7327).Normalize(), VectorIj(127.7778, -29.2893).Normalize(),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			"Curve(-515t^3+825t^2-468t+285, -232t^3-36t^2+261t+39, t, 0, 1)",
			[]float64{-515, 825, -268, 285}, []float64{-232, -36, 261, 39},
			PtXy(201.894945, 112.872216), PtXy(192.875, 131.5), PtXy(186.889555, 127.932584),
			VectorIj(-91.7505, 161.4456).Normalize(), VectorIj(-56.0505, -99.6744).Normalize(),
		},
	}
	for h, test := range cubicTests {
		a := ParamCubic(test.p1, test.p2, test.p3, test.p4)
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		if p := a.PtAtT(0); !IsEqualPair(p, test.p1) {
			t.Errorf("[%d](%s).PtAtT(0) failed. %v != %v",
				h, a, p, test.p1)
		}
		if p := a.PtAtT(0.33); !IsEqualPair(p, test.p33) {
			t.Errorf("[%d](%s).PtAtT(0.33) failed. %v != %v",
				h, a, p, test.p33)
		}
		if p := a.PtAtT(0.50); !IsEqualPair(p, test.p50) {
			t.Errorf("[%d](%s).PtAtT(0.50) failed. %v != %v",
				h, a, p, test.p50)
		}
		if p := a.PtAtT(0.67); !IsEqualPair(p, test.p67) {
			t.Errorf("[%d](%s).PtAtT(0.67) failed. %v != %v",
				h, a, p, test.p67)
		}
		if p := a.PtAtT(2); !IsEqualPair(p, test.p4) {
			t.Errorf("[%d](%s).PtAtT(2) failed. %v != %v",
				h, a, p, test.p3)
		}
		if tt, _ := a.TangentAtT(0.33); !IsEqualPair(tt, test.t33) {
			t.Errorf("[%d](%s).TangentAtT(0.33) failed. %v != %v",
				h, a, tt, test.t33)
		}
		if tt, _ := a.TangentAtT(0.67); !IsEqualPair(tt, test.t67) {
			t.Errorf("[%d](%s).TangentAtT(0.67) failed. %v != %v",
				h, a, tt, test.t67)
		}
	}

}
