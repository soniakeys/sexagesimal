# Sexagesimal

This package provides sexagesimal formatting for the four types defined in the
external package github.com/soniakeys/unit.  Those types are Angle, HourAngle,
RA, and Time.

## Motivation

Code in this package was written to support the external package
github.com/soniakeys/meeus, a collection of astronomy routines.  Actually
though, sexagesimal is used only for nicely readable examples in the meeus
package.  It is thus only a test dependency.  Neither this package nor
meeus depend on each other.  Both depend on unit.

## Install

### Go get

As usual,

    go get github.com/soniakeys/sexagesimal

### Vgo

Experimentally, you can try [vgo](https://research.swtch.com/vgo).

To run package tests, clone the repository -- anywhere! it doesn't have to
be under GOPATH -- and from the cloned directory run

    vgo test

Vgo will fetch the unit package dependency as needed and run the sexagesimal
package tests.

### Or don't install it

If you only need `sexagesimal` as dependency of some other package that you
are installing, the normal installation of that package will likely install
`sexagesimal` for you.  Try that first.

## License

MIT
