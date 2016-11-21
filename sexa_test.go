// Public domain.

package sexa_test

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/soniakeys/sexagesimal"
)

func ExampleCombineUnit() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	// Note that some software may not render the combining dot well.
	fmt.Println("Degree unit with combining form of decimal point:",
		sexa.CombineUnit(formatted, "°"))
	// Output:
	// Decimal point: 1.25
	// Degree unit with combining form of decimal point: 1°̣25
}

func ExampleDMSToDeg() {
	// Typical usage:  non-negative d, m, s, and '-' to indicate negative:
	fmt.Println(sexa.DMSToDeg('-', 20, 30, 0))
	// Putting the minus sign on d is not equivalent:
	fmt.Println(sexa.DMSToDeg(' ', -20, 30, 0))
	fmt.Println()
	// Other combinations can give the same result though:
	fmt.Println(sexa.DMSToDeg(' ', -20, -30, 0))
	fmt.Println(sexa.DMSToDeg(' ', -21, 30, 0))
	fmt.Println(sexa.DMSToDeg(' ', -22, 90, 0))
	fmt.Println(sexa.DMSToDeg('-', 22, -90, 0))
	// Output:
	// -20.5
	// -19.5
	//
	// -20.5
	// -20.5
	// -20.5
	// -20.5
}

// see comments below at TestInsertUnit
func TestCombineUnit_No_DecSep(t *testing.T) {
	// test case of DecSep not present
	formattedNoDecSep := "0125"
	got := sexa.CombineUnit(formattedNoDecSep, "°")
	want := "0125°"
	if got != want {
		t.Error("got", got, "want", want)
	}
	// test case of empty DecSep.  same result wanted
	sexa.DecSep = ""
	got = sexa.CombineUnit(formattedNoDecSep, "°")
	if got != want {
		t.Error("got", got, "want", want)
	}
	sexa.DecSep = "." // restore package variable for other tests
}

func ExampleInsertUnit() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	fmt.Println("Degree unit with decimal point: ",
		sexa.InsertUnit(formatted, "°"))
	// Output:
	// Decimal point: 1.25
	// Degree unit with decimal point:  1°.25
}

// For coverage, exercise unusual cases of adding unit with DecSep not present
// in the formatted input or with DecSep set to the empty string.
// An empty DecSep might be useful for formatting numbers in a fixed column
// format with an implicit decimal point.  Maybe four columns with implied
// two decimal places would format 1.25 as 0125.  InsertUnit is documented
// to put the unit symbol at the end in this case, like 0125°.  Really it
// doesn't make much sense that if you are implying the decimal point that
// you wouldn't also be implying the unit symbol, but anyway this is how the
// function is designed to work.
func TestInsertUnit_No_DecSep(t *testing.T) {
	// test case of DecSep not present
	formattedNoDecSep := "0125"
	got := sexa.InsertUnit(formattedNoDecSep, "°")
	want := "0125°"
	if got != want {
		t.Error("got", got, "want", want)
	}
	// test case of empty DecSep.  same result wanted
	sexa.DecSep = ""
	got = sexa.InsertUnit(formattedNoDecSep, "°")
	if got != want {
		t.Error("got", got, "want", want)
	}
	sexa.DecSep = "." // restore package variable for other tests
}

func ExampleStripUnit_combine() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	u := sexa.CombineUnit(formatted, "°")
	// (Note combining dot doesn't display well with some software)
	fmt.Println("With degree unit:", u)
	s, ok := sexa.StripUnit(u, "°")
	fmt.Println("Degree unit stripped:", s, ok)
	// Output:
	// Decimal point: 1.25
	// With degree unit: 1°̣25
	// Degree unit stripped: 1.25 true
}

func ExampleStripUnit_insert() {
	formatted := "1.25"
	fmt.Println("Decimal point:", formatted)
	u := sexa.InsertUnit(formatted, "°")
	fmt.Println("With degree unit:", u)
	s, ok := sexa.StripUnit(u, "°")
	fmt.Println("Degree unit stripped:", s, ok)
	// Output:
	// Decimal point: 1.25
	// With degree unit: 1°.25
	// Degree unit stripped: 1.25 true
}

func ExampleStripUnit_missingUnit() {
	formatted := "1.25"
	fmt.Println("Formatted:   ", formatted)
	s, ok := sexa.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)
	// Output:
	// Formatted:    1.25
	// Strip result: 1.25 false
}

