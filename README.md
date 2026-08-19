# algorithms

Practice ground for data structures and algorithms interview prep, in Go.

Each problem gets its own folder, named after the **pattern** rather than the problem story —
recognising the pattern is the skill being drilled, so the folder name is part of the
exercise. A folder holds the implementation, a `README.md` with the statement, recognition
signals, approach, traps and review notes, and a test file.

```
arrays/
  binary_search/
    binary_search.go
    binary_search_test.go
    README.md
strings/
  two_pointers/
    ...
```

Run everything:

```sh
go test ./...
go vet ./...
```

- [AGENTS.md](./AGENTS.md) — how problems get worked here, plus the running pattern catalogue.
- [COLD_REDO.md](./COLD_REDO.md) — spaced-repetition queue. Statements only, no hints.
