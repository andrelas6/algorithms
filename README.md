# algorithms

Daily practice for data structures and algorithms, in Go.

Each problem gets its own folder, grouped by where it came from and by topic:
`<source>/<topic>/<problem>/`. A folder holds the implementation and a test file, plus a
`README.md` with the statement, recognition signals, approach and traps when the problem
went through the full teaching loop.

```
neetcodeio/
  stack/
    car_fleet/
      car_fleet.go
      car_fleet_test.go
hackerrank/
  arrays/
    two_sum/
      two_sum.go
      two_sum_test.go
      README.md
cracking_the_coding_interview/
  arrays_lists/
    urlify/
      ...
```

Run everything:

```sh
go test ./...
go vet ./...
```

`AGENTS.md` describes how problems get worked here and holds the running pattern catalogue.
