# Count Elements Greater Than Previous Average

**Source:** HackerRank · **Pattern:** Running aggregate / prefix sum

## Problem

Given an array of positive integers, return the number of elements that are strictly greater
than the average of all previous elements. Skip the first element.

```
responseTimes = [100, 200, 150, 300]
output = 2
```

**Constraints**

```
0 <= responseTimes.length <= 1000
1 <= responseTimes[i] <= 10^9
```

## Recognition signals

- "Compare each element against something about *all previous* elements" → maintain a
  running aggregate instead of recomputing per index.
- Recomputing the average each step is O(n²); a running sum makes it O(n).

## Idea

Sweep left to right keeping `sum` and `count` of the elements seen so far. At index `i > 0`,
compare `a[i]` against `sum / count` before folding `a[i]` into the aggregate.

Complexity: O(n) time, O(1) space.

## Traps

- **Overflow.** 1000 elements × 10^9 exceeds int32 — accumulate in `int64`.
- Integer division truncates. Comparing `a[i] > sum/count` is not the same as
  `a[i]*count > sum`; the latter avoids the truncation entirely and is the safer form.
- Empty array and single-element array both return 0.

## Notes
