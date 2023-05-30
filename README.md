Converts CSV to [pprof](https://github.com/google/pprof) profile format.

Use csv2pprof when you have some data in a database or a spreadsheet you'd like
to turn into a pprof profile.

## Installation

```
go install github.com/google/csv2pprof@latest
```

## Usage

```
csv2pprof < input.csv > output.pprof.gz
```

Input CSVs must have:
- a header row
- a semicolon-delimited `stack` column (similar to Brendan Gregg's Folded Stacks
  format)
- one or more integer measurement columns (e.g. samples, or time). Measurements
  can have units given after a forward-slash `/`.

Example CSV input:

```
stack,cpu-time/milliseconds
main;foo,1000
main;foo;bar,2000
main;baz,4000
```

Or you can have many measurement columns.

```
stack,cpu-time/milliseconds,samples,instructions
main;foo,1000,10,100
main;foo;bar,2000,20,200
main;baz,4000,40,400
```


## See Also / Prior Art

- [Brendan Gregg's Folded Stacks Format](https://github.com/brendangregg/FlameGraph) uses semicolon-separated stacks, and a space-separated measurement.
- [felixge's pprofutils](https://github.com/felixge/pprofutils) converts folded stack format to pprof.

The folded stack format is just a little harder to generate from databases and spreadsheets.
