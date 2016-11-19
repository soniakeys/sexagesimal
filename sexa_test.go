// Public domain.

package sexa_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/soniakeys/sexagesimal"
)

func ExampleInsertUnit() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	fmt.Println("Degree unit with decimal point: ",
		sexa.InsertUnit(formatted, "°"))
	// Output:
	// Decimal point: 1.25
	// Degree unit with decimal point:  1°.25
}

func ExampleCombineUnit() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	// Note that some software may not rendering the combining dot well.
	fmt.Println("Degree unit with combining form of decimal point:",
		sexa.CombineUnit(formatted, "°"))
	// Output:
	// Decimal point: 1.25
	// Degree unit with combining form of decimal point: 1°̣25
}

// For various numbers and symbols, test both Insert and Combine.
// See that the functions do something, and that Strip returns
// the original number.
func TestStrip(t *testing.T) {
	var d string
	var sym string
	t1 := func(fName string, f func(string, string) string) {
		ad := f(d, sym)
		if ad == d {
			t.Fatalf("%s(%s, %s) had no effect", fName, d, sym)
		}
		if sd := sexa.StripUnit(ad, sym); sd != d {
			t.Fatalf("StripUnit(%s, %s) returned %s expected %s",
				ad, sym, sd, d)
		}
	}
	for _, d = range []string{"1.25", "1.", "1", ".25"} {
		for _, sym = range []string{"°", `"`, "h", "ʰ"} {
			t1("InsertUnit", sexa.InsertUnit)
			t1("CombineUnit", sexa.CombineUnit)
		}
	}
}

func ExampleDMSToDeg() {
	// Example p. 7.
	fmt.Printf("%.8f\n", sexa.DMSToDeg(' ', 23, 26, 49))
	// Output:
	// 23.44694444
}

func ExampleNewAngle() {
	// Example negative values, p. 9.
	a := sexa.NewAngle('-', 13, 47, 22)
	fmt.Println(sexa.NewFmtAngle(a.Rad()))
	a = sexa.NewAngle('-', 0, 32, 41)
	// use # flag to force output of all three components
	fmt.Printf("%#s\n", sexa.NewFmtAngle(a.Rad()))
	// Output:
	// -13°47′22″
	// -0°32′41″
}

func ExampleNewRA() {
	// Example 1.a, p. 8.
	a := sexa.NewRA(9, 14, 55.8)
	fmt.Printf("%.6f\n", math.Tan(a.Rad()))
	// Output:
	// -0.877517
}

func ExampleFmtAngle() {
	// Example p. 6
	a := new(sexa.FmtAngle).SetDMS(' ', 23, 26, 44)
	fmt.Println(a)
	// Output:
	// 23°26′44″
}

func ExampleFmtTime() {
	// Example p. 6
	a := new(sexa.FmtTime).SetHMS(' ', 15, 22, 7)
	fmt.Printf("%0s\n", a)
	// Output:
	// 15ʰ22ᵐ07ˢ
}

func TestOverflow(t *testing.T) {
	a := new(sexa.FmtAngle).SetDMS(' ', 23, 26, 44)
	if f := fmt.Sprintf("%03s", a); f != "023°26′44″" {
		t.Fatal(f)
	}
	a.SetDMS(' ', 4423, 26, 44)
	if f := fmt.Sprintf("%03s", a); f != "**********" {
		t.Fatal(f)
	}
}

func TestFmtLeadingZero(t *testing.T) {
	// regression test
	a := sexa.NewFmtAngle(.089876 * math.Pi / 180)
	got := fmt.Sprintf("%.6h", a)
	want := "0.089876°"
	if got != want {
		t.Fatalf("Format %%.6h = %s, want %s", got, want)
	}
}
