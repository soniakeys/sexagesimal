// Public domain.

package sexa_test

import (
	"fmt"

	"github.com/soniakeys/sexagesimal"
)

func Example_types() {
	// Various construction techniques shown.
	// Similar techniques work across all sexagesimal types.

	v := "%.1v\n" // output shown with common format

	// Angle
	a := sexa.NewAngle(' ', 12, 34, 45.6) // construct from components
	fa := &sexa.FmtAngle{Angle: a}        // Angle is an embedded field
	fmt.Printf(v, fa)

	// HourAngle
	ha := sexa.NewHourAngle('-', 12, 34, 45.6)
	fha := ha.Fmt() // Fmt_ constructor
	fmt.Printf(v, fha)

	// RA
	var fra sexa.FmtRA
	fra.SetHMS(36, 34, 45.6) // RA has no sign, values wrapped to 24 hours
	fmt.Printf(v, &fra)      // custom formatters need pointer receivers

	// Time
	fmt.Printf(v, new(sexa.FmtTime).SetHMS(' ', 12, 34, 45.6))

	// Output:
	// 12°34′45.6″
	// -12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
	// 12ʰ34ᵐ45.6ˢ
}

func Example_verbs() {
	a := sexa.NewAngle(' ', 12, 34, 45.6).Fmt()
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
	a := sexa.NewAngle(' ', 0, 1, 2).Fmt()
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
	a := sexa.NewAngle(' ', 0, 1, 2.34).Fmt()
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
