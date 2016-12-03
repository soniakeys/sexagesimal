// Public domain.

package sexa_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"
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
	var noSep sexa.Symbols
	got = noSep.CombineUnit(formattedNoDecSep, "°")
	if got != want {
		t.Error("got", got, "want", want)
	}
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
	var noSep sexa.Symbols
	got = noSep.InsertUnit(formattedNoDecSep, "°")
	if got != want {
		t.Error("got", got, "want", want)
	}
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

	// Output:
	// Formatted:    1.25ʰ
	// Strip result: 1.25ʰ false
	// Formatted:    1°25′44.5″
	// Strip result: 1°25′44.5″ false
	// Formatted:    1°25
	// Strip result: 1°25 false
}

func ExampleSymbols_StripUnit_strange() {
	// Empty DecSep.  StripUnit needs to validate the presence of a non-empty
	// separator before it removes the unit.
	formatted := "1°.25"
	fmt.Println("Formatted:   ", formatted)
	var noSep sexa.Symbols
	s, ok := noSep.StripUnit(formatted, "°")
	fmt.Println("Strip result:", s, ok)
	// Output:
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

func ExampleAngle() {
	f := sexa.FmtAngle(unit.NewAngle(' ', 180, 0, 0))
	fmt.Println(f)
	fmt.Printf("%#v\n", *f)
	// Output:
	// 180°0′0″
	// sexa.Angle{Angle:3.141592653589793, Sym:(*sexa.Symbols)(nil), Err:error(nil)}
}

func ExampleFmtAngle() {
	a := unit.NewAngle('-', 13, 47, 22)
	f := sexa.FmtAngle(a)
	fmt.Println(reflect.TypeOf(f), f)
	// Output:
	// *sexa.Angle -13°47′22″
}

func ExampleAngle_String() {
	a := sexa.FmtAngle(unit.NewAngle(' ', 23, 26, 44))
	s := a.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "23°26′44″"
}

func ExampleHourAngle() {
	f := sexa.FmtHourAngle(unit.NewHourAngle(' ', 4, 0, 0))
	fmt.Println(f)
	fmt.Printf("%#v\n", *f)
	// Output:
	// 4ʰ0ᵐ0ˢ
	// sexa.HourAngle{HourAngle:1.0471975511965976, Sym:(*sexa.Symbols)(nil), Err:error(nil)}
}

func ExampleFmtHourAngle() {
	h := unit.NewHourAngle('-', 1, 47, 22)
	f := sexa.FmtHourAngle(h)
	fmt.Println(reflect.TypeOf(f), f)
	// Output:
	// *sexa.HourAngle -1ʰ47ᵐ22ˢ
}

func ExampleHourAngle_String() {
	h := sexa.FmtHourAngle(unit.NewHourAngle('-', 2, 34, 45.6))
	s := h.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "-2ʰ34ᵐ46ˢ"
}

func ExampleRA() {
	f := sexa.FmtRA(unit.NewRA(4, 0, 0))
	fmt.Println(f)
	fmt.Printf("%#v\n", *f)
	// Output:
	// 4ʰ0ᵐ0ˢ
	// sexa.RA{RA:1.0471975511965976, Sym:(*sexa.Symbols)(nil), Err:error(nil)}
}

func ExampleFmtRA() {
	ra := unit.NewRA(1, 47, 22)
	f := sexa.FmtRA(ra)
	fmt.Println(reflect.TypeOf(f), f)
	// Output:
	// *sexa.RA 1ʰ47ᵐ22ˢ
}

func ExampleRA_String() {
	ra := sexa.FmtRA(unit.NewRA(12, 34, 45.6))
	s := ra.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "12ʰ34ᵐ46ˢ"
}

func ExampleTime() {
	t := sexa.FmtTime(unit.NewTime(' ', 0, 1, 30))
	fmt.Println(t)
	fmt.Printf("%#v\n", *t)
	// Output:
	// 1ᵐ30ˢ
	// sexa.Time{Time:90, Sym:(*sexa.Symbols)(nil), Err:error(nil)}
}

func ExampleFmtTime() {
	t := unit.NewTime('-', 15, 22, 7)
	f := sexa.FmtTime(t)
	fmt.Println(reflect.TypeOf(f), f)
	// Output:
	// *sexa.Time -15ʰ22ᵐ7ˢ
}

func ExampleTime_String() {
	t := sexa.FmtTime(unit.NewTime(0, 12, 34, 45.6))
	s := t.String()
	fmt.Printf("%T %q\n", s, s)
	// Output:
	// string "12ʰ34ᵐ46ˢ"
}

func TestOverflow(t *testing.T) {
	a := sexa.FmtAngle(unit.NewAngle(' ', 23, 26, 44))
	if f := fmt.Sprintf("%03s", a); f != " 023°26′44″" {
		t.Fatal(f)
	}
	a.Angle = unit.NewAngle(' ', 4423, 26, 44)
	if f := fmt.Sprintf("%03s", a); f != "***********" {
		t.Fatal(f)
	}
}

func TestLeadingZero(t *testing.T) {
	// regression test
	a := unit.AngleFromDeg(.089876)
	got := fmt.Sprintf("%.6h", sexa.FmtAngle(a))
	want := "0.089876°"
	if got != want {
		t.Fatalf("Format %%.6h = %s, want %s", got, want)
	}
}

func ExampleSymbols_CombineUnit() {
	formatted := "1,25"
	fmt.Println("Decimal comma:", formatted)
	c := sexa.Default
	c.DecSep = ","
	c.DecCombine = '\u0326' // combining comma below
	// Note that some software may not render the combining comma well.
	fmt.Println("Degree unit with combining form of decimal comma:",
		c.CombineUnit(formatted, "°"))
	// Output:
	// Decimal comma: 1,25
	// Degree unit with combining form of decimal comma: 1°̦25
}

func ExampleSymbols_InsertUnit() {
	formatted := "1,25"
	fmt.Println("Decimal comma:", formatted)
	commaSep := &sexa.Symbols{DecSep: ","}
	fmt.Println("Degree unit with decimal comma: ",
		commaSep.InsertUnit(formatted, "°"))
	// Output:
	// Decimal comma: 1,25
	// Degree unit with decimal comma:  1°,25
}

func ExampleSymbols_StripUnit() {
	formatted := "1,25"
	fmt.Println("Decimal comma:", formatted)
	c := &sexa.Symbols{DecSep: ","}
	u := c.InsertUnit(formatted, "°")
	fmt.Println("With degree unit:", u)
	s, ok := c.StripUnit(u, "°")
	fmt.Println("Degree unit stripped:", s, ok)
	// Output:
	// Decimal comma: 1,25
	// With degree unit: 1°,25
	// Degree unit stripped: 1,25 true
}

func ExampleSymbols_FmtAngle() {
	a := unit.NewAngle('-', 13, 47, 22)
	s := sexa.Symbols{DMSUnits: sexa.UnitSymbols{"d ", "m ", "s"}}
	fmt.Println(s.FmtAngle(a))
	// Output:
	// -13d 47m 22s
}

func ExampleSymbols_FmtHourAngle() {
	a := unit.NewHourAngle('-', 1, 47, 22)
	s := sexa.Symbols{HMSUnits: sexa.UnitSymbols{"hr ", "min ", "sec"}}
	fmt.Println(s.FmtHourAngle(a))
	// Output:
	// -1hr 47min 22sec
}

func ExampleSymbols_FmtRA() {
	a := unit.NewRA(1, 47, 22)
	s := sexa.Symbols{HMSUnits: sexa.UnitSymbols{"h", "m", "s"}}
	fmt.Println(s.FmtRA(a))
	// Output:
	// 1h47m22s
}

func ExampleSymbols_FmtTime() {
	a := unit.NewTime('-', 0, 22, 7)
	s := sexa.Symbols{HMSUnits: sexa.UnitSymbols{":", ":", ""}}
	// Careful using ":" for clock-style formatting.  You probably want the
	// '#' flag to show all segments and the '0' flag to show 2 digit minutes
	// and seconds.
	fmt.Printf("%#0s\n", s.FmtTime(a))
	// Output:
	// -0:22:07
}

func TestCoverage(t *testing.T) {
	f := sexa.FmtAngle(unit.AngleFromDeg(9))
	want := "******************"
	got := fmt.Sprintf("%.15h", f)
	if got != want {
		t.Error(got, want)
	}
	if f.Err != sexa.ErrLossOfPrecision {
		t.Error(f.Err, sexa.ErrLossOfPrecision)
	}
	want = "*****************"
	got = fmt.Sprintf("%.14m", f)
	if got != want {
		t.Error(got, want)
	}
	if f.Err != sexa.ErrLossOfPrecision {
		t.Error(f.Err, sexa.ErrLossOfPrecision)
	}

	// flags with non-fixed width
	f.Angle = unit.AngleFromDeg(1.02)
	want = "+1.020°"
	got = fmt.Sprintf("%+.3h", f)
	if got != want {
		t.Error(got, want)
	}
	want = " 1.020°"
	got = fmt.Sprintf("% .3h", f)
	if got != want {
		t.Error(got, want)
	}

	// fixed width
	want = " +1.020°"
	got = fmt.Sprintf("%+2.3h", f)
	if got != want {
		t.Error(got, want)
	}
	f.Angle = -f.Angle
	want = "-01.020°"
	got = fmt.Sprintf("% 02.3h", f)
	if got != want {
		t.Error(got, want)
	}
	f.Angle *= 100
	want = "********"
	got = fmt.Sprintf("% 02.3h", f)
	if got != want {
		t.Error(got, want)
	}
	if f.Err != sexa.ErrDegreeOverflow {
		t.Error(f.Err, sexa.ErrDegreeOverflow)
	}
	tf := sexa.FmtTime(unit.TimeFromHour(102))
	want = "********"
	got = fmt.Sprintf("% 02.3h", tf)
	if got != want {
		t.Error(got, want)
	}
	if tf.Err != sexa.ErrHourOverflow {
		t.Error(tf.Err, sexa.ErrHourOverflow)
	}

	// a coverage test, but the point is that these render with the same width
	// as long as the combining dot is rendered actually combining.
	wantValid := "|  1°̣02|"
	wantOverf := "|******|"
	f.Angle = unit.AngleFromDeg(1.02)
	got = fmt.Sprintf("|%2.2i|", f)
	if got != wantValid {
		t.Error(got, wantValid)
	}
	f.Angle *= 100
	got = fmt.Sprintf("|%2.2i|", f)
	if got != wantOverf {
		t.Error(got, wantOverf)
	}
}
