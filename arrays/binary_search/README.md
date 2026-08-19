# Target Index Search

**Source:** HackerRank · **Pattern:** Binary search (textbook form)

## Problem

Given a sorted array of **distinct** integers and a target, return the index of the target,
or `-1` if absent.

```
nums = [1, 2, 3, 4, 5], target = 3   -> 2
nums = [2,4,6,8,10,12,14,16], target = 16 -> 7
```

### Constraints

- `0 <= n <= 10^6`
- `-10^9 <= nums[i], target <= 10^9`
- strictly increasing: `nums[i] < nums[i+1]`

### Samples

| n | nums | target | Output |
|---|---|---|---|
| 0 | `[]` | 5 | -1 |
| 1 | `[10]` | 10 | 0 |

## Recognition signals

- **"sorted" + "find"** → binary search. Close to a reflex; sorted input that you are not
  asked to sort is there for a reason.
- `n <= 10^6` with a lookup → O(log n) expected. A linear scan also passes at this size, so
  the constraint is not what forces it — the word *sorted* is.
- Distinct values remove the "find the first/last occurrence" complication. With duplicates
  this becomes the lower-bound variant instead.

## Idea

Invariant: **if the target is present, its index is inside `[lo, hi]`.** Each step looks at
the middle and discards the half that cannot contain it.

```
lo = 0, hi = n - 1
while lo <= hi:
    mid = lo + (hi - lo) / 2
    if nums[mid] == target: return mid
    if nums[mid] < target:  lo = mid + 1     // target is right of mid
    else:                   hi = mid - 1     // target is left of mid
return -1
```

The range halves every iteration, so it empties after `log2(n)` steps — 20 for `n = 10^6`.

Complexity: **O(log n) time, O(1) space** (iterative; a recursive version costs O(log n)
stack).

## Traps

- **`lo <= hi`, not `lo < hi`**, when `hi` is an *inclusive* bound. With `<` a
  single-candidate range `lo == hi` is never examined, so `[10]`/`10` returns -1.
- **`mid + 1` / `mid - 1`**, not `mid`. Leaving `mid` in the range means a range that stops
  shrinking → infinite loop.
- Empty array: `hi = -1`, the loop never runs, `-1` falls out. No special case needed.
- **`mid = lo + (hi-lo)/2`, not `(lo+hi)/2`** — the latter overflows for large bounds. Not
  reachable with `int` indices in Go, but it is the canonical bug and interviewers ask.
- Pick **one** convention and stay in it: inclusive `[lo, hi]` with `lo <= hi` and
  `hi = mid-1`, **or** half-open `[lo, hi)` with `hi = n`, `lo < hi` and `hi = mid`. Mixing
  them is where off-by-ones come from.

## Notes

Solved 2026-08-19, after one wrong attempt that timed out. Carry forward:

1. **The state is the range, not the midpoint.** First attempt tracked only `midIndex` and
   recomputed it from `0` and `len` — constants that never change — so nothing was ever
   discarded and the index oscillated (`6 → 3 → 6 → 3`) forever. `mid` is *derived*: compute
   it inside the loop from `lo`/`hi`, use it, throw it away. A range needs two numbers; one
   number cannot say which half you already ruled out.
2. **Termination is structural once `lo`/`hi` are the state.** Every iteration does
   `lo = mid+1` or `hi = mid-1`, so `hi - lo` strictly decreases. An infinite loop becomes
   impossible — you don't need a clever loop condition to prevent one.
3. **Get the range right and the small cases stop being special.** Empty (`hi = -1`, loop
   never runs) and single-element (`lo == hi == mid`, examined because the condition is
   `lo <= hi`) both fall out. The `if len(nums) == 1` guard from the first attempt was
   papering over a loop condition that excluded index `0` and `n-1` entirely.
4. **`(low+high)/2` is not wrong here** — same value, and `n <= 10^6` can't overflow. But
   `low + (high-low)/2` is unconditionally safe (`hi-lo <= hi`, which is a valid index by
   definition) and is the standard follow-up question. Make it the default.

Final bug was `return 1` instead of `return middle` — every found case returned the same
value. Caught by testing several indices of the *same* array rather than one case per array.