func ExampleStripUnit_strange() {
	// Attempt to strip wrong unit
	formatted := "1.25ʰ"
	fmt.Println("Formatted:   ", formatted)
	s, ok := sexa.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)

	// Multiple segments.  StripUnit isn't meaningful for multiple segments,
	formatted = "1°25′44.5″"
	fmt.Println("Formatted:   ", formatted)
	s, ok = sexa.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)

	// Missing decimal separator.  Not a standard format, unclear how to
	// interpret this, unclear how a result "125" might be interpreted.
	formatted = "1°25"
	fmt.Println("Formatted:   ", formatted)
	s, ok = sexa.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)

	// Empty DecSep.  StripUnit needs to validate the presense of a non-empty
	// separator before it removes the unit.
	formatted = "1°.25"
	fmt.Println("Formatted:   ", formatted)
	sexa.DecSep = ""
	s, ok = sexa.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)
	sexa.DecSep = "."

	// Output:
	// Formatted:    1.25ʰ
	// Strip result: 1.25ʰ false
	// Formatted:    1°25′44.5″
	// Strip result: 1°25′44.5″ false
	// Formatted:    1°25
	// Strip result: 1°25 false
	// Formatted:    1°.25
	// Strip result: 1°.25 false
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
		if sd, ok := sexa.StripUnit(ad, sym); sd != d || !ok {
			t.Fatalf("StripUnit(%s, %s) returned %s false expected %s true",
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

func ExampleNewAngle() {
	a := sexa.NewAngle('-', 13, 47, 22)
	fmt.Println(a.Fmt())
	fmt.Println(sexa.NewAngle('-', 0, 32, 41).Fmt())
	fmt.Println(sexa.NewAngle(' ', 0, 32, 41).Fmt())
	fmt.Println(sexa.NewAngle('+', 0, 32, 41).Fmt())
	fmt.Println(sexa.NewAngle(0, 0, 32, 41).Fmt())
	// Output:
	// -13°47′22″
	// -32′41″
	// 32′41″
	// 32′41″
	// 32′41″
}

func ExampleAngle_Deg() {
	a := sexa.NewAngle(' ', 3, 30, 0)
	fmt.Println(a.Deg())
	// Output:
	// 3.5
}

func ExampleAngle_Fmt() {
	a := sexa.NewAngle(' ', 180, 0, 0)
	f := a.Fmt()
	fmt.Println(f)
	fmt.Printf("%#v\n", *f)
	// Output:
	// 180°0′0″
	// sexa.FmtAngle{Angle:3.141592653589793, Err:error(nil)}
}

func ExampleFmtAngle() {
	f := sexa.NewAngle(' ', 135, 0, 0).Fmt()
	if s := fmt.Sprint(f); f.Err != nil {
		log.Println(f.Err)
	} else {
		fmt.Println(s)
	}
	// Output:
	// 135°0′0″
}

func ExampleFmtAngle_verb() {
	f := sexa.NewAngle(' ', 135, 0, 0).Fmt()
	fmt.Printf("%z\n", f) // produces no output
	if f.Err != nil {
		fmt.Println(f.Err)
	}
	// Output:
	// %!z(BADVERB)
}

func ExampleFmtAngle_width() {
	f := sexa.NewAngle(' ', 135, 0, 0).Fmt()
	fmt.Printf("%2s\n", f) // fills field with *s
	if f.Err != nil {
		fmt.Println(f.Err)
	}
	// Output:
	// *********
	// Degrees overflow width
}

func ExampleAngle_Rad() {
	a := sexa.NewAngle(' ', 180, 0, 0)
	fmt.Println(a.Rad())
	// Output:
	// 3.141592653589793
}

func ExampleFmtAngle_SetDMS() {
	var a sexa.FmtAngle
	a.SetDMS(' ', 23, 26, 44)
	fmt.Println(&a)
	// Output:
	// 23°26′44″
}

func ExampleFmtAngle_String() {
	a := sexa.NewAngle(' ', 23, 26, 44).Fmt()
	s := a.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "23°26′44″"
}

func ExampleNewRA() {
	a := sexa.NewRA(9, 14, 55.8)
	fmt.Printf("%.6f\n", math.Tan(a.Rad()))
	// Output:
	// -0.877517
}

func ExampleFmtHourAngle_String() {
	h := sexa.NewHourAngle('-', 12, 34, 45.6).Fmt()
	s := h.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "-12ʰ34ᵐ46ˢ"
}

func ExampleFmtRA_String() {
	h := new(sexa.FmtRA).SetHMS(12, 34, 45.6)
	s := h.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "12ʰ34ᵐ46ˢ"
}

func ExampleFmtTime() {
	a := new(sexa.FmtTime).SetHMS(' ', 15, 22, 7)
	fmt.Printf("%0s\n", a)
	// Output:
	// 15ʰ22ᵐ07ˢ
}

func ExampleFmtTime_String() {
	h := new(sexa.FmtTime).SetHMS('-', 12, 34, 45.6)
	s := h.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "-12ʰ34ᵐ46ˢ"
}

func ExampleNewTime() {
	t := sexa.NewTime('-', 12, 34, 45.6)
	fmt.Println(t.Fmt())
	// Output:
	// -12ʰ34ᵐ46ˢ
}

func ExampleTime_Sec() {
	t := sexa.NewTime(0, 0, 1, 30)
	fmt.Println(t.Sec())
	// Output:
	// 90
}

func ExampleTime_Min() {
	t := sexa.NewTime(0, 0, 1, 30)
	fmt.Println(t.Min())
	// Output:
	// 1.5
}

func ExampleTime_Hour() {
	t := sexa.NewTime(0, 2, 15, 0)
	fmt.Println(t.Hour())
	// Output:
	// 2.25
}

func ExampleTime_Day() {
	t := sexa.NewTime(0, 48, 36, 0)
	fmt.Println(t.Day())
	// Output:
	// 2.025
}

func ExampleTime_Rad() {
	t := sexa.NewTime(0, 12, 0, 0)
	fmt.Println(t.Rad())
	// Output:
	// 3.141592653589793
}

func TestOverflow(t *testing.T) {
	a := sexa.NewAngle(' ', 23, 26, 44).Fmt()
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
	a := sexa.Angle(.089876 * math.Pi / 180)
	got := fmt.Sprintf("%.6h", a.Fmt())
	want := "0.089876°"
	if got != want {
		t.Fatalf("Format %%.6h = %s, want %s", got, want)
	}
}
