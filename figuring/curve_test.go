package figuring

import (
	"log"
	"math"
	"testing"
)

func TestParamCurve(t *testing.T) {
	linearTests := []struct {
		p1, p2   Pt
		s        string
		p33, p67 Pt
	}{
		{
			PtXy(0, 10), PtXy(20, 15), "Curve(20t+0, 5t+10, t, 0, 1)",
			PtXy(6.6, 11.65), PtXy(13.4, 13.35),
		}, {
			PtXy(-10, -10), PtXy(-30, -5), "Curve(-20t-10, 5t-10, t, 0, 1)",
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
		p33, p50, p67 Pt
		t33, t67      Vector
	}{
		{
			PtXy(70, 250), PtXy(20, 110), PtXy(220, 60),
			"Curve(250t^2-100t+70, 90t^2-280t+250, t, 0, 1)",
			PtXy(64.225, 167.401), PtXy(82.5, 132.5), PtXy(115.225, 102.801),
			VectorIj(65, -220.6), VectorIj(235, -159.4),
		}, {
			PtXy(95, 235), PtXy(46, 25), PtXy(217, 217),
			"Curve(220t^2-98t+95, 402t^2-420t+235, t, 0, 1)",
			PtXy(86.618, 140.1778), PtXy(101, 125.5), PtXy(128.098, 134.0578),
			VectorIj(47.2, -154.68), VectorIj(196.8, 118.68),
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
		p33, p50, p67  Pt
		t33, t67       Vector
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			"Curve(-85t^3+120t^2+0t+10, -35t^3-75t^2+90t+10, t, 0, 1)",
			PtXy(20.013355, 30.274705), PtXy(29.375, 31.875), PtXy(38.303145, 26.105795),
			VectorIj(51.4305, 29.0655), VectorIj(46.3305, -57.6345),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			"Curve(-740t^3+870t^2+330t-10, -240t^3-1080t^2+1230t-10, t, 0, 1)",
			PtXy(167.04962, 269.66312), PtXy(280, 305), PtXy(379.07838, 257.10488),
			VectorIj(662.442, 438.792), VectorIj(499.242, -540.408),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			"Curve(-3.25t^3-6t^2+3.9t-0.1, -1.2t^3-11.4t^2+12.6t-0.1, t, 0, 1)",
			PtXy(0.41680475, 2.7734156), PtXy(-0.05625, 3.2), PtXy(-1.15787975, 2.8636244),
			VectorIj(-1.121775, 4.68396), VectorIj(-8.516775, -4.29204),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			"Curve(-202t^3+345t^2-42t+51, 37t^3-396t^2+396t+113, t, 0, 1)",
			PtXy(67.451226, 201.885269), PtXy(91, 216.625), PtXy(116.976374, 211.683831),
			VectorIj(119.7066, 146.7279), VectorIj(148.2666, -84.8121),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			"Curve(-455t^3+810t^2-255t+110, -300t^3+60t^2+120t+150, t, 0, 1)",
			PtXy(97.707665, 185.3529), PtXy(128.125, 187.5), PtXy(165.911835, 167.1051),
			VectorIj(130.9515, 61.59), VectorIj(217.6515, -203.61),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			"Curve(-466t^3+1251t^2-921t+396, -479t^3+267t^2+258t+34, t, 0, 1)",
			PtXy(211.557258, 131.002477), PtXy(190, 169.875), PtXy(200.348342, 182.650823),
			VectorIj(-247.5822, 277.7327), VectorIj(127.7778, -29.2893),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			"Curve(-515t^3+825t^2-468t+285, -232t^3-36t^2+261t+39, t, 0, 1)",
			PtXy(201.894945, 112.872216), PtXy(192.875, 131.5), PtXy(186.889555, 127.932584),
			VectorIj(-91.7505, 161.4456), VectorIj(-56.0505, -99.6744),
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

	quarticTests := []struct {
		p1, p2, p3, p4, p5 Pt
		s                  string
		p33, p50, p67      Pt
		t33, t67           Vector
	}{
		{
			PtXy(-2.42, -8.24), PtXy(-0.14, -2.94), PtXy(5.74, -8.84),
			PtXy(9.96, 0.4), PtXy(13.78, -5.2),
			"Curve(6.52t^4-21.04t^3+21.6t^2+9.12t-2.42, -56.32t^4+105.36t^3-67.2t^2+21.2t-8.24, t, 0, 1)",
			PtXy(2.2630475692, -5.4436683872), PtXy(5.3175, -4.79), PtXy(8.3724395692, -3.8628016672),
			VectorIj(17.43946896, 3.17322464),
			VectorIj(17.57333104, 5.28442336),
		},
	}
	for h, test := range quarticTests {
		a := ParamQuartic(test.p1, test.p2, test.p3, test.p4, test.p5)
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
		if p := a.PtAtT(2); !IsEqualPair(p, test.p5) {
			t.Errorf("[%d](%s).PtAtT(2) failed. %v != %v",
				h, a, p, test.p5)
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

	lengthTests := []struct {
		a          ParamCurve
		asteps     int
		al, l      Length
		trimTo, tl Length
	}{
		{
			ParamLinear(PtXy(1, 1), PtXy(1, 11)),
			16, 10.0, 10.0,
			5., 5.,
		}, {
			ParamCubic(PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10)),
			16, 81.79, 81.7889377631191,
			65., 65.,
		}, {
			ParamCubic(PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80)),
			16, 306.21, 306.2137924899652,
			100., 100.,
		}, {
			ParamQuartic(PtXy(-2.42, -8.24), PtXy(-0.14, -2.94), PtXy(5.74, -8.84),
				PtXy(9.96, 0.4), PtXy(13.78, -5.2)),
			16, 18.18, 18.1824891479602,
			5., 5.,
		},
	}
	for h, test := range lengthTests {
		a := test.a
		length := a.ApproxLength(test.asteps)
		if !IsEqual(length.Round(), test.al.Round()) {
			t.Errorf("[%d](%s).ApproxLength(%d) failed. %v != %v",
				h, a, test.asteps, length.Round(), test.al.Round())
		}
		length = a.Length()
		if !IsEqual(length, test.l) {
			t.Errorf("[%d](%s).Length() failed. %v != %v",
				h, a, length, test.l)
		}
		b, _ := a.SplitAtLength(test.trimTo)
		length = b.Length() * 10
		tl := test.tl * 10
		if !IsEqual(length.Round(), tl.Round()) {
			t.Errorf("[%d](%s).TrimToLength() failed. %v != %v",
				h, a, length, test.tl)
		}
	}

	// Bounding Box
	boxTesting := []struct {
		a   ParamCurve
		box Rectangle
	}{
		{
			ParamQuartic(
				PtXy(-2.42, -8.24), PtXy(-0.14, -2.94),
				PtXy(5.74, -8.84), PtXy(9.96, 0.4),
				PtXy(13.78, -5.2)),
			RectanglePt(PtXy(-2.42, -8.24), PtXy(13.78, -3.409209141)),
		}, {
			ParamCubic(
				PtXy(10, 10), PtXy(10, 40),
				PtXy(50, 45), PtXy(45, -10)),
			RectanglePt(PtXy(10, -10), PtXy(45.432526, 32.126252)),
		}, {
			ParamCubic(
				PtXy(-10, -10), PtXy(100, 400),
				PtXy(500, 450), PtXy(450, -100)),
			RectanglePt(PtXy(-10, -100), PtXy(454.303137, 305.156522)),
		}, {
			ParamQuadratic(
				PtXy(-0.10, -0.10), PtXy(0.5, 4.50),
				PtXy(-5.45, -0.1)),
			RectanglePt(PtXy(-5.45, -0.1), PtXy(-0.045038168, 2.2)),
		},
	}
	for h, test := range boxTesting {
		a := test.a
		box := a.BoundingBox()
		if !IsEqualPts(box, test.box) {
			t.Errorf("[%d](%s).BoundingBox() failed. %v != %v",
				h, a, box, test.box)
		}
	}

	// Roots
	rootsTests := []struct {
		p1, p2, p3, p4 Pt
		xroots, yroots []float64
	}{
		{
			//0
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			[]float64{0.3706095719294, 0},
			[]float64{1, 0},
		}, {
			//1
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			[]float64{0},
			[]float64{1, 0},
		}, {
			//2
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			[]float64{0.5094282273212036, 0},
			[]float64{1, 0},
		}, {
			//3
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			[]float64{0},
			[]float64{1, 0},
		}, {
			//4
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			[]float64{0.5846511999394, 0},
			[]float64{1, 0.2198586446693, 0},
		}, {
			//5
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			[]float64{0},
			[]float64{1, 0.0840609753683, 0},
		}, {
			//6
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			[]float64{0},
			[]float64{1, 0},
		},
	}
	for h, test := range rootsTests {
		line := LineFromPt(test.p1, test.p4)
		pts := []Pt{test.p1, test.p2, test.p3, test.p4}
		pts = RotateOrTranslateToXAxis(line, pts)
		pts = TranslatePts(pts[0].VectorTo(PtOrig), pts)
		a := ParamCubic(pts[0], pts[1], pts[2], pts[3])
		xroots, yroots := a.Roots()
		if len(xroots) != len(test.xroots) {
			t.Fatalf("[%d](%s).Roots() (X) (length) failed. %v != %v",
				h, a.X, xroots, test.xroots)
		}
		for i := 0; i < len(xroots); i++ {
			if !IsEqual(xroots[i], test.xroots[i]) {
				t.Errorf("[%d][%d](%s).Roots() (X) failed. %v != %v",
					h, i, a.X, xroots[i], test.xroots[i])
			}
		}

		if len(yroots) != len(test.yroots) {
			t.Fatalf("[%d](%s).Roots() (Y) (length) failed. %v != %v",
				h, a.Y, yroots, test.yroots)
		}
		for i := 0; i < len(yroots); i++ {
			if !IsEqual(yroots[i], test.yroots[i]) {
				t.Errorf("[%d][%d](%s).Roots() (Y) failed. %v != %v",
					h, i, a.Y, yroots[i], test.yroots[i])
			}
		}
	}
}

func TestBezier(t *testing.T) {
	// Identity (points, string, p33, p66)
	identityTests := []struct {
		p1, p2, p3, p4 Pt
		s              string
		p33, p50, p67  Pt
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			"Bezier[ Curve(-85t^3+120t^2+0t+10, -35t^3-75t^2+90t+10, t, 0, 1) ]",
			PtXy(20.013355, 30.274705), PtXy(29.375, 31.875), PtXy(38.303145, 26.105795),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			"Bezier[ Curve(-740t^3+870t^2+330t-10, -240t^3-1080t^2+1230t-10, t, 0, 1) ]",
			PtXy(167.04962, 269.66312), PtXy(280, 305), PtXy(379.07838, 257.10488),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			"Bezier[ Curve(-3.25t^3-6t^2+3.9t-0.1, -1.2t^3-11.4t^2+12.6t-0.1, t, 0, 1) ]",
			PtXy(0.41680475, 2.7734156), PtXy(-0.05625, 3.2), PtXy(-1.15787975, 2.8636244),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			"Bezier[ Curve(-202t^3+345t^2-42t+51, 37t^3-396t^2+396t+113, t, 0, 1) ]",
			PtXy(67.451226, 201.885269), PtXy(91, 216.625), PtXy(116.976374, 211.683831),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			"Bezier[ Curve(-455t^3+810t^2-255t+110, -300t^3+60t^2+120t+150, t, 0, 1) ]",
			PtXy(97.707665, 185.3529), PtXy(128.125, 187.5), PtXy(165.911835, 167.1051),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			"Bezier[ Curve(-466t^3+1251t^2-921t+396, -479t^3+267t^2+258t+34, t, 0, 1) ]",
			PtXy(211.557258, 131.002477), PtXy(190, 169.875), PtXy(200.348342, 182.650823),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			"Bezier[ Curve(-515t^3+825t^2-468t+285, -232t^3-36t^2+261t+39, t, 0, 1) ]",
			PtXy(201.894945, 112.872216), PtXy(192.875, 131.5), PtXy(186.889555, 127.932584),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			"Bezier[ Curve(455t^3-450t^2+150t+70, -410t^3+945t^2-705t+250, t, 0, 1) ]",
			PtXy(86.846335, 105.52633), PtXy(89.375, 82.5), PtXy(105.342165, 78.54767),
		},
	}
	for h, test := range identityTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		if p := a.PtAtT(0); !IsEqualPair(p, test.p1) {
			t.Errorf("[%d](%s).PtAtT(0) failed. %v != %v",
				h, a, p, test.p1)
		}
		if p := a.Begin(); !IsEqualPair(p, test.p1) {
			t.Errorf("[%d](%s).Begin() failed. %v != %v",
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
		if p := a.PtAtT(1); !IsEqualPair(p, test.p4) {
			t.Errorf("[%d](%s).PtAtT(1) failed. %v != %v",
				h, a, p, test.p3)
		}
		if p := a.End(); !IsEqualPair(p, test.p4) {
			t.Errorf("[%d](%s).End() failed. %v != %v",
				h, a, p, test.p3)
		}
	}

	// Flattening
	// TODO depends on polygon / linear spline

	// Splitting
	splittingTests := []struct {
		p1, p2, p3, p4    Pt
		t33left, t33right Bezier
		t67left, t67right Bezier
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			BezierPt(PtXy(10, 10), PtXy(10, 19.9),
				PtXy(14.356, 27.0775), PtXy(20.013355, 30.274705)),
			BezierPt(PtXy(20.013355, 30.274705), PtXy(31.4995, 36.766),
				PtXy(48.35, 26.85), PtXy(45, -10)),
			BezierPt(PtXy(10, 10), PtXy(10, 30.1),
				PtXy(27.956, 38.9775), PtXy(38.303145, 26.105795)),
			BezierPt(PtXy(38.303145, 26.105795), PtXy(43.3995, 19.766),
				PtXy(46.65, 8.15), PtXy(45, -10)),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			BezierPt(PtXy(-10, -10), PtXy(26.3, 125.3),
				PtXy(94.181, 221.396), PtXy(167.04962, 269.66312)),
			BezierPt(PtXy(167.04962, 269.66312), PtXy(314.995, 367.66),
				PtXy(483.5, 268.5), PtXy(450, -100)),
			BezierPt(PtXy(-10, -10), PtXy(63.7, 264.7),
				PtXy(267.581, 377.796), PtXy(379.07838, 257.10488)),
			BezierPt(PtXy(379.07838, 257.10488), PtXy(433.995, 197.66),
				PtXy(466.5, 81.5), PtXy(450, -100)),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			BezierPt(PtXy(-0.1, -0.1), PtXy(0.329, 1.286),
				PtXy(0.5402, 2.25818), PtXy(0.41680475, 2.7734156)),
			BezierPt(PtXy(0.41680475, 2.7734156), PtXy(0.166275, 3.8195),
				PtXy(-1.4635, 2.982), PtXy(-5.45, -0.1)),
			BezierPt(PtXy(-0.1, -0.1), PtXy(0.771, 2.714),
				PtXy(0.7442, 3.82218), PtXy(-1.15787975, 2.8636244)),
			BezierPt(PtXy(-1.15787975, 2.8636244), PtXy(-2.094725, 2.3915),
				PtXy(-3.4865, 1.418), PtXy(-5.45, -0.1)),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			BezierPt(PtXy(51, 113), PtXy(46.38, 156.56),
				PtXy(54.2835, 185.7452), PtXy(67.451226, 201.885269)),
			BezierPt(PtXy(67.451226, 201.885269), PtXy(94.1857, 234.6545),
				PtXy(142.62, 213.65), PtXy(152, 150)),
			BezierPt(PtXy(51, 113), PtXy(41.62, 201.44),
				PtXy(83.8635, 230.6252), PtXy(116.976374, 211.683831)),
			BezierPt(PtXy(116.976374, 211.683831), PtXy(133.2857, 202.3545),
				PtXy(147.38, 181.35), PtXy(152, 150)),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			BezierPt(PtXy(110, 150), PtXy(81.95, 163.2),
				PtXy(83.303, 178.578), PtXy(97.707665, 185.3529)),
			BezierPt(PtXy(97.707665, 185.3529), PtXy(126.9535, 199.108),
				PtXy(210, 177.4), PtXy(210, 30)),
			BezierPt(PtXy(110, 150), PtXy(53.05, 176.8),
				PtXy(117.303, 212.578), PtXy(165.911835, 167.1051)),
			BezierPt(PtXy(165.911835, 167.1051), PtXy(189.8535, 144.708),
				PtXy(210, 102.6), PtXy(210, 30)),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			BezierPt(PtXy(396, 34), PtXy(294.69, 62.38),
				PtXy(238.7913, 100.4521), PtXy(211.557258, 131.002477)),
			BezierPt(PtXy(211.557258, 131.002477), PtXy(156.2639, 193.029),
				PtXy(219.13, 224.05), PtXy(260, 80)),
			BezierPt(PtXy(396, 34), PtXy(190.31, 91.62),
				PtXy(171.8113, 189.1921), PtXy(200.348342, 182.650823)),
			BezierPt(PtXy(200.348342, 182.650823), PtXy(214.4039, 179.429),
				PtXy(239.87, 150.95), PtXy(260, 80)),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			BezierPt(PtXy(285, 39), PtXy(233.52, 67.71),
				PtXy(211.9875, 95.1132), PtXy(201.894945, 112.872216)),
			BezierPt(PtXy(201.894945, 112.872216), PtXy(181.404, 148.9284),
				PtXy(208.07, 145.23), PtXy(127, 32)),
			BezierPt(PtXy(285, 39), PtXy(180.48, 97.29),
				PtXy(199.4075, 150.1932), PtXy(186.889555, 127.932584)),
			BezierPt(PtXy(186.889555, 127.932584), PtXy(180.724, 116.9684),
				PtXy(166.93, 87.77), PtXy(127, 32)),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			BezierPt(PtXy(70, 250), PtXy(86.5, 172.45),
				PtXy(86.665, 129.2035), PtXy(86.846335, 105.52633)),
			BezierPt(PtXy(86.846335, 105.52633), PtXy(87.2145, 57.4545),
				PtXy(87.65, 90.05), PtXy(225, 80)),
			BezierPt(PtXy(70, 250), PtXy(103.5, 92.55),
				PtXy(69.665, 76.5035), PtXy(105.342165, 78.54767)),
			BezierPt(PtXy(105.342165, 78.54767), PtXy(122.9145, 79.5545),
				PtXy(157.35, 84.95), PtXy(225, 80)),
		},
	}
	for h, test := range splittingTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		left, right := a.SplitAtT(0.33)
		if !IsEqualPts(left, test.t33left) || !IsEqualPts(right, test.t33right) {
			t.Errorf("[%d](%s).SplitAtT(0.33) failed. %v != %v || %v != %v",
				h, a, left, test.t33left, right, test.t33right)
			log.Printf("0.33\n%+v\n%+v\n\n%+v\n%+v",
				left.pts, test.t33left.pts, right.pts, test.t33right.pts)
		}
		left, right = a.SplitAtT(0.67)
		if !IsEqualPts(left, test.t67left) || !IsEqualPts(right, test.t67right) {
			t.Errorf("[%d](%s).SplitAtT(0.67) failed. %v != %v || %v != %v",
				h, a, left, test.t67left, right, test.t67right)
			log.Printf("0.67\n%+v\n%+v\n\n%+v\n%+v",
				left.pts, test.t67left.pts, right.pts, test.t67right.pts)
		}
	}

	// Tangents and Normals
	tangentTests := []struct {
		p1, p2, p3, p4 Pt
		t33, n33       Vector
		t67, n67       Vector
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			VectorIj(51.4305, 29.0655), VectorIj(-29.0655, 51.4305),
			VectorIj(46.3305, -57.6345), VectorIj(57.6345, 46.3305),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			VectorIj(662.442, 438.792), VectorIj(-438.792, 662.442),
			VectorIj(499.242, -540.408), VectorIj(540.408, 499.242),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			VectorIj(-1.121775, 4.68396), VectorIj(-4.68396, -1.121775),
			VectorIj(-8.516775, -4.29204), VectorIj(4.29204, -8.516775),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			VectorIj(119.7066, 146.7279), VectorIj(-146.7279, 119.7066),
			VectorIj(148.2666, -84.8121), VectorIj(84.8121, 148.2666),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			VectorIj(130.9515, 61.59), VectorIj(-61.59, 130.9515),
			VectorIj(217.6515, -203.61), VectorIj(203.61, 217.6515),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			VectorIj(-247.5822, 277.7307), VectorIj(-277.7307, -247.5822),
			VectorIj(127.7778, -29.2893), VectorIj(29.2893, 127.7778),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			VectorIj(-91.7505, 161.4456), VectorIj(-161.4456, -91.7505),
			VectorIj(-56.0505, -99.6744), VectorIj(99.6744, -56.0505),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			VectorIj(1.6485, -215.247), VectorIj(215.247, 1.6485),
			VectorIj(159.7485, 9.153), VectorIj(-9.153, 159.7485),
		},
	}
	for h, test := range tangentTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		tangent, normal := a.TangentAtT(0.33)
		if !IsEqualPair(tangent, test.t33) || !IsEqualPair(normal, test.n33) {
			t.Errorf("[%d](%s).TangentAtT(0.33) failed. \n%v != \n%v || \n%v != \n%v",
				h, a, tangent, test.t33, normal, test.n33)
		}
		tangent, normal = a.TangentAtT(0.67)
		if !IsEqualPair(tangent, test.t67) || !IsEqualPair(normal, test.n67) {
			t.Errorf("[%d](%s).TangentAtT(0.67) failed. \n%v != \n%v || \n%v != \n%v",
				h, a, tangent, test.t67, normal, test.n67)
		}
	}

	// Bounding Box
	boxTesting := []struct {
		p1, p2, p3, p4 Pt
		box            Rectangle
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			RectanglePt(PtXy(10, -10), PtXy(45.432526, 32.126252)),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			RectanglePt(PtXy(-10, -100), PtXy(454.303137, 305.156522)),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			RectanglePt(PtXy(-5.45, -0.1), PtXy(0.451705, 3.201703)),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			RectanglePt(PtXy(49.672082, 113), PtXy(152, 217.192920)),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			RectanglePt(PtXy(87.6645332689289, 30), PtXy(210, 188.8623458218187)),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			RectanglePt(PtXy(189.825127, 34), PtXy(396, 182.963675)),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			RectanglePt(PtXy(127, 32), PtXy(285, 133.130906)),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			RectanglePt(PtXy(70, 78.391972622), PtXy(225, 250)),
		},
	}
	for h, test := range boxTesting {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		box := a.BoundingBox()
		if !IsEqualPts(box, test.box) {
			t.Errorf("[%d](%s).BoundingBox() failed. %v != %v",
				h, a, box, test.box)
		}
	}

	// Aligning
	aligningTests := []struct {
		p1, p2, p3, p4 Pt
		trans          Vector
		theta          Radians
		scale          Length
		ax             Bezier
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			VectorIj(-10, -10), -5.764037121173873, 40.311288741,
			BezierPt(PtXy(0, 0), PtXy(-0.369230769, 0.646153846),
				PtXy(0.430769231, 1.246153846), PtXy(1, 0)),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			VectorIj(10, 10), -1.938498874567 * math.Pi, 468.72166581,
			BezierPt(PtXy(0, 0), PtXy(0.062357761, 0.903504779),
				PtXy(0.879380974, 1.172052799), PtXy(1, 0)),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			VectorIj(0.10, 0.10), -1 * math.Pi, 5.35,
			BezierPt(PtXy(0, 0), PtXy(-0.242990654, -0.785046729),
				PtXy(-0.112149533, -0.859813084), PtXy(1, 0)),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			VectorIj(-51, -113), -0.1117757396712 * math.Pi, 107.563934476,
			BezierPt(PtXy(0, 0), PtXy(0.29991357, 1.197061366),
				PtXy(1.18159032, 0.874070873), PtXy(1, 0)),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			VectorIj(-110, -150), -1.7211420616237 * math.Pi, 156.204993518,
			BezierPt(PtXy(0, 0), PtXy(-0.545081967, -0.254098361),
				PtXy(-0.081967213, 0.901639344), PtXy(1, 0)),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			VectorIj(-396, -34), -0.8961813805266 * math.Pi, 143.568798839,
			BezierPt(PtXy(0, 0), PtXy(2.217543179, 0.117698428),
				PtXy(1.882301572, -1.282456821), PtXy(1, 0)),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			VectorIj(-285, -39), -1.0140931207676 * math.Pi, 158.154987275,
			BezierPt(PtXy(0, 0), PtXy(0.961060249, -0.59321153),
				PtXy(0.188382041, -1.033662496), PtXy(1, 0)),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			VectorIj(-70, -250), -1.7353191928108 * math.Pi, 230.054341407,
			BezierPt(PtXy(0, 0), PtXy(0.90127539, -0.527633444),
				PtXy(0.351440718, -0.61454889), PtXy(1, 0)),
		},
	}
	for h, test := range aligningTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		trans, theta, scale, aligned := a.AlignOnX()
		if !IsEqualPair(trans, test.trans) {
			t.Errorf("[%d](%s).AlignOnX() (translate) failed. %v != %v",
				h, a, trans, test.trans)
		}
		if !IsEqual(theta, test.theta) {
			t.Errorf("[%d](%s).AlignOnX() (angle) failed. %v != %v",
				h, a, theta, test.theta)
		}
		if !IsEqual(scale, test.scale) {
			t.Errorf("[%d](%s).AlignOnX() (scale) failed. %v != %v",
				h, a, scale, test.scale)
		}
		if !IsEqualPts(aligned, test.ax) {
			t.Errorf("[%d](%s).AlignOnX() failed. %v != %v / %+v != %+v",
				h, a, aligned, test.ax, aligned.Points(), test.ax.Points())
		}
	}

	// Roots
	rootsTests := []struct {
		p1, p2, p3, p4 Pt
		xroots, yroots []float64
	}{
		{
			//0
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			[]float64{0.3706095719294, 0},
			[]float64{1, 0},
		}, {
			//1
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			[]float64{0},
			[]float64{1, 0},
		}, {
			//2
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			[]float64{0.5094282273212036, 0},
			[]float64{1, 0},
		}, {
			//3
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			[]float64{0},
			[]float64{1, 0},
		}, {
			//4
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			[]float64{0.5846511999394, 0},
			[]float64{1, 0.2198586446693, 0},
		}, {
			//5
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			[]float64{0},
			[]float64{1, 0.0840609753683, 0},
		}, {
			//6
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			[]float64{0},
			[]float64{1, 0},
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			[]float64{0},
			[]float64{1, 0},
		},
	}
	for h, test := range rootsTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		_, _, _, b := a.AlignOnX()
		xroots, yroots := b.Roots()
		if len(xroots) != len(test.xroots) {
			t.Fatalf("[%d](%s).Roots() (X) (length) failed. %v != %v",
				h, b.x, xroots, test.xroots)
		}
		for i := 0; i < len(xroots); i++ {
			if !IsEqual(xroots[i], test.xroots[i]) {
				t.Errorf("[%d][%d](%s).Roots() (X) failed. %v != %v",
					h, i, b.x, xroots[i], test.xroots[i])
			}
		}

		if len(yroots) != len(test.yroots) {
			t.Fatalf("[%d](%s).Roots() (Y) (length) failed. %v != %v",
				h, b.y, yroots, test.yroots)
		}
		for i := 0; i < len(yroots); i++ {
			if !IsEqual(yroots[i], test.yroots[i]) {
				t.Errorf("[%d][%d](%s).Roots() (Y) failed. %v != %v",
					h, i, b.y, yroots[i], test.yroots[i])
			}
		}
	}

	// TightBox
	tightBoxTests := []struct {
		p1, p2, p3, p4 Pt
		box            Polygon
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			PolygonPt(
				PtXy(6.712347988, 11.878658293),
				PtXy(45, -10),
				PtXy(59.532851496, 15.432490118),
				PtXy(21.245199484, 37.311148411),
			),
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			PolygonPt(
				PtXy(-10, -10),
				PtXy(450, -100),
				PtXy(520.340831688, 259.519806406),
				PtXy(60.340831688, 349.519806406),
			),
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			PolygonPt(
				PtXy(0.451704996, 3.201702789),
				PtXy(-5.45, 3.201702789),
				PtXy(-5.45, -0.1),
				PtXy(0.451704996, -0.1),
			),
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			PolygonPt(
				PtXy(51, 113),
				PtXy(154.464544376, 150.90285289),
				PtXy(125.554932261, 229.818280554),
				PtXy(22.090387886, 191.915427664),
			),
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			PolygonPt(
				PtXy(82.568952748, 173.783261538),
				PtXy(205.50787123, 26.256559359),
				PtXy(251.989423465, 64.991186221),
				PtXy(129.050504983, 212.5178884),
			),
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			PolygonPt(
				PtXy(421.056718144, 108.080731903),
				PtXy(193.102651037, 185.182842836),
				PtXy(167.718784596, 110.13488988),
				PtXy(395.672851703, 33.032778947),
			),
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			PolygonPt(
				PtXy(280.65385548, 137.098690599),
				PtXy(122.65385548, 130.098690599),
				PtXy(127, 32),
				PtXy(285, 39),
			),
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			PolygonPt(
				PtXy(-2.919229782, 183.514819905),
				PtXy(152.080770218, 13.514819905),
				PtXy(225, 80),
				PtXy(70, 250),
			),
		},
	}
	for h, test := range tightBoxTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		box := a.TightBox()
		if !IsEqualPts(box, test.box) {
			t.Errorf("[%d](%s).TightBox() failed. %v != %v",
				h, a, box, test.box)
		}
	}

	// Inflections
	inflectionTests := []struct {
		p1, p2, p3, p4 Pt
		inflects       []float64
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			nil,
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			nil,
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			nil,
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			nil,
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			nil,
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			nil,
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			[]float64{0.43807908584189087, 0.7193516086422476},
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			[]float64{0.32665059013775993, 0.7295669472896766},
		},
	}
	for h, test := range inflectionTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		pts := a.InflectionPts()
		if len(pts) != len(test.inflects) {
			t.Fatalf("[%d](%s).InflectionPts() (length) failed. %v != %v",
				h, a, pts, test.inflects)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqual(pts[i], test.inflects[i]) {
				t.Errorf("[%d][%d](%s).InflectionPts() failed. %v != %v",
					h, i, a, pts[i], test.inflects[i])
			}
		}
	}

	// Curve Type
	curvetypeTests := []struct {
		p1, p2, p3, p4 Pt
		curvetype      BezierCurveType
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			BEZIER_CURVE_TYPE_PLAIN,
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			BEZIER_CURVE_TYPE_PLAIN,
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			BEZIER_CURVE_TYPE_PLAIN,
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			BEZIER_CURVE_TYPE_PLAIN,
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			BEZIER_CURVE_TYPE_PLAIN,
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			BEZIER_CURVE_TYPE_LOOP,
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			BEZIER_CURVE_TYPE_DOUBLEINFLECTION,
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			BEZIER_CURVE_TYPE_DOUBLEINFLECTION,
		},
	}
	for h, test := range curvetypeTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		curvetype := a.CurveType()
		if curvetype != test.curvetype {
			t.Errorf("[%d](%s).CurveType() failed. %d != %d",
				h, a, curvetype, test.curvetype)
		}
	}

	// Finding Y given X
	// TODO Not sure of the value of this for my application.

	// Arc Length
	lengthTests := []struct {
		p1, p2, p3, p4       Pt
		length, approxLength Length
	}{
		{
			PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10),
			81.7889377631191, 81.79,
		}, {
			PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100),
			944.927455012432, 944.93,
		}, {
			PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1),
			10.0199019804689, 10.02,
		}, {
			PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150),
			221.2661140122809, 221.27,
		}, {
			PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30),
			272.8700297821006, 272.87,
		}, {
			PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80),
			398.1352944481076, 398.14,
		}, {
			PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32),
			255.8115465025523, 255.81,
		}, {
			PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80),
			306.2137924899652, 306.21,
		},
	}

	for h, test := range lengthTests {
		a := BezierPt(test.p1, test.p2, test.p3, test.p4)
		length := a.Length()
		if !IsEqual(length, test.length) {
			t.Errorf("[%d](%s).Length() failed. %f != %f",
				h, a, length, test.length)
		}
		length = a.ApproxLength(32)
		if length.Round() != test.approxLength.Round() {
			t.Errorf("[%d](%s).ApproxLength() failed. %f != %f",
				h, a, length, test.approxLength)
		}
	}

	// Distance Intervals
	// TODO Not sure how I want this to work.
}

func BenchmarkBezierLength(b *testing.B) {
	lengthTests := []Bezier{
		BezierPt(PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10)),
		BezierPt(PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100)),
		BezierPt(PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1)),
		BezierPt(PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150)),
		BezierPt(PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30)),
		BezierPt(PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80)),
		BezierPt(PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32)),
		BezierPt(PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80)),
	}
	max := len(lengthTests)
	for h := 0; h < b.N; h++ {
		lengthTests[h%max].Length()
	}
}
func BenchmarkBezierApproxLength(b *testing.B) {
	lengthTests := []Bezier{
		BezierPt(PtXy(10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10)),
		BezierPt(PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100)),
		BezierPt(PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1)),
		BezierPt(PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150)),
		BezierPt(PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30)),
		BezierPt(PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80)),
		BezierPt(PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32)),
		BezierPt(PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80)),
	}
	max := len(lengthTests)
	for h := 0; h < b.N; h++ {
		lengthTests[h%max].ApproxLength(16)
	}
}
