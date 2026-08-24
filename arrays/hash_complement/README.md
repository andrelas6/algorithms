# Two Sum

**Source:** HackerRank · **Pattern:** Hash map — remember what you've seen (complement lookup)

## Problem

Given an array of positive integers and a target, return the indices of two elements summing
to the target, or `[-1, -1]`.

```
taskDurations = [2, 7, 11, 15], slotLength = 9   -> [0, 1]
taskDurations = [1, 2, 3, 4],   slotLength = 8   -> [-1, -1]
```

### Constraints

- `0 <= n <= 1000`
- `1 <= taskDurations[i] <= 1_000_000` (all **positive**)
- `1 <= slotLength <= 2_000_000`
- must return `0 <= i < j < n` — **distinct, ascending**
- any valid pair is accepted if several exist

### Samples

| n | input | Output |
|---|---|---|
| 0 | `[]`, 10 | -1 -1 |
| 1 | `[5]`, 5 | -1 -1 |

## Recognition signals

- **"find two things that sum to X"** → hash map complement lookup. The canonical one.
- More generally: *"have I already seen the thing I need?"* → hash map. Don't search for the
  partner, **remember what you've walked past**.
- `n <= 1000` means O(n²) also passes here — so the brute force is not *wrong*, it just isn't
  the answer they're looking for. Say the O(n) version out loud even if you write the loop.

## Idea

One pass. At index `i` holding value `x`, the partner you need is `slotLength - x`. Ask the
map whether you've seen it already.

```
seen = {}                      # value -> index where it was seen
for i, x in enumerate(durations):
    need = slotLength - x
    if need in seen:
        return [seen[need], i]   # seen[need] < i, so ascending is automatic
    seen[x] = i
return [-1, -1]
```

Complexity: **O(n) time, O(n) space.** The brute force is O(n²) time, O(1) space.

## Traps

- **`i < j` is required.** Returning `[j, i]` fails. Checking the map *before* inserting the
  current value gives you ascending order for free — the stored index is always older.
- **Don't pair an element with itself.** Insert `x` only *after* checking for its complement;
  otherwise `[3]` with target 6 matches index 0 against itself.
- **`0` is a valid index, so it cannot be your "not found" sentinel.** Use `-1`, or a
  separate bool.
- Empty and single-element inputs → `[-1, -1]`. Loop never fires, or fires once with an empty
  map.
- Duplicate values are fine: `[3, 3]` with target 6 → `[0, 1]`. The map holds the *first*
  index for a value, which is all you need.

## Notes

2026-08-21. Brute force took four attempts; the hash map version worked first try once the
trick was known.

**The trick, in one line:** don't search for the partner — **remember what you've walked
past**. Not something to derive under pressure; something to know.

**Say the complexity precisely: O(n) *average*,** because map lookups are O(1) average, not
worst case. That chains straight into the "why are hash lookups O(1)?" answer Nebius
reportedly asks — uniform hash, resizing bounds the chain, Go randomises the seed.

The four brute-force attempts, each a different flavour of the same bug — **the number I have
is measured against a different thing than the number I need**:

1. `if v1 != 0` as the "found" flag — index 0 is a *valid answer*, so found and not-found were
   indistinguishable. Use `-1`, or better, `return` from inside the loop and delete the flag.
2. Inner loop from `j := 0` — allowed `i == j` (element paired with itself) and descending
   pairs.
3. `for j, y := range taskDurations[i+1:]` — **re-slicing renumbers.** `s[k:]` shares memory
   but rebases index 0 to the original's `k`. The *value* was right, so the sum check passed
   and it misreported the index. True index is `i + 1 + j`. Don't re-slice and translate —
   index explicitly.
4. `for j := 1` — anchored to a constant instead of to `i`.

Same tell as `index+1` in `trees/dfs_recursion/` and `len(code)` vs `codeLowercase` in
`strings/two_pointers/`.

**Control flow beat the flags.** `v1`, `v2`, both `break`s and the outer check all vanished
once the match returned directly. When you add a variable to remember "did I find it", check
whether you can just return.

Caveat on the tests here: they validate the *contract* (any `i < j` pair summing to target),
because the statement says any valid pair is accepted. If HackerRank exact-matches one
expected answer, these tests will pass while the judge fails — brute force returns the pair
with the smallest `i`, the hash map returns the pair with the smallest `j`, and both are
"correct". e.g. `[1,2,4,6,3,9]` target 10 → brute `[0,5]`, hash `[2,3]`.
