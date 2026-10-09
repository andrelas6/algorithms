# Non-Trivial String Rotation

**Source:** HackerRank · **Pattern:** Concatenation trick (`s1+s1` contains every rotation)

## Problem

Given two strings `s1` and `s2`, return 1 if `s2` is a rotation of `s1` **but not identical**
to `s1`, otherwise 0.

```
s1 = abcde
s2 = cdeab      -> 1   (rotate left by 2)
```

### Constraints

- `1 <= |s1|, |s2| <= 1000`, and `|s1| == |s2|`
- lowercase English letters only

### Samples

| s1 | s2 | Output |
|---|---|---|
| `a` | `a` | 0 (identical) |
| `a` | `b` | 0 (not a rotation) |

## Recognition signals

- "is X a rotation of Y" → **double one of them**. This is a named trick, not something to
  re-derive under pressure.
- Anything **circular / wrap-around** hints at the same move: concatenate the sequence to
  itself, or index with `% n`, so the wrap becomes a plain linear scan.
- "but not identical" is a separate predicate bolted on top — handle it independently of the
  rotation test.

## Idea

Every rotation of `s1` appears in `s1+s1` as a contiguous window of length `n`: rotating left
by `k` gives `s1[k:] + s1[:k]`, which is exactly `(s1+s1)[k : k+n]`. Sliding `k` from `0` to
`n-1` enumerates all `n` rotations, and nothing else of length `n` starts before index `n`.

```
s1 = abcde
s1+s1 = abcdeabcde
         k=0 abcde
          k=1 bcdea
           k=2 cdeab   <- s2
```

So, given equal lengths:

```
if len(s1) != len(s2):      return false
if s1 == s2:                return false      // trivial rotation excluded
return s1+s1 contains s2
```

Complexity: O(n) extra space for the doubled string. Time is the substring search — O(n) with
KMP/Z-algorithm, and in practice `strings.Contains` (Rabin-Karp) is O(n) average, O(n²) worst.
At `n <= 1000` any of it passes.

**O(1)-space alternative:** for each `k` in `1..n-1`, compare `s2[i] == s1[(i+k)%n]` for all
`i`. O(n²) time, no allocation. Worth naming as the trade-off.

## Traps

- **The length check is load-bearing**, not defensive noise. `s1="abc"`, `s2="ab"` is a
  substring of `"abcabc"` but is not a rotation. The constraints promise equal lengths; the
  trick is only correct *because* of that.
- **`s1 == s2` must be rejected** — that's the "non-trivial" clause, and sample 0 exists
  purely to catch it. Rotation by `0` and by `n` are the identity.
- All-identical strings (`"aaaa"`) — every rotation equals the original, so the answer is
  always 0. The identity check handles it; make sure it runs *before* the contains test.
- Don't build `s2+s2` and check `s1` in it by mistake — same answer here since lengths match,
  but be deliberate about which one you double.

## Notes

Solved 2026-08-19 in 21 min, all tests passed. Carry forward:

1. **The trick ends in a library call: `strings.Contains(s1+s1, s2)`.** I hand-rolled the
   substring search as a sliding window — correct, but O(n²) instead of Rabin-Karp's O(n)
   average, and 15 lines instead of 3. Habit to build: once the idea is in place, ask *"is
   the remaining step stdlib?"* before writing the loop.
2. **The `s1 == s2` guard is not redundant with starting the loop at `k=1`.** For a periodic
   string a *non-zero* rotation can still equal `s1` — `"abab"` at `k=2` is `"abab"`. The
   identity check must run first.
3. That loop is a **fixed-size sliding window**, not two pointers. One index, constant width.
   Naming it `pointer1`/`pointer2` hid which pattern I was actually in.
4. `len(s1) == 1` early return was dead code — the loop never runs for `n=1` anyway. And
   `pointer2` was derived state (`pointer1 + len(s1)`); inline it so it can't drift.
