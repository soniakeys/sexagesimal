// License: MIT

package sexa_test

import (
	"fmt"
	"math"

	"github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"
)

func Example_flags() {
	a := sexa.FmtAngle(unit.NewAngle(' ', 0, 1, 2))
	fmt.Printf("%+s\n", a)    // + sign for non-negative
	fmt.Printf("% s\n", a)    // space for non-negative
	fmt.Printf("%0s\n", a)    // 0 pad segments after the first
	fmt.Printf("%#s\n", a)    // output all segments
	fmt.Printf("%+ #0s\n", a) // all four flags
	// Output:
	// +1′2″
	//  1′2″
	// 1′02″
	// 0°1′2″
	// +0°01′02″
}

// Various construction techniques are shown below.  Similar value construction
// techniques work across all sexagesimal types.
func Example_types() {
	fs := "%.1s\n" // a common format string for all examples

	// Angle
	a := unit.NewAngle(' ', 12, 34, 45.6) // construct unit.Angle from components
	fa := &sexa.Angle{Angle: a}           // unit.Angle is embedded in sexa.Angle
	fmt.Printf(fs, fa)

	// HourAngle
	ha := unit.NewHourAngle('-', 12, 34, 45.6)
	fha := sexa.FmtHourAngle(ha) // all four types also have Fmt- constructor
	fmt.Printf(fs, fha)

	// RA
	var fra sexa.RA // zero value of type
	// (unit.RA has no sign, values wrapped to 24 hours)
	fra.RA = unit.NewRA(36, 34, 45.6) // assign RA field of existing sexa.RA
	fmt.Printf(fs, &fra)              // custom formatters need pointer receivers

	// Time
	fmt.Printf(fs, sexa.FmtTime(unit.NewTime(' ', 12, 34, 45.6))) // one-liner

	// Output:
	// 12°34′45.6″
	// -12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
}

// Comments give some possible mnemonics
func Example_verbs() {
	a := sexa.FmtAngle(unit.NewAngle(' ', 12, 34, 45.6))
	fmt.Println("Full sexagesimal formats")
	fmt.Printf("%.1s\n", a) // (s)exagesimal format, decimal point in (s)econds
	fmt.Printf("%.1c\n", a) // (c)ombining form
	fmt.Printf("%.1d\n", a) // (d)ecimal point following unit symbol
	fmt.Println("Decimal minute formats")
	fmt.Printf("%.2m\n", a) // decimal (m)inute
	fmt.Printf("%.2n\n", a) // m(n)o, just alphabetic order
	fmt.Printf("%.2o\n", a) // mn(o)
	fmt.Println("Decimal hour/degree formats")
	fmt.Printf("%.3h\n", a) // decimal (h)our/degree
	fmt.Printf("%.3i\n", a) // h(i)j, alphabetic order
	fmt.Printf("%.3j\n", a) // hi(j)
	// (Remember the combining form may not format well especially with
	// monospace fonts.)

	// Output:
	// Full sexagesimal formats
	// 12°34′45.6″
	// 12°34′45″̣6
	// 12°34′45″.6
	// Decimal minute formats
	// 12°34.76′
	// 12°34′̣76
	// 12°34′.76
	// Decimal hour/degree formats
	// 12.579°
	// 12°̣579
	// 12°.579
}

