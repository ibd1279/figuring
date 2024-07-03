/*
Package figuring is a math and geometry library used for doing 2D vector
manipulations. It includes functions for Length, Radians, Vectors, Points,
Lines, Curves, and Polygons.
*/
package figuring

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/mathgl/mgl64"
)

// These unit of measures values were copied from https://www.npl.co.uk/si-units
const (
	Micrometer Length = 1                  // Micrometer (µm) unit of measure
	Millimeter        = 1000               // Millimeter (mm) unit of measure
	Centimeter        = 10000              // Centimeter (cm) unit of measure
	Decimeter         = 100000             // Decimeter (dm) unit of measure
	Meter             = 1000000            // Meter (m) unit of measure
	Dekameter         = 10 * Meter         // Dekameter (dam) unit of measure
	Hectometer        = 100 * Meter        // Hectometer (hm) unit of measure
	Kilometer         = 1000 * Meter       // Kilometer (km) unit of measure
	Megameter         = 1000000 * Meter    // Megameter (Mm) unit of measure
	Gigameter         = 1000000000 * Meter // Gigameter (Gm) unit of measure
)

// These unit of measure labels were copied from https://www.npl.co.uk/si-units
const (
	MicrometerLabel     string = "µm"
	MicrometerLabelAlt1        = "um"
	MicrometerLabelAlt2        = "micro"
	MillimeterLabel            = "mm"
	MillimeterLabelAlt1        = "milli"
	CentimeterLabel            = "cm"
	CentimeterLabelAlt1        = "centi"
	DecimeterLabel             = "dm"
	DecimeterLabelAlt1         = "deci"
	MeterLabel                 = "m"
	MeterLabelAlt1             = "meter"
	DekameterLabel             = "dam"
	DekameterLabelAlt1         = "deka"
	HectometerLabel            = "hm"
	HectometerLabelAlt1        = "hecto"
	KilometerLabel             = "km"
	KilometerLabelAlt1         = "kilo"
	MegameterLabel             = "Mm"
	MegameterLabelAlt1         = "mega"
	GigameterLabel             = "Gm"
	GigameterLabelAlt1         = "giga"
)

const (
	// equalEpsilon is used by the equals methods to compare floats.
	// differences less than this are considered equal
	equalEpsilon = 1e-5

	// zeroEpsilon is used to check some values against zero.
	zeroEpsilon = 1e-9
)

var (
	uomParseLabels = map[string]Length{
		MicrometerLabel:     Micrometer,
		MicrometerLabelAlt1: Micrometer,
		MicrometerLabelAlt2: Micrometer,
		MillimeterLabel:     Millimeter,
		MillimeterLabelAlt1: Millimeter,
		CentimeterLabel:     Centimeter,
		CentimeterLabelAlt1: Centimeter,
		DecimeterLabel:      Decimeter,
		DecimeterLabelAlt1:  Decimeter,
		MeterLabel:          Meter,
		MeterLabelAlt1:      Meter,
		DekameterLabel:      Dekameter,
		DekameterLabelAlt1:  Dekameter,
		HectometerLabel:     Hectometer,
		HectometerLabelAlt1: Hectometer,
		KilometerLabel:      Kilometer,
		KilometerLabelAlt1:  Kilometer,
		MegameterLabel:      Megameter,
		MegameterLabelAlt1:  Megameter,
		GigameterLabel:      Gigameter,
		GigameterLabelAlt1:  Gigameter,
	}
)

// Radians is used for angle measurements.
type Radians float64

// RadiansFromDegrees creates a Radian value from a degrees value.
func RadiansFromDegrees(f float64) Radians { return Radians(f * math.Pi / 180) }

// RadiansFromFloat creates and normalizes a Radian value to be betwee 0 <= r <
// 2*math.Pi
func RadiansFromFloat(f float64) Radians { return Radians(f).Normalize() }

// Degrees create a degree value from a radian value.
func (r Radians) Degrees() float64 { return float64(r) * 180 / math.Pi }

// Normalize the radians to between 0 <= r < 2*math.Pi
func (r Radians) Normalize() Radians {
	n := Radians(math.Mod(float64(r), 2*math.Pi))
	if Signbit(n) {
		n = Radians(2*math.Pi) + n
	}
	if IsEqual(n, Radians(2*math.Pi)) {
		n = 0
	}
	return n
}

// String for outputting.
func (r Radians) String() string { return fmt.Sprintf("θ(%0.5fπ)", float64(r/math.Pi)) }

// OrErr tests if the value is a NaN or Inf value and returns an error if it is.
func (r Radians) OrErr() (Radians, *FloatingPointError) {
	f := float64(r)
	if math.IsNaN(f) || math.IsInf(f, -1) || math.IsInf(f, 1) {
		return r, &FloatingPointError{v: f}
	}
	return r, nil
}

// Length is used for distance measurements.  The value is based on a float64
// and supports around 9 Gm before percision starts to degrade.
type Length float64

// ParseUnitOfMeasure parses a string UOM value and returns the UOM constant.
// Falls back to the \c fallback UOM if the label cannot be matched or the
// number cannot be parsed.
func ParseUnitOfMeasure(s string, fallback Length) Length {
	if uom, ok := uomParseLabels[s]; ok {
		return uom
	}
	return fallback
}

