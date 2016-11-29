// Public domain.

package sexa_test

import (
	"fmt"

	"github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"
)

func Example_types() {
	// Various construction techniques shown.
	// Similar techniques work across all sexagesimal types.

	v := "%.1v\n" // output shown with common format

	// Angle
	a := unit.NewAngle(' ', 12, 34, 45.6) // construct from components
	fa := &sexa.Angle{Angle: a}           // Angle is an embedded field
	fmt.Printf(v, fa)

	// HourAngle
	ha := unit.NewHourAngle('-', 12, 34, 45.6)
	fha := sexa.FmtHourAngle(ha) // _ constructor
	fmt.Printf(v, fha)

	// RA
	var fra sexa.RA
	// RA has no sign, values wrapped to 24 hours
	fra.RA = unit.NewRA(36, 34, 45.6)
	fmt.Printf(v, &fra) // custom formatters need pointer receivers

	// Time
	fmt.Printf(v, sexa.FmtTime(unit.NewTime(' ', 12, 34, 45.6)))

	// Output:
	// 12°34′45.6″
	// -12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
}

func Example_verbs() {
	a := sexa.FmtAngle(unit.NewAngle(' ', 12, 34, 45.6))
	fmt.Println("Full sexagesimal formats")
	fmt.Printf("%.1s\n", a)
	fmt.Printf("%.1c\n", a)
	fmt.Printf("%.1d\n", a)
	fmt.Println("Decimal minute formats")
	fmt.Printf("%.2m\n", a)
	fmt.Printf("%.2n\n", a)
	fmt.Printf("%.2o\n", a)
	fmt.Println("Decimal hour/degree formats")
	fmt.Printf("%.3h\n", a)
	fmt.Printf("%.3i\n", a)
	fmt.Printf("%.3j\n", a)
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

func Example_flags() {
	a := sexa.FmtAngle(unit.NewAngle(' ', 0, 1, 2))
	fmt.Printf("%+s\n", a)
	fmt.Printf("% s\n", a)
	fmt.Printf("%#s\n", a)
	fmt.Printf("%0s\n", a)
	fmt.Printf("%+ #0s\n", a)
	// Output:
	// +1′2″
	//  1′2″
	// 0°1′2″
	// 1′02″
	// +0°01′02″
}

func Example_width() {
	// fixed width formats
	a := sexa.FmtAngle(unit.NewAngle(' ', 0, 1, 2.34))
	fmt.Printf("|%2.3s|\n", a)
	fmt.Printf("|%02.3s|\n", a)

	// example with implicit units and decimal position.

	// save package values, restore on function return.
	defer func(u sexa.UnitSymbols, d string) {
		sexa.DMSUnits = u
		sexa.DecSep = d
	}(sexa.DMSUnits, sexa.DecSep)

	// empty package variables.
	sexa.DMSUnits = sexa.UnitSymbols{}
	sexa.DecSep = ""
	fmt.Printf("\n|%02.3s|\n", a)

	// Output:
	// | 0° 1′ 2.340″|
	// |00°01′02.340″|
	//
	// |000102340|
}
