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

### Git clone

Alternatively, you can clone the repository into an appropriate place under
your GOPATH.  To clone into the same place as `go get` for example, assuming
the default GOPATH of ~/go, you would cd to `~/go/src/github.com/soniakeys`
before running the clone command.

    cd <somewhere under GOPATH>
    git clone https://github.com/soniakeys/sexagesimal

### Dep

Sexagesimal imports `github.com/soniakeys/unit`.  You can also use dep
(https://golang.github.io/dep/) to "vendor" this package.

To use dep, first read about dep on the website linked above and install it.
Then install `unit` with either `go get` or `git clone`.  Finally, from the
installed `sexagesimal` directory, type

    dep ensure

This will "vendor" the unit package, installing it under the `vendor`
subdirectory and also installing a specific version of unit known to
work with the version of sexagesimal that you just installed.

### Or don't install it

If you only need `sexagesimal` as dependency of some other package that you
are installing, the normal installation of that package will likely install
`sexagesimal` for you.  Try that first.

## License

MIT
