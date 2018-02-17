// License: MIT

package sexa

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/soniakeys/unit"
)

// Predefined errors indicate that a value could not be formatted.
// Custom formatters of Angle, HourAngle, RA, and Time types
// may store these in the Err field of the value being formatted.
var (
	ErrLossOfPrecision = errors.New("Loss of precision")
	ErrDegreeOverflow  = errors.New("Degrees overflow width")
	ErrHourOverflow    = errors.New("Hours overflow width")
	ErrPosInf          = errors.New("+Inf")
	ErrNegInf          = errors.New("-Inf")
	ErrNaN             = errors.New("NaN")
)

// UnitSymbols holds symbols for formatting Angle, HourAngle, RA,
// and Time types.
type UnitSymbols struct {
	HrDeg, Min, Sec string
}

// Symbols species unit and decimal indicators
//
// Set these as needed, for example to ASCII symbols.
// It is valid to use multiple character strings for DMSUnits and HMSUnits.
// It is valid to use empty strings with a fixed width format.
// DecCombine should be a rune of Unicode category "Mn" (mark, nonspacing).
type Symbols struct {
	DMSUnits   UnitSymbols
	HMSUnits   UnitSymbols
	DecSep     string
	DecCombine rune
}

// Default symbols are used by package top-level functions.
var Default = &Symbols{
	DMSUnits:   UnitSymbols{"°", "′", "″"},
	HMSUnits:   UnitSymbols{"ʰ", "ᵐ", "ˢ"},
	DecSep:     ".",
	DecCombine: '\u0323',
}

// CombineUnit inserts a unit indicator into a formatted decimal number,
// combining it if possible with the decimal separator.
//
// The decimal separator is identified by the package variable Default.
// If Default.DecSep is non-empty and occurrs in d, the occurrence is replaced
// with argument 'unit' and the symbol Default.DecCombine.  Otherwise unit is
// appended to the end of d.
//
// See also InsertUnit, StripUnit, and Symbols.CombineUnit.
func CombineUnit(d, unit string) string {
	return Default.CombineUnit(d, unit)
}

// InsertUnit inserts a unit indicator into a formatted decimal number.
//
// The indicator is inserted just before the decimal separator if one is
// present, or at the end of the number otherwise.
//
// The decimal separator is identified by the package variable Default.
// If Default.DecSep is non-empty and occurrs in d, unit is added just before
// the occurrence.  Otherwise unit is appended to the end of d.
//
// See also CombineUnit, StripUnit, and Symbols.InsertUnit.
func InsertUnit(d, unit string) string {
	return Default.InsertUnit(d, unit)
}

// StripUnit reverses the action of InsertUnit or CombineUnit,
// removing the specified unit indicator and restoring a following
// DecCombine to DecSep.
//
// More specifically, if the specified unit is found at the end of string d
// or if DecSep is non-empty and ocurrs immediately following, then the unit
// is Removed.  If the specified unit is found with DecCombine immediately
// following, then both the unit and the DecCombine rune are replaced with
// DecSep.  As with InsertUnit and CombineUnit, the decimal separators are
// identified by the package variable Default.
//
// StripUnit returns ok = true if the unit was found and removed.  Otherwise it
// returns d unchanged and ok = false.
func StripUnit(d, unit string) (stripped string, ok bool) {
	return Default.StripUnit(d, unit)
}

// Angle is represents a formattable angle.
type Angle struct {
	unit.Angle
	Sym *Symbols
	Err error // set each time the value is formatted.
}

// FmtAngle constructs an formattable Angle containing the value a.
func FmtAngle(a unit.Angle) *Angle { return &Angle{Angle: a} }

// Format implements fmt.Formatter
func (a *Angle) Format(f fmt.State, c rune) {
	s := state{
		State:  f,
		verb:   c,
		hrDeg:  a.Deg(),
		caller: fsAngle,
		sym:    a.Sym,
	}
	a.Err = s.writeFormatted()
}

// String implements fmt.Stringer
func (a *Angle) String() string { return fmt.Sprintf("%s", a) }

// HourAngle represents a formattable angle hour.
type HourAngle struct {
	unit.HourAngle
	Sym *Symbols
	Err error // set each time the value is formatted.
}

// FmtHourAngle constructs an formattable HourAngle containing the value h.
func FmtHourAngle(h unit.HourAngle) *HourAngle {
	return &HourAngle{HourAngle: h}
}

// Format implements fmt.Formatter
func (ha *HourAngle) Format(f fmt.State, c rune) {
	s := &state{
		State:  f,
		verb:   c,
		hrDeg:  ha.Hour(),
		caller: fsHourAngle,
		sym:    ha.Sym,
	}
	ha.Err = s.writeFormatted()
}

