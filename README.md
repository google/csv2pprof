Converts CSV to [pprof](https://github.com/google/pprof) profile format.

Use when you have some data in a database or a spreadsheet you'd like to turn
into a pprof profile.

## Installation

```
go install github.com/mhansen/csv2pprof@latest
```

## Usage

```
csv2pprof < input.csv > output.pprof
```

Generate a stack column of semicolon-separated frame names, similar to Brendan
Gregg's "Folded Stacks" format, and an integer measurement column.

Example CSV input:

```
stack,cpu-time/milliseconds
main;foo,1000
main;foo;bar,2000
main;baz,4000
```

Or you can have many measurement columns. Give units after the forward-slash
`/`:

```
stack,cpu-time/milliseconds,samples,instructions
main;foo,1000,10,100
main;foo;bar,2000,20,200
main;baz,4000,40,400
```

CSVs must have:
- a header row
- a semicolon-delimited `stack` column
- one or more measurement columns


## See Also / Prior Art

- [Brendan Gregg's Folded Stacks Format](https://github.com/brendangregg/FlameGraph) uses semicolon-separated stacks, and a space-separated measurement.
- [felixge's pprofutils](https://github.com/felixge/pprofutils) converts folded stack format to pprof.

The folded stack format is just a little harder to generate from databases and spreadsheets.
