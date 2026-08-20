# Maximum Number of Non-Overlapping Intervals

**Source:** HackerRank · **Pattern:** Greedy — sort by end time (activity selection)

## Problem

Given intervals with a start and end, return the maximum number you can pick that don't
overlap.

```
[[1,2],[2,3],[3,4],[1,3]]                  -> 3
[[0,5],[0,1],[1,2],[2,3],[3,5],[4,6]]      -> 4
```

Touching is fine: `[1,2]` and `[2,3]` do **not** overlap.

### Constraints

- `0 <= n <= 1000`
- `0 <= start < end <= 10^9`
- meetings may share start or end times

### Samples

| meetings | Output |
|---|---|
| `[[5,10]]` | 1 |
| `[[1,2],[2,3],[3,4]]` | 3 |

## Recognition signals

- **Intervals + "maximum how many"** → greedy, sorted by end time. This one is a named
  classic (*activity selection*); learn it, don't re-derive it.
- "Maximum number of X" with no weights/values attached → greedy. If each interval had a
  *value* and you wanted max total value, greedy breaks and it becomes DP.
- `n <= 1000` — the sort dominates, so O(n log n). Nothing here needs to be clever.

## Idea

Sort by **end** time. Walk left to right, take an interval whenever it starts at or after the
last one you took ended.

```
sort meetings by end ascending
count = 0
lastEnd = -infinity

for each [start, end] in meetings:
    if start >= lastEnd:
        count++
        lastEnd = end

return count
```

**Why end time?** Taking the interval that finishes earliest leaves the most room for
everything after it. Formally: if some optimal answer exists, you can always swap its first
interval for the earliest-finishing one without breaking anything — the swap can't create an
overlap, since the replacement ends no later. Repeat that argument and the greedy choice is
optimal at every step.

Complexity: **O(n log n)** time (the sort; the scan is O(n)), **O(1)** extra space beyond
the sort.

## Traps

- **`start >= lastEnd`, not `>`.** `[1,2]` and `[2,3]` are both selectable — Example 1 depends
  on this.
- **Sorting by start is wrong.** `[[0,5],[1,2],[3,4]]` → sorted by start you take `[0,5]` and
  get 1. Sorting by end gives 2.
- **Sorting by duration is also wrong.** `[[0,10],[9,11],[10,20]]` → shortest is `[9,11]`,
  which blocks both others for a total of 1. Answer is 2.
- Empty input → 0. The loop never runs, `count` stays 0.
- `lastEnd` must start below every possible start. `0 <= start`, so `-1` works; so does
  taking the first interval before the loop.
- Sorting mutates the caller's slice. Doesn't matter for the judge, worth knowing.

## Notes

Solved 2026-08-19 after two failed attempts. All tests pass, including a fuzz test against an
exhaustive subset search.

**The first failure was comprehension, not code.** I read it as "count the intervals that
don't overlap" rather than "choose the largest subset that is mutually compatible". You get
to *discard* the troublemakers. The statement said so: Example 1 has 4 meetings and answers
3, so one was dropped. **Check your restatement against the given examples before coding** —
if your reading doesn't reproduce their output, your reading is wrong.

Then two state bugs:

1. **Sorted `copy`, then read `meetings` inside the loop.** The sort was thrown away. Same
   shape as the palindrome bug (size a pointer from one string, index another). Recurring
   personal tell: derive from one thing, read from another.
2. **Invented `lastDiscardedIndex` to reconstruct "the last one I kept".** Signal to watch
   for: *if you're adding a variable to reconstruct something you already had, the framing is
   wrong.* The index was right there — it's `i`. The detour existed because I was counting
   discards and subtracting instead of counting what I take.

Also: don't update `lastEnd` when you **skip**. Only when you take. Skipping leaves the
previous kept interval as the thing that blocks the next one.

**Init matters.** Starting with `count = 1` and `lastEnd = copy[0][1]` forces a
`len <= 1` guard, because you have to read `copy[0]`. Starting with `count = 0` and
`lastEnd = -1` handles every meeting in one uniform loop and needs no guard.
