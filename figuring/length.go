package figuring

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	// These unit of measures values were copied from https://www.npl.co.uk/si-units
	UOM_MICROMETER Length = 1                      // Micrometer(µm) unit of measure
	UOM_MILLIMETER        = 1000                   // Millimeter(mm) unit of measure
	UOM_CENTIMETER        = 10000                  // Centimeter(cm) unit of measure
	UOM_DECIMETER         = 100000                 // Decimeter(dm) unit of measure
	UOM_METER             = 1000000                // Meter (m) unit of measure
	UOM_DEKAMETER         = 10 * UOM_METER         // Dekameter(dam) unit of measure
	UOM_HECTOMETER        = 100 * UOM_METER        // Hectometer(hm) unit of measure
	UOM_KILOMETER         = 1000 * UOM_METER       // Kilometer(km) unit of measure
	UOM_MEGAMETER         = 1000000 * UOM_METER    // Megameter(Mm) unit of measure
	UOM_GIGAMETER         = 1000000000 * UOM_METER // Gigameter(Gm) unit of measure

	// These unit of measure labels were copied from https://www.npl.co.uk/si-units
	UOM_MICROMETER_LABEL      string = "µm"
	UOM_MICROMETER_LABEL_ALT1        = "um"
	UOM_MICROMETER_LABEL_ALT2        = "micro"
	UOM_MILLIMETER_LABEL             = "mm"
	UOM_MILLIMETER_LABEL_ALT1        = "milli"
	UOM_CENTIMETER_LABEL             = "cm"
	UOM_CENTIMETER_LABEL_ALT1        = "centi"
	UOM_DECIMETER_LABEL              = "dm"
	UOM_DECIMETER_LABEL_ALT1         = "deci"
	UOM_METER_LABEL                  = "m"
	UOM_METER_LABEL_ALT1             = "meter"
	UOM_DEKAMETER_LABEL              = "dam"
	UOM_DEKAMETER_LABEL_ALT1         = "deka"
	UOM_HECTOMETER_LABEL             = "hm"
	UOM_HECTOMETER_LABEL_ALT1        = "hecto"
	UOM_KILOMETER_LABEL              = "km"
	UOM_KILOMETER_LABEL_ALT1         = "kilo"
	UOM_MEGAMETER_LABEL              = "Mm"
	UOM_MEGAMETER_LABEL_ALT1         = "mega"
	UOM_GIGAMETER_LABEL              = "Gm"
	UOM_GIGAMETER_LABEL_ALT1         = "giga"

	UOM_DEFAULT = UOM_MICROMETER
)

const (
	// floatTOLERANCE is used by the equals methods to compare floats.
	// differences less than this are considered equal
	floatTOLERANCE = 1e-5

	// floatZERO_TOLERANCE is used to check some values against zero.
	floatZERO_TOLERANCE = 1e-9

	// floatHUMAN_DENOMINATOR is used to decide when to switch up a level
	// for UOM.
	lengthHUMAN_DENOMINATOR = 1.9
)

var (
	uomParseLabels map[string]Length = map[string]Length{
		UOM_MICROMETER_LABEL:      UOM_MICROMETER,
		UOM_MICROMETER_LABEL_ALT1: UOM_MICROMETER,
		UOM_MICROMETER_LABEL_ALT2: UOM_MICROMETER,
		UOM_MILLIMETER_LABEL:      UOM_MILLIMETER,
		UOM_MILLIMETER_LABEL_ALT1: UOM_MILLIMETER,
		UOM_CENTIMETER_LABEL:      UOM_CENTIMETER,
		UOM_CENTIMETER_LABEL_ALT1: UOM_CENTIMETER,
		UOM_DECIMETER_LABEL:       UOM_DECIMETER,
		UOM_DECIMETER_LABEL_ALT1:  UOM_DECIMETER,
		UOM_METER_LABEL:           UOM_METER,
		UOM_METER_LABEL_ALT1:      UOM_METER,
		UOM_DEKAMETER_LABEL:       UOM_DEKAMETER,
		UOM_DEKAMETER_LABEL_ALT1:  UOM_DEKAMETER,
		UOM_HECTOMETER_LABEL:      UOM_HECTOMETER,
		UOM_HECTOMETER_LABEL_ALT1: UOM_HECTOMETER,
		UOM_KILOMETER_LABEL:       UOM_KILOMETER,
		UOM_KILOMETER_LABEL_ALT1:  UOM_KILOMETER,
		UOM_MEGAMETER_LABEL:       UOM_MEGAMETER,
		UOM_MEGAMETER_LABEL_ALT1:  UOM_MEGAMETER,
		UOM_GIGAMETER_LABEL:       UOM_GIGAMETER,
		UOM_GIGAMETER_LABEL_ALT1:  UOM_GIGAMETER,
	}
)

// Radians is used for angle measurements.
type Radians float64

// Create a Radian value from a degrees value.
func RadiansFromDegrees(f float64) Radians { return Radians(f * math.Pi / 180) }

// Create and normalize a Radian value to be betwee 0 <= r < 2*math.Pi
func RadiansFromFloat(f float64) Radians { return Radians(f).Normalize() }

// Create a Degree value from a radian value.
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

// Parse a string UOM value and return the UOM constant.  Falls back to the \c
// fallback UOM if the label cannot be matched or the number cannot be parsed.
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

// HumanUnitLabel returns the UOM and label that would make the most sense for the value.
func (d Length) HumanUnitLabel() (Length, string) {
	u := d + 1
	if u > UOM_GIGAMETER {
		return UOM_GIGAMETER, UOM_GIGAMETER_LABEL
	} else if u > UOM_MEGAMETER {
		return UOM_MEGAMETER, UOM_MEGAMETER_LABEL
	} else if u > UOM_KILOMETER {
		return UOM_KILOMETER, UOM_KILOMETER_LABEL
	} else if u > UOM_METER {
		return UOM_METER, UOM_METER_LABEL
	} else if u > UOM_MILLIMETER {
		return UOM_MILLIMETER, UOM_MILLIMETER_LABEL
	}
	return UOM_MICROMETER, UOM_MICROMETER_LABEL
}

// String generates a human readable string for the length. Includes the UOM label.
func (d Length) Text(uom Length) string {
	uom, label := uom.HumanUnitLabel()
	return fmt.Sprintf("%0.03f%s", d.Float(uom), label)
}

type FloatingPointError struct {
	v float64
}

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

func (e *FloatingPointError) IsNaN() bool    { return math.IsNaN(e.v) }
func (e *FloatingPointError) IsInf() bool    { return math.IsInf(e.v, 0) }
func (e *FloatingPointError) IsPosInf() bool { return math.IsInf(e.v, 1) }
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

// ClampRadians clamps value v between min and max.
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
	return mgl64.FloatEqualThreshold(float64(a), float64(b), floatTOLERANCE)
}

// IsZero tests if a value is within a tolerance of zero.
func IsZero[T Radians | Length | float64](a T) bool {
	if -floatZERO_TOLERANCE < a && a < floatZERO_TOLERANCE {
		return true
	}
	return false
}

// Signbit tests if the (negative) sign bit is set on a value.
func Signbit[T Radians | Length | float64](a T) bool { return math.Signbit(float64(a)) }

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
