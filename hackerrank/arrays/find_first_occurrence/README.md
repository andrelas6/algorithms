# Find First Occurrence

**Source:** HackerRank · **Pattern:** Binary search — lower bound (leftmost) variant

## Problem

Sorted array that **may contain duplicates**. Return the index of the **first** occurrence of
`target`, or -1 if it's not there.

```
nums = [1, 2, 3, 4, 5], target = 3  ->  2
```

### Constraints

- `0 <= n <= 1000`
- `-10^9 <= nums[i], target <= 10^9`
- non-decreasing: `nums[i] <= nums[i+1]` (note: `<=`, not `<`)

### Samples

| n | nums | target | Output |
|---|---|---|---|
| 0 | `[]` | 5 | -1 |
| 1 | `[3]` | 3 | 0 |

## Recognition signals

- Sorted + find → binary search. Same as before.
- **"first occurrence" / "last occurrence" / duplicates allowed** → the *lower bound* variant,
  not plain binary search. This is the tell. Plain binary search returns *some* match; which
  one you get depends on where the midpoints land, so it's useless here.
- `n <= 1000` is tiny — a linear scan passes. The word *sorted* is what forces binary search,
  not the size.

## Idea

One change from plain binary search: **finding the target doesn't end the search.**

When `nums[mid] == target`, you've found *an* occurrence, but maybe not the first. There
could be more to the left. So write it down and keep looking left.

```
lo = 0, hi = n - 1
best = -1

while lo <= hi:
    mid = lo + (hi - lo)/2
    if nums[mid] == target:
        best = mid          # remember it
        hi = mid - 1        # but keep hunting to the LEFT
    else if nums[mid] < target:
        lo = mid + 1
    else:
        hi = mid - 1

return best
```

`best` starts at -1, so "not found" falls out with no extra check.

Complexity: **O(log n) time, O(1) space.** You don't stop early on a match, so it's always
the full ~10 steps for `n = 1000` — same big-O, slightly worse constant than plain search.

## Traps

- **Don't `return mid` on a match.** That's the plain-binary-search reflex and it gives you an
  arbitrary occurrence, not the first.
- After recording, you still must do `hi = mid - 1`. Setting `hi = mid` leaves `mid` in range
  → infinite loop.
- Empty array: `hi = -1`, loop never runs, `best` is still -1. No special case.
- All elements smaller than target, or all bigger → `best` stays -1. Also free.
- The array is **non-decreasing**, so `nums[i] == nums[i+1]` is legal. That's the whole point
  of the problem — don't assume distinct.

## Notes

Solved 2026-08-19 in 6 min, first attempt, all tests passed (including a fuzz test against a
brute-force linear scan over 2000 random duplicate-heavy arrays).

1. **"Found it" and "too big" collapse into the same branch.** Both mean *go left*. So you
   don't need a three-way if — record the match, then let `==` fall into the `else` that does
   `high = mid - 1`. Tidy, but subtle enough to deserve a comment.
2. **`found` comes out leftmost automatically.** After a match you set `high = mid - 1`, so
   every later `mid` is strictly smaller. Matches can only be recorded in decreasing order,
   and the last one written wins.
3. **Naming the state and invariant before coding worked.** Wrote `state: low, high, best`
   and `invariant: if target is in array, it is within low-high bounds` as comments first —
   and there was no state bug, unlike `arrays/binary_search/`.
4. Wrote `(low+high)/2` on the first pass, fixed to `low + (high-low)/2` after review. Same
   value either way; the second can't overflow because `high-low <= high`, and `high` is a
   valid index by definition. Make it the default so it's automatic under pressure.

**On the invariant, for future me:** an invariant is a *sentence*, not a variable. Not "low
and high" — "if the target is in the array, it's between low and high". It has to be (a)
about your state, (b) true at the top of every iteration, (c) enough that invariant + exit
condition hands you the answer. "The length of the array" fails (c) and isn't something the
loop works to preserve, so it tells you nothing.

Strictly, this problem needs the weaker version: *"if a match exists further left than
`found`, it's within `[low, high]`"*. The plain sentence goes false the moment you record a
match and discard everything to the right of it.