// String implements fmt.Stringer
func (ha *HourAngle) String() string { return fmt.Sprintf("%s", ha) }

// RA represents a formattable right ascension.
type RA struct {
	unit.RA
	Sym *Symbols
	Err error // set each time the value is formatted.
}

// FmtRA constructs an formattable RA containing the value ra.
func FmtRA(ra unit.RA) *RA { return &RA{RA: ra} }

// Format implements fmt.Formatter, formatting to hours, minutes, and seconds.
func (ra *RA) Format(f fmt.State, c rune) {
	s := &state{
		State: f,
		verb:  c,
		// pmod in case ra.RA was directly set to something out of range
		hrDeg:  unit.PMod(ra.Hour(), 24),
		caller: fsRA,
		sym:    ra.Sym,
	}
	ra.Err = s.writeFormatted()
}

// String implements fmt.Stringer
func (ra *RA) String() string { return fmt.Sprintf("%s", ra) }

// Time represents a formattable duration or relative time.
type Time struct {
	unit.Time
	Sym *Symbols
	Err error // set each time the value is formatted.
}

// FmtTime constructs an formattable Time containing the value t.
func FmtTime(t unit.Time) *Time { return &Time{Time: t} }

// Format implements fmt.Formatter, formatting to hours, minutes, and seconds.
func (t *Time) Format(f fmt.State, c rune) {
	s := &state{
		State:  f,
		verb:   c,
		hrDeg:  t.Hour(),
		caller: fsTime,
		sym:    t.Sym,
	}
	t.Err = s.writeFormatted()
}

// String implements fmt.Stringer
func (t *Time) String() string { return fmt.Sprintf("%s", t) }

// FmtAngle constructs an formattable Angle containing the value a.
func (sym *Symbols) FmtAngle(a unit.Angle) *Angle { return &Angle{a, sym, nil} }

// FmtHourAngle constructs an formattable HourAngle containing the value h.
func (sym *Symbols) FmtHourAngle(h unit.HourAngle) *HourAngle {
	return &HourAngle{h, sym, nil}
}

// FmtRA constructs an formattable RA containing the value ra.
func (sym *Symbols) FmtRA(ra unit.RA) *RA { return &RA{ra, sym, nil} }

// FmtTime constructs an formattable Time containing the value t.
func (sym *Symbols) FmtTime(t unit.Time) *Time {
	return &Time{t, sym, nil}
}

// CombineUnit inserts a unit indicator into a formatted decimal number,
// combining it if possible with the decimal separator.
//
// If sym.DecSep is non-empty and occurrs in d, the occurrence is replaced with
// argument 'unit' and the symbol sym.DecCombine.  Otherwise unit is appended
// to the end of d.
//
// See also InsertUnit, StripUnit, and the corresponding top-level functions
// that use package default symbols.
func (sym *Symbols) CombineUnit(d, unit string) string {
	if sym.DecSep == "" || sym.DecCombine == 0 {
		return d + unit // DecSep empty, append unit
	}
	i := strings.Index(d, sym.DecSep)
	if i < 0 {
		return d + unit // no DecSep found, append unit
	}
	// insert unit, replace DecSep occurrence with DecCombine
	return d[:i] + unit + string(sym.DecCombine) + d[i+len(sym.DecSep):]
}

// InsertUnit inserts a unit indicator into a formatted decimal number.
//
// The indicator is inserted just before the decimal separator if one is
// present, or at the end of the number otherwise.
//
// If sym.DecSep is non-empty and occurrs in d, unit is added just before the
// occurrence.  Otherwise unit is appended to the end of d.
//
// See also CombineUnit, StripUnit, and the corresponding top-level functions
// that use package default symbols.
func (sym *Symbols) InsertUnit(d, unit string) string {
	if sym.DecSep == "" {
		return d + unit // DecSep empty, append unit
	}
	i := strings.Index(d, sym.DecSep)
	if i < 0 {
		return d + unit // no DecSep found, append unit
	}
	// insert unit before DecSep
	return d[:i] + unit + d[i:]
}

