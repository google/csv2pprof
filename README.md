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
cpu-time/milliseconds,stack
1000,main;foo
2000,main;foo;bar
4000,main;baz
```

Or you can have many measurement columns:

```
cpu-time/milliseconds,samples,instructions,stack
1000,10,100,main;foo
2000,20,200,main;foo;bar
4000,40,400,main;baz
```

If you want to use a different separator for the stack than semicolon, use
`--stacksep`:

```
$ csv2pprof --stacksep="\n" < input.csv > pprof.pb.gz
```


## See Also / Prior Art

- [Brendan Gregg's Folded Stacks Format](https://github.com/brendangregg/FlameGraph) uses semicolon-separated stacks, and a space-separated measurement.
- [felixge's pprofutils](https://github.com/felixge/pprofutils) converts folded stack format to pprof.

The folded stack format is just a little harder to generate from databases and spreadsheets.