// LengthUom creates a distance based on \c f of \c uom.
func LengthUom(f float64, uom Length) Length { return Length(f) * uom }

// Round a length value to the closest micrometer.
func (d Length) Round() Length { return Length(math.Round(float64(d))) }

// Float returns \c d as a float64, in the given UOM.
func (d Length) Float(uom Length) float64 { return float64(d / uom) }

// Float32 returns \c d as a float32, in the given UOM.
func (d Length) Float32(uom Length) float32 { return float32(d / uom) }

// Int returns \c d as a int, in the given UOM.
func (d Length) Int(uom Length) int { return int(d / uom) }

// Int64 returns \c d as an int64, in the given UOM.
func (d Length) Int64(uom Length) int64 { return int64(d / uom) }

// OrErr tests if a length is a NaN or Inf value and returns an error if it is.
func (d Length) OrErr() (Length, *FloatingPointError) {
	f := float64(d)
	if math.IsNaN(f) || math.IsInf(f, -1) || math.IsInf(f, 1) {
		return d, &FloatingPointError{v: f}
	}
	return d, nil
}

// HumanUnitLabel returns the UOM and label that would make the most sense for
// the value.
func (d Length) HumanUnitLabel() (Length, string) {
	u := d + 1
	if u > Gigameter {
		return Gigameter, GigameterLabel
	} else if u > Megameter {
		return Megameter, MegameterLabel
	} else if u > Kilometer {
		return Kilometer, KilometerLabel
	} else if u > Meter {
		return Meter, MeterLabel
	} else if u > Millimeter {
		return Millimeter, MillimeterLabel
	}
	return Micrometer, MicrometerLabel
}

// Text generates a human readable string for the length. Includes the UOM
// label.
func (d Length) Text(uom Length) string {
	uom, label := uom.HumanUnitLabel()
	return fmt.Sprintf("%0.03f%s", d.Float(uom), label)
}

// FloatingPointError provides an error interfaced wrapper for floats.
type FloatingPointError struct {
	v float64
}

// Error implements the error interface.
func (e *FloatingPointError) Error() string {
	if math.IsNaN(e.v) {
		return "NaN encountered"
	}
	if math.IsInf(e.v, -1) {
		return "Negative Inf encountered"
	}
	if math.IsInf(e.v, 1) {
		return "Positive Inf encountered"
	}
	return fmt.Sprintf("%g resulted in an error", e.v)
}

// IsNaN tests if the error was because of a NaN value.
func (e *FloatingPointError) IsNaN() bool { return math.IsNaN(e.v) }

// IsInf tests if the error was because of a Inf value, positive or negative.
func (e *FloatingPointError) IsInf() bool { return math.IsInf(e.v, 0) }

// IsPosInf tests if the error was because of a positive Inf value.
func (e *FloatingPointError) IsPosInf() bool { return math.IsInf(e.v, 1) }

// IsNegInf tests if the error was because of a negative Inf value.
func (e *FloatingPointError) IsNegInf() bool { return math.IsInf(e.v, -1) }

// Minimum returns the smallest value from a set of values. Discards NaN values.
func Minimum[T Radians | Length | float64](vals ...T) (ret T) {
	if len(vals) < 1 {
		return ret
	}

	ret = vals[0]
	for _, v := range vals {
		if v < ret || math.IsNaN(float64(ret)) {
			ret = v
		}
	}
	return ret
}

// Maximum returns the largest value from a set of values. Discards NaN values.
func Maximum[T Radians | Length | float64](vals ...T) (ret T) {
	if len(vals) < 1 {
		return ret
	}

	ret = vals[0]
	for _, v := range vals {
		if v > ret || math.IsNaN(float64(ret)) {
			ret = v
		}
	}
	return ret
}

// Clamp value v between min and max. Preserves NaN values.
func Clamp[T Radians | Length | float64](min, v, max T) T {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// IsEqual tests if two values are within a tolerance of each other.
func IsEqual[T Radians | Length | float64](a, b T) bool {
	return mgl64.FloatEqualThreshold(float64(a), float64(b), equalEpsilon)
}

// IsZero tests if a value is within a tolerance of zero.
func IsZero[T Radians | Length | float64](a T) bool {
	if -zeroEpsilon < a && a < zeroEpsilon {
		return true
	}
	return false
}

// Signbit tests if the (negative) sign bit is set on a value.
func Signbit[T Radians | Length | float64](a T) bool { return math.Signbit(float64(a)) }

// HumanFormat outputs the floating point value with the desired percision.
// Trailing zeros are trimmed.
func HumanFormat[T Radians | Length | float64](percision int, v T) string {
	fmtstr := fmt.Sprintf("%%.%df", percision)
	str := fmt.Sprintf(fmtstr, v)
	idx := strings.LastIndexAny(str, "123456789.")
	if idx > -1 {
		str = str[:idx+1]
	}
	if strings.HasSuffix(str, ".") {
		str = str[:len(str)-1]
	}
	return str
}