// StripUnit reverses the action of InsertUnit or CombineUnit,
// removing the specified unit indicator and restoring a following
// DecCombine to DecSep.
//
// More specifically, if the specified unit is found at the end of string d
// or if sym.DecSep is non-empty and ocurrs immediately following, then the unit
// is Removed.  If the specified unit is found with sym.DecCombine immediately
// following, then both the unit and the DecCombine rune are replaced with
// sym.DecSep.
//
// StripUnit returns ok = true if the unit was found and removed.  Otherwise it
// returns d unchanged and ok = false.
func (sym *Symbols) StripUnit(d, unit string) (stripped string, ok bool) {
	xu := strings.Index(d, unit)
	if xu < 0 {
		return d, false
	}
	xd := xu + len(unit)
	if xd == len(d) {
		return d[:xu], true // string ends with unit.  just remove the unit.
	}
	if sym.DecSep != "" && strings.HasPrefix(d[xd:], sym.DecSep) {
		return d[:xu] + d[xd:], true // remove unit, retain DecSep
	}
	if r, sz := utf8.DecodeRuneInString(d[xd:]); r == sym.DecCombine {
		// replace unit and DecCombine with DecSep
		return d[:xu] + sym.DecSep + d[xd+sz:], true
	}
	return d, false // otherwise don't mess with it
}

const (
	secAppend    = 's'
	secCombine   = 'c'
	secInsert    = 'd'
	minAppend    = 'm'
	minCombine   = 'n'
	minInsert    = 'o'
	hrDegAppend  = 'h'
	hrDegCombine = 'i'
	hrDegInsert  = 'j'
)

const (
	fsAngle = iota
	fsHourAngle
	fsRA
	fsTime
)

type state struct {
	fmt.State         // 'f' in fmt.Formatter doc.  kind of handy to embed this.
	verb      rune    // 'c' in fmt.Formatter doc
	hrDeg     float64 // input, value to format
	prec      int     // f.Precision with a default of 0
	caller    int     // use fs constants
	sym       *Symbols
	units     UnitSymbols
}

func (s *state) writeFormatted() error {
	if s.sym == nil {
		s.sym = Default
	}
	switch {
	case s.caller == fsAngle:
		s.units = s.sym.DMSUnits
	default:
		s.units = s.sym.HMSUnits
	}

	// valiate verb, pick formatting method in the process
	var f func() (string, error)
	switch s.verb {
	case 'v':
		fallthrough
	case secAppend, secCombine, secInsert:
		f = s.decimalSec // it's a method value! see the spec.
	case minAppend, minCombine, minInsert:
		f = s.decimalMin
	case hrDegAppend, hrDegCombine, hrDegInsert:
		f = s.decimalHrDeg
	default:
		fmt.Fprintf(s, "%%!%c(BADVERB)", s.verb)
		return nil // not a value error
	}

	// validate precision, storing it in the receiver.
	// 0 is our default if it's not specified.
	// (the docs don't define what prec is returned for the !ok case)
	var ok bool
	switch s.prec, ok = s.Precision(); {
	case !ok:
		s.prec = 0
	case s.prec > 15:
		// limit of 15 set by max power of 10 that is exactly representable
		// as a float64.  later code depends on prec being in this range.
		fmt.Fprintf(s, "%%!(BADPREC %d)", s.prec)
		return nil // not a value error
	}

	// format validated, now preliminary checks on value:
	var (
		r   string
		err error
	)
	switch {
	case math.IsNaN(s.hrDeg):
		err = ErrNaN
		goto valErr
	case !math.IsInf(s.hrDeg, 0): // normal path
	case math.IsInf(s.hrDeg, 1):
		err = ErrPosInf
		goto valErr
	default:
		err = ErrNegInf
		goto valErr
	}
	// and then call the formatting method picked above
	if r, err = f(); err == nil {
		s.Write([]byte(r))
		return nil // normal return
	}

	// If there was a value error, we output all '*'s
	// but we need a length.  The strategy here is to replace the invalid
	// value with something valid and call format again to get a mock
	// result, then use len(mock) for the number of '*'s to output.
valErr:
	s.hrDeg = 0
	width := 10 // default, defensive in case f somehow fails on 0.
	if mock, err2 := f(); err2 == nil {
		width = utf8.RuneCountInString(mock)
		if strings.IndexRune(mock, s.sym.DecCombine) >= 0 {
			width--
		}
	}
	s.Write(bytes.Repeat([]byte{'*'}, width))
	return err
}

var (
	tenf = [16]float64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5,
		1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15}
	teni = [16]int64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5,
		1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15}
)

// sig verifies and returns significant digits of a number at a precision.
//
// x must be >= 0.  prec must be 0..15.
//
// the digits are returned as xs = int64(x * 10**prec + .5), as long as
// the result xs is small enough that all digits are significant given
// float64 representation.
// if xs does not represent a fully significant result -1 is returned.
func sig(x float64, prec int) int64 {
	xs := x*tenf[prec] + .5
	if !(xs <= 1<<52) { // 52 mantissa bits in float64
		return -1
	}
	return int64(xs)
}

