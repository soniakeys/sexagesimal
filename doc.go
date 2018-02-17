// License: MIT

// Package sexa provides sexagesimal formatting for four types defined in the
// external package github.com/soniakeys/unit, Angle, HourAngle, RA, and Time.
//
// Four types in this package wrap the four types from package unit, adding
// formatting symbol options and an error value.  Formatting is done with Go
// custom formatters.
//
// Formatting symbols
//
// Unit indicators and decimal separators are defined in the Symbols type.
// There is a package default but symbols can also be defined as needed.
// The symbols used for degrees, minutes, and seconds for the FmtAngle type
// are taken from Symbols.DMSUnits.  The symbols for hours, minutes, and
// seconds for the FmtHourAngle, FmtRA, and FmtTime types are taken from
// Symbols.HMSUnits.
//
// A sexagesimal format can have up to three numeric segments.
// The decimal separator, if it appears, is always in the last segment.
// Symbols used for decimal separators are taken from Symbols.DecSep and
// Symbols.DecCombine.
//
// Decimal unit indication
//
// Three conventions are supported for unit indication on the decimal segment.
// By default (with %v, for example) the unit follows the segment.
//
// 1°23′45.6″
//
// There is a typgraphic convention however, for moving the final unit
// indicator to the decimal separator and placing the decimal separator
// directly under the unit symbol.  This sometimes can be approximated in
// Unicode with codes of the category "Mn", for example "combining dot below"
// u+0323.  Example (that may or may not look right*)
//
// 1°23′45″̣6
//
// For cases where software does not render this satisfactorily, an
// alternative convention is to simply insert the unit symbol ahead of the
// decimal separator as in
//
// 1°23′45″.6
//
//   * Footnote about combining dot.  The combining dot only looks right
//     to the extent that software (such as fonts and browsers) can render it.
//     See http://www.unicode.org/faq/char_combmark.html#12b for a description
//     of the issues.  It seems that monospace fonts are more problematic.
//     The examples above are aligned flush left to avoid godoc coding
//     them monospace in the HTML.  For example 1°23′45″̣6 is less likely to
//     look right.  Other contexts likely to use monospace fonts and so likely
//     to have trouble with the combining dot are operating system shells and
//     code text editors.
//
// Format specifiers
//
// The syntax of a format specifier is
//
//    %[flags][width][.precision]verb
//
// The syntax is set by the Go fmt package, but this package customizes
// the meaning of all format specifier components.
//
// Verbs specify one of the above decimal unit conventions and also the unit
// of the decimal (right most) segment.  The decimal unit determines the
// the potential number of segments.  Full sexagesimal format has three
// segments with the decimal separator in seconds.  Decimal minutes format has
// an hour or degrees segment, a minutes segment with the decimal separator,
// and no seconds segment.  Decimal hour or degree format has only a single
// decimal segment.
//
// This table gives the verbs for the combinations of decimal unit indication
// and decimal segment:
//
//    decimal-unit indication:             following  combined  inserted
//
//    three segments, decimal in seconds:      %s        %c        %d
//    two segments, decimal in minutes:        %m        %n        %o
//    one segment, decimal in hr/degs:         %h        %i        %j
//
// Also %v is equivalent to %s.
//
// The following flags are supported:
//  +   always print leading sign
//  ' ' (space) leave space for elided + sign
//  #   display all segments, even if 0
//  0   pad displayed segments with leading zeros
//
// A + flag takes precedence over a ' ' (space) flag.
//
// The # flag forces output to have all segments, even if 0.  Without it,
// leading zero segments are elided.  (Consider formatting coordinates with #;
// distances and durations without.)
//
// The 0 flag pads with a leading zero on non-first (sexagesimal) segments.
// If a width is specfied, the 0 flag pads with leading zeros on the first
// (hr/deg) segment as well.
//
// For the RA type, sign formatting flags '+' and ' ' are ignored.
//
// Specifying width forces a fixed width format.  Flag '#' is implied, ' ' is
// implied unless '+' is given, and segments are space padded unless '0' is
// given.  The width number specifies the number of digits in the integer part
// of the most significant segment, hours or degrees — not the total width.
// For example you would typically use the number 2 for RA, 3 for longitude.
// Also with fixed width consider avoiding the combining dot verbs unless you
// also control output rendering. (See note above on rendering of the combining
// dot.)  With fixed width sexagesimal formats, the sign indicator is always
// the left-most column; with fixed width space padded decimal hour or degree
// formats, the sign indicator is formatted immediately in front of the number
// within the space padded field.
//
// Precision specifies the number of places past the decimal separator
// of the decimal segment.  The default is 0.  There is no variable precision
// format.
//
// Without a specified width the format is not fixed width but of course you
// can always format the result into a fixed width string with an additional
// printf step.  Since all formats of this package have fixed precision, if
// you right justify the string (the default) then at least the decimal points
// will align.  Additionally if you use the '0' flag, then all segments after
// the first will have fixed width so that all unit indicators will align as
// well.
//
// Errors
//
// A value that cannot be expressed the in the requested format represents
// an overflow condition.  In this case, the custom formatters emit all
// asterisks "*************" and leave a more descriptive error in the
// Err field of the value.
//
// If you specify width, digits of the integer part of the first segment must
// fit in the specified width.  Larger values cause overflow.
//
// Overflow also happens if more precision is requested than is represented
// in the underlying float64.  In the case of an angle formatted with the
// decimal separator in seconds, precision of 15 is possible only for angles
// less than a few arc seconds.  As angle values increase, fewer digits of
// precision are possible.  At one degree, you can get 12 digits of precision
// in the seconds segment of a full sexagesimal number, at 360 degrees,
// you can get 9.  For all formats, an angle too large for the specified
// precision causes overflow.
//
// +Inf, -Inf, and NaN always cause overflow.
//
// Only errors related to the value being formatted are handled as overflow
// and leave a non-nil Err field.  Errors of format specification are handled
// with the standard Printf convention of emitting the error in the formatted
// result and leave the Err field nil.
package sexa
