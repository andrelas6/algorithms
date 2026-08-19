# Find the Smallest Missing Positive Integer

**Source:** HackerRank · **Pattern:** Cyclic sort

## Problem

Given an unsorted array of integers, find the smallest positive integer not present in the
array, in **O(n)** time and **O(1)** extra space.

```
orderNumbers = [3, 4, -1, 1]
output = 2
```

## Recognition signals

- The answer must lie in `1..n+1`, so the values that matter map onto array indices.
- O(1) extra space rules out a hash set — that forces an in-place approach.
- "Smallest missing from a `1..n` range" is the canonical cyclic-sort tell.

## Idea

Place every value `v` where `1 <= v <= n` at index `v-1` by swapping. Then scan once: the
first index `i` where `arr[i] != i+1` gives the answer `i+1`. If every slot is correct, the
answer is `n+1`.

Complexity: O(n) time (each swap puts one value permanently home), O(1) space.

## Traps

- **Duplicates cause an infinite loop.** Swapping a value into a slot that already holds the
  same value makes no progress. Guard with "if the target slot already holds this value,
  move on" — that check is in the implementation.
- Skip out-of-range values (`v < 1` or `v > n`) instead of swapping them.
- Empty array → answer is 1.
- Only advance `i` when no swap happened; a successful swap must re-examine the same index.

## Notes