func (s *state) decimalHrDeg() (string, error) {
	i := sig(math.Abs(s.hrDeg), s.prec)
	if i < 0 {
		return "", ErrLossOfPrecision
	}
	if s.hrDeg < 0 {
		i = -i
	}
	var r, f string
	if wid, widSpec := s.Width(); !widSpec {
		if s.Flag('+') {
			f = "%+0*d"
		} else if s.Flag(' ') { // sign space if requested
			f = "% 0*d"
		} else {
			f = "%0*d"
		}
		// +1 forces at least one place left of decimal point
		r = fmt.Sprintf(f, s.prec+1, i)
	} else {
		// fixed width a little more involved
		if s.Flag('+') {
			f = "%+"
		} else {
			f = "% " // sign space forced with fixed width
		}
		if s.Flag('0') {
			f += "0*d"
		} else {
			f += "*d"
		}
		wf := s.prec + wid + 1 // +1 here is required space for sign
		r = fmt.Sprintf(f, wf, i)
		if len(r) > wf {
			if s.caller == fsAngle {
				return "", ErrDegreeOverflow
			}
			return "", ErrHourOverflow
		}
	}
	if s.prec > 0 {
		split := len(r) - s.prec
		r = r[:split] + s.sym.DecSep + r[split:]
	}
	switch s.verb {
	case hrDegAppend:
		r += string(s.units.HrDeg)
	case hrDegCombine:
		r = s.sym.CombineUnit(r, s.units.HrDeg)
	case hrDegInsert:
		r = s.sym.InsertUnit(r, s.units.HrDeg)
	}
	return r, nil
}

func (s *state) decimalMin() (string, error) {
	i := sig(math.Abs(s.hrDeg)*60, s.prec) // hrDeg*60 gets minutes
	if i < 0 {
		return "", ErrLossOfPrecision
	}
	p60 := 60 * teni[s.prec]
	min := i / p60
	sec := i % p60

	r, minEl, err := s.firstSeg(min)
	if err != nil {
		return "", err
	}
	return r + s.lastSeg(sec, s.units.Min, minEl), nil
}

func (s *state) firstSeg(x int64) (r string, elided bool, err error) {
	wid, widSpec := s.Width()
	switch {
	case widSpec:
		f := "%*d"
		if s.Flag('0') {
			f = "%0*d"
		}
		r = fmt.Sprintf(f, wid, x)
		if len(r) > wid {
			if s.caller == fsAngle {
				return "", false, ErrDegreeOverflow
			}
			return "", false, ErrHourOverflow
		}
		r += s.units.HrDeg
	case x > 0 || s.Flag('#'):
		r = fmt.Sprintf("%d%s", x, s.units.HrDeg)
	default:
		elided = true
	}
	switch {
	case s.hrDeg < 0:
		r = "-" + r
	case s.Flag('+'):
		r = "+" + r
	case s.Flag(' ') || widSpec:
		r = " " + r
	}
	return r, elided, nil
}

func (s *state) lastSeg(sec int64, unit string, first bool) string {
	wid := s.prec + 1
	_, widSpec := s.Width()
	if s.Flag('0') && (widSpec || !first) {
		wid++
	}
	r := fmt.Sprintf("%0*d", wid, sec)
	if widSpec && len(r) < s.prec+2 {
		r = " " + r
	}
	if s.prec > 0 {
		split := len(r) - s.prec
		r = r[:split] + s.sym.DecSep + r[split:]
	}
	switch s.verb {
	case secCombine, minCombine:
		return s.sym.CombineUnit(r, unit)
	case secInsert, minInsert:
		return s.sym.InsertUnit(r, unit)
	}
	return r + unit
}

func (s *state) decimalSec() (string, error) {
	i := sig(math.Abs(s.hrDeg)*3600, s.prec) // hrDeg*3600 gets seconds
	if i < 0 {
		return "", ErrLossOfPrecision
	}
	p60 := 60 * teni[s.prec]
	sec := i % p60
	i /= p60
	min := i % 60
	hrDeg := i / 60
	r, firstEl, err := s.firstSeg(hrDeg)
	if err != nil {
		return "", err
	}
	f := "%s%d%s"
	minEl := false
	if s.Flag('0') && !firstEl {
		f = "%s%02d%s"
	} else {
		switch _, widSpec := s.Width(); {
		case widSpec:
			f = "%s%2d%s"
		case firstEl && min == 0:
			minEl = true
			goto last
		}
	}
	r = fmt.Sprintf(f, r, min, s.units.Min)
last:
	return r + s.lastSeg(sec, s.units.Sec, minEl), nil
}
