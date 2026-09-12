# Cold Redo Queue

Problems to re-solve **from a blank file**, with no notes, no README, and no folder names.

Rules:

- Do not open the problem's folder until you are done. The folder name is the answer.
- Start a timer. If you blow past the target, stop and re-read the folder's `## Notes`.
- Passing cold means: correct on the first run, within the target time. Then log it and
  schedule the next repeat (3 days → 1 week → 1 month).

| # | Problem | Solved | Target | Redo on | Result |
|---|---|---|---|---|---|
| 1 | **Longest Substring Without Duplicate Characters** | 2026-09-08 — first version correct but too slow; needed hints, then the answer, to fix it | 15 min | **2026-09-11** | |
| 2 | **Merge Intervals** | 2026-09-11 — landed on the right idea after a nudge; lost time to an out-of-bounds index | 15 min | **2026-09-14** | |
| 3 | **Encode and Decode Strings** | 2026-09-12 — overcomplicated the first pass with a structure that lost order and duplicates | 15 min | **2026-09-15** | |

---

## 1. Longest Substring Without Duplicate Characters

Given a string `s`, find the length of the longest substring without duplicate characters.

A substring is a contiguous run of characters. The empty string has length 0.

**Input format:** a single STRING parameter, `s`.

**Constraints:**

- `0 <= s.length <= 50,000`
- `s` may consist of printable ASCII characters (codes 32 to 126 inclusive)
- `s` may be empty

**Output format:** a single integer — the length of the longest substring of `s` that
contains no repeated character.

```
Input:  zxyzxyz    Output: 3
Input:  xxxx       Output: 1
Input:  abcabcbb   Output: 3
Input:  pwwkew     Output: 3
Input:  (empty)    Output: 0
```

Edge cases to think about: empty string, single character, every character identical, every
character distinct, spaces and symbols (not just letters), the answer sitting at the very
start or the very end of the string, and a repeat whose earlier copy has already fallen out
of consideration.

---

## 2. Merge Intervals

Given an array of intervals where `intervals[i] = [start_i, end_i]`, merge all overlapping
intervals and return an array of the non-overlapping intervals that cover all the intervals
in the input.

You may return the answer in any order.

Note: intervals are non-overlapping if they have no common point. `[1, 2]` and `[3, 4]` are
non-overlapping, but `[1, 2]` and `[2, 3]` are overlapping.

**Input format:** `intervals` (INTEGER_ARRAY of INTEGER_ARRAY), each inner array holding
exactly two values.

**Constraints:**

- `1 <= intervals.length <= 1000`
- `intervals[i].length == 2`
- `0 <= start <= end <= 1000`

**Output format:** an array of two-element arrays — the merged, non-overlapping intervals.

```
Input:  [[1,3],[1,5],[6,7]]            Output: [[1,5],[6,7]]
Input:  [[1,2],[2,3]]                  Output: [[1,3]]
Input:  [[1,3],[2,6],[8,10],[15,18]]   Output: [[1,6],[8,10],[15,18]]
```

Edge cases to think about: a single interval, input arriving in any order, intervals that
touch at exactly one endpoint, one interval fully containing another, several intervals
sharing the same start, zero-width intervals like `[4,4]`, exact duplicates, nothing
overlapping at all, everything collapsing into one, and whether the caller's input may be
modified.

---

## 3. Encode and Decode Strings

Design an algorithm to encode a list of strings into a single string. The encoded string is
then sent over the network and decoded back into the original list of strings.

Machine 1 (sender) has a function `encode(strs)` returning a single string. Machine 2
(receiver) has a function `decode(encoded)` returning the list. The list Machine 2 produces
must be the same as the list Machine 1 was given.

**Input format:** `strs` (STRING_ARRAY) for encode; a single STRING for decode.

**Constraints:**

- `0 <= strs.length < 100`
- `0 <= strs[i].length < 200`
- `strs[i]` may contain any UTF-8 characters

**Output format:** encode returns a single STRING; decode returns a STRING_ARRAY equal to
the original input.

```
Input:  ["neet","code","love","you"]    Output: ["neet","code","love","you"]
Input:  ["we","say",":","yes"]          Output: ["we","say",":","yes"]
Input:  ["Hello","World"]               Output: ["Hello","World"]
```

Edge cases to think about: an empty list, a list holding one empty string, several empty
strings, duplicate strings, whether order survives, strings containing whatever character
you chose to structure the output with, strings that look like an already-encoded payload,
strings long enough that any count you write is more than one character wide, and non-ASCII
text where the number of characters and the number of bytes disagree.
