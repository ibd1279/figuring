package figuring

import "fmt"

// Circle represents a geometric circle defined as a center point and a raidus.
type Circle struct {
	c Pt
	r Length
}

// CirclePt creates a circle at specific point.
func CirclePt(c Pt, r Length) Circle {
	if r < 0 {
		r = -r
	}
	return Circle{
		c: c,
		r: r,
	}
}

// BoundingBox returns the bounding box for the circle.
func (c Circle) BoundingBox() Rectangle {
	v := VectorIj(c.r, c.r)
	least, most := c.c.Add(v), c.c.Add(v.Invert())
	return RectanglePt(least, most)
}

// OrErr returns a floating point error if either the center or the radius are
// in error.
func (c Circle) OrErr() (Circle, *FloatingPointError) {
	_, cerr := c.c.OrErr()
	_, rerr := c.r.OrErr()
	if cerr != nil && cerr.IsNaN() {
		return c, cerr
	} else if rerr != nil && rerr.IsNaN() {
		return c, rerr
	} else if cerr != nil {
		return c, cerr
	} else if rerr != nil {
		return c, rerr
	}
	return c, nil
}

// PtAtTheta returns the point on the circle, at the provided angle.
func (c Circle) PtAtTheta(theta Radians) Pt {
	v := VectorFromTheta(theta).Scale(c.r)
	return c.c.Add(v)
}

// String returns the implicit formula of this circle.
func (c Circle) String() string {
	x, y := c.c.XY()
	r := c.r
	xop, yop := '-', '-'

	if x < 0 {
		xop = '+'
		x = -x
	}
	if y < 0 {
		yop = '+'
		y = -y
	}
	return fmt.Sprintf("(x%c%s)^2+(y%c%s)^2=%s^2",
		xop,
		HumanFormat(9, x),
		yop,
		HumanFormat(9, y),
		HumanFormat(9, r),
	)
}
