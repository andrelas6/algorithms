# Check Palindrome by Filtering Non-Letters

**Source:** HackerRank · **Pattern:** Two pointers (opposite ends)

## Problem

Given a string containing letters, digits, and symbols, determine if it reads the same
forwards and backwards when considering **only alphabetic characters**, case-insensitively.

```
code = "A1b2B!a"
output = 1        // letters -> A b B a -> lowercase -> abba -> palindrome
```

Return a boolean; the judge prints `1` for true and `0` for false.

### Constraints

- `0 <= len(code) <= 1000`
- every character is printable ASCII, `33 <= ASCII(code[i]) <= 126`

### Samples

| Input | Output |
|---|---|
| `Z` | 1 |
| `abc123cba` | 1 |

## Recognition signals

- "reads the same forwards and backwards" → opposite-ends two pointers.
- "ignore some characters" is not a new pattern — it's a *skip rule* on the pointer advance.
- Case-insensitivity is normalisation **at comparison time**, not a pre-pass.

## Approach

```
l = 0, r = len(code) - 1
while l < r:
    while l < r and code[l] is not a letter: l++
    while l < r and code[r] is not a letter: r--
    if lower(code[l]) != lower(code[r]): return false
    l++; r--
return true
```

O(n) time (pointers only move inward), **O(1) space**.

## Traps

- **Guard `l < r` inside the skip loops too**, or `"!!!"` runs off the end.
- Move **both** pointers after a match; moving one is an infinite loop.
- `""`, single char, and no-letters-at-all are all palindromes → default `true` is correct.
- Classify **then** fold: `| 0x20` on a non-letter mangles `@ [ \ ] ^ _`.

## Notes

Solved 2026-08-19 in 36 min, all tests passed. The four things to carry forward:

1. **Fold with a bitwise OR instead of `strings.ToLower`.** `const caseDiff = 'a' - 'A'`
   (== 32 == a single bit), then compare `code[l]|caseDiff != code[r]|caseDiff`. `|` is
   idempotent so it's safe on either case — `+` is not (`'a'+32` overflows past `'z'`).
   Valid only because both bytes are already known to be letters.
2. **Inner `for` loops to skip to the next letter, then compare.** My version used three
   sequential `if`s that happened not to double-increment; skip-first makes the mutual
   exclusion *structural* instead of incidental. This shape transfers to "valid palindrome
   II", `removeDuplicates`, and every skip-filtered pointer walk.
3. **Bug: pointer sized from `len(code)` but indexing `codeLowercase`.** `strings.ToLower`
   is not length-preserving in general. Always derive the index from the string you index.
4. **`strings.ToLower` made it O(n) space**, not the O(1) the pattern targets. Fine at
   `n <= 1000`, but name the trade-off out loud before being asked.

Redo cold without `strings.ToLower`, target 10 minutes.

**Neighbouring pattern:** fast/slow pointers (cycle detection) — same direction, different
speeds. Here they converge from opposite ends.