func Example_width() {
	// fixed width format
	p := sexa.FmtAngle(unit.NewAngle(' ', 0, 1, 2.34))
	n := sexa.FmtAngle(unit.NewAngle('-', 0, 1, 2.34))
	fmt.Printf("|%2.3s|\n", p)
	fmt.Printf("|%2.3s|\n", n)

	// same with '0' flag
	fmt.Println()
	fmt.Printf("|%02.3s|\n", p)
	fmt.Printf("|%02.3s|\n", n)

	// fixed width with no unit symbols or decimal separators packs columns
	fmt.Println()
	ff := func(a unit.Angle) {
		noSep := sexa.Symbols{}
		fmt.Printf("|%+02.3s|\n", noSep.FmtAngle(a))
	}
	ff(0)
	ff(unit.NewAngle('-', 0, 1, 2.34))
	ff(unit.NewAngle(' ', 23, 45, 16.7))

	// no width specifier but additional printf step to format into fixed
	// width string.  '0' flag aligns unit symbols.
	fmt.Println()
	ff = func(a unit.Angle) {
		fmt.Printf("|%15s|\n", fmt.Sprintf("%0.3s", sexa.FmtAngle(a)))
	}
	ff(0)
	ff(unit.NewAngle('-', 0, 1, 2.34))
	ff(unit.NewAngle('-', 123, 45, 16.7))

	// same with '#' flag
	fmt.Println()
	ff = func(a unit.Angle) {
		fmt.Printf("|%15s|\n", fmt.Sprintf("%#0.3s", sexa.FmtAngle(a)))
	}
	ff(0)
	ff(unit.NewAngle('-', 0, 1, 2.34))
	ff(unit.NewAngle('-', 123, 45, 16.7))

	// Output:
	// |  0° 1′ 2.340″|
	// |- 0° 1′ 2.340″|
	//
	// | 00°01′02.340″|
	// |-00°01′02.340″|
	//
	// |+000000000|
	// |-000102340|
	// |+234516700|
	//
	// |         0.000″|
	// |     -1′02.340″|
	// |-123°45′16.700″|
	//
	// |   0°00′00.000″|
	// |  -0°01′02.340″|
	// |-123°45′16.700″|
}

func Example_withOverflow() {
	// Angle example
	af := sexa.FmtAngle(unit.NewAngle(' ', 35, 0, 0))
	fmt.Printf("|%2s|\n", af) // 35 fits in two digits
	af = sexa.FmtAngle(unit.NewAngle(' ', 135, 0, 0))
	// 135 doesn't fit.  Yes there's another space there but that's for
	// the sign.  The specified width for the degrees field is 2.
	fmt.Printf("|%2s|\n", af)
	fmt.Println("Err:", af.Err)

	// Time example
	tf := sexa.FmtTime(unit.NewTime(' ', 12, 0, 0))
	fmt.Printf("\n|%2m|\n", tf) // 12 fits in two digits
	tf = sexa.FmtTime(unit.NewTime(' ', 125, 0, 0))
	fmt.Printf("|%2m|\n", tf) // 125 doesn't.
	fmt.Println("Err:", tf.Err)
	// Output:
	// | 35° 0′ 0″|
	// |**********|
	// Err: Degrees overflow width
	//
	// | 12ʰ 0ᵐ|
	// |*******|
	// Err: Hours overflow width
}

func Example_withInvalidVerb() {
	f := sexa.FmtAngle(unit.NewAngle(' ', 135, 0, 0))
	fmt.Printf("\nFormatted: %q\n", f)
	fmt.Println("Err:", f.Err)
	// Output:
	// Formatted: %!q(BADVERB)
	// Err: <nil>
}

func Example_withLossOfPrecision() {
	f := sexa.FmtAngle(unit.NewAngle(' ', 135, 0, 0))
	fmt.Printf("%.16s\n", f) // 16 is always too much to ask
	fmt.Printf("%.10s ", f)  // but 10 is still too much for 135°
	fmt.Println(f.Err)
	fmt.Printf("%.9s\n", f) // 9 is ok.  all digits are significant.
	// Output:
	// %!(BADPREC 16)
	// ************* Loss of precision
	// 135°0′0.000000000″
}

func Example_withInfNaN() {
	f := sexa.FmtAngle(unit.Angle(math.Inf(1)))
	fmt.Printf("%s", f)
	fmt.Println("", f.Err)
	f.Angle = unit.Angle(math.Inf(-1))
	fmt.Printf("%s", f)
	fmt.Println("", f.Err)
	f.Angle = unit.Angle(math.NaN())
	fmt.Printf("%s", f)
	fmt.Println("", f.Err)
	// Output:
	// ** +Inf
	// ** -Inf
	// ** NaN
}
