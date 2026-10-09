# Validate Properly Nested Brackets

**Source:** HackerRank · **Pattern:** Stack — matching pairs

## Problem

Given a string, check that `()`, `{}` and `[]` are all matched and properly nested. Return 1
if valid, 0 otherwise. Everything that isn't a bracket is ignored.

```
if (a[0] > b[1]) { doSomething(); }   -> 1
int x = 42; // no brackets here       -> 1
() {} []                              -> 1
```

### Constraints

- `0 <= len <= 1000`
- printable ASCII, codes **32–126** (note: whitespace is allowed this time)
- may be empty

## Recognition signals

- **"properly nested" / "matched pairs"** → stack. Nesting is last-in-first-out by
  definition: the bracket you opened most recently must be the one you close first.
- More generally: whenever you need "the most recent unmatched thing", that's a stack.
- **Three bracket types is what forces the stack.** With only one type you'd just keep a
  counter (depth), increment and decrement, fail if it goes negative. Add a second type and
  a counter can't tell `([)]` from `([])` — you need to remember *what* was opened, in order.

## Idea

Walk the string once.

- **Opener** (`(`, `{`, `[`) → push it.
- **Closer** (`)`, `}`, `]`) → the top of the stack must be its matching opener. If the stack
  is empty, or the top doesn't match, fail. Otherwise pop.
- **Anything else** → skip.

At the end the stack must be **empty**. Leftovers mean unclosed openers.

```
stack = []
for each ch:
    if ch is an opener:  push ch
    if ch is a closer:
        if stack empty:              return false
        if top != matching opener:   return false
        pop
return stack is empty
```

Complexity: **O(n) time, O(n) space** (worst case `"((((("` pushes everything).

Go: a slice is the stack. `push = append(s, ch)`, `top = s[len(s)-1]`,
`pop = s[:len(s)-1]`. A `map[byte]byte` from closer to opener keeps the matching tidy.

## Traps

- **Check the stack is empty at the end.** `"("` pushes and never fails inside the loop —
  returning true there is the classic bug.
- **Closer on an empty stack** → false. `")"` must fail, and `s[len(s)-1)]` panics if you
  don't check first.
- **The types must match**, not just the counts. `"([)]"` is invalid; a depth counter says
  it's fine.
- Empty string → true. The loop never runs, the stack is empty.
- Ignore non-brackets — letters, digits, spaces, `;`, `>` all just pass through.

## Notes

Solved 2026-08-19 in 27 min, **first attempt, no bugs** — 22 table cases plus 20,000 fuzz
trials against an independent oracle (strip non-brackets, repeatedly delete adjacent pairs,
valid iff it collapses to nothing).

What made it clean:

1. **Checked `len(stack) == 0` before peeking.** Without it, `")"` panics instead of
   returning false.
2. **Ended with `return len(stack) == 0`.** The most-forgotten line in this problem — and it
   came straight from the invariant, not from memory. *"The stack holds exactly the unclosed
   openers"* + *"the loop is over"* = *"empty means valid"*.
3. **Wrote the state down before coding.** Two for two now: every problem where I named the
   state first had no state bug.

**State vs derived, again.** My first answer listed "last opening element" and "current
closing element" as state. Neither is: the first is `stack[len(stack)-1]`, the second is the
loop variable. The state is *just the stack*. Storing the top separately means two things to
keep in sync.

**Invariant template that works:** `<state>` **holds exactly** `<what it should hold>`. If
your sentence contains "if ... then invalid", you've written the rule, not the invariant.

Review nits: hoist the map to a package-level `var` (it's rebuilt per call); make the two
branches `else if` since a char can't be both; rename `tail` → `top` and `v` → `wantOpener`.
