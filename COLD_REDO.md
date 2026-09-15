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
| 4 | **Product of Array Except Self** | 2026-09-14 — brute force came easily, the improvement needed every hint | 15 min | **2026-09-17** | |
| 5 | **Valid Binary Search Tree** | 2026-09-15 — had only a brute-force idea and never wrote it; needed the full solution walkthrough | 15 min | **2026-09-18** | |
| 6 | **Kth Largest Element in an Array** | 2026-09-15 — the easy version came quickly; the faster one needed the full answer and still slows to a crawl on some large inputs | 20 min | **2026-09-18** | |
| 7 | **Permutation in String** | 2026-09-15 — found the idea without help, but the implementation was shaky and had leftover steps it didn't need | 20 min | **2026-09-18** | |

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

---

## 4. Product of Array Except Self

Given an integer array `nums`, return an array `output` where `output[i]` is the product of
all elements of `nums` except `nums[i]`.

Each product is guaranteed to fit in a 32-bit integer.

**Input format:** `nums` (INTEGER_ARRAY).

**Constraints:**

- `1 <= nums.length <= 1000`
- `-20 <= nums[i] <= 20`

**Output format:** an INTEGER_ARRAY of the same length as `nums`.

```
Input:  [1,2,4,6]        Output: [48,24,12,8]
Input:  [-1,0,1,2,3]     Output: [0,-6,0,0,0]
Input:  [1,2,3,4]        Output: [24,12,8,6]
```

Follow-up worth attempting on the redo: solve it without using the division operator, in
O(n) time, and with O(1) extra space — the returned array does not count towards the space.

Edge cases to think about: a single element, exactly two elements, one zero anywhere in the
array, two or more zeroes, all zeroes, all negative values with an odd and an even count,
mixed signs, and every element being 1.

---

## 5. Valid Binary Search Tree

Given the `root` of a binary tree, return `true` if it is a valid binary search tree,
otherwise return `false`.

A valid binary search tree satisfies all of the following:

- The left subtree of every node contains only nodes with keys **strictly less** than that
  node's key.
- The right subtree of every node contains only nodes with keys **strictly greater** than
  that node's key.
- Both the left and right subtrees are also binary search trees.

**Input format:** `root` (TreeNode), given in level order where `null` marks an absent
child.

**Constraints:**

- The number of nodes is in the range `[1, 10^4]`
- `-2^31 <= Node.val <= 2^31 - 1`

**Output format:** a single BOOLEAN.

```
Input:  [2,1,3]                  Output: true
Input:  [5,1,4,null,null,3,6]    Output: false
Input:  [1,2,3]                  Output: false
```

Edge cases to think about: a single node, a node with only one child on either side, equal
values anywhere in the tree, a node that looks fine next to its parent but not next to a
node higher up, values at the very limits of the allowed range, and a tree that is one long
path to the left or to the right.

---

## 6. Kth Largest Element in an Array

Given an integer array `nums` and an integer `k`, return the `k`th largest element in the
array.

Note that it is the `k`th largest element in sorted order, not the `k`th distinct element.

**Input format:** `nums` (INTEGER_ARRAY) and `k` (INTEGER).

**Constraints:**

- `1 <= k <= nums.length <= 10^5`
- `-10^4 <= nums[i] <= 10^4`

**Output format:** a single INTEGER.

```
Input:  nums = [2,3,1,5,4], k = 2              Output: 4
Input:  nums = [2,3,1,1,5,5,4], k = 3          Output: 4
Input:  nums = [3,2,1,5,6,4], k = 2            Output: 5
Input:  nums = [3,2,3,1,2,4,5,5,6], k = 4      Output: 4
```

Follow-up worth attempting on the redo: can you solve it without sorting the whole array?

Edge cases to think about: a single element, `k = 1`, `k` equal to the length, duplicates
at and around the answer, every value equal, negatives, and input that arrives already in
ascending or descending order at the maximum length.

---

## 7. Permutation in String

Given two strings `s1` and `s2`, return `true` if `s2` contains a permutation of `s1`, or
`false` otherwise.

In other words, return `true` if one of `s1`'s permutations is a substring of `s2`.

**Input format:** two STRING parameters, `s1` and `s2`.

**Constraints:**

- `1 <= s1.length, s2.length <= 10^4`
- `s1` and `s2` consist of lowercase English letters

**Output format:** a single BOOLEAN.

```
Input:  s1 = "abc", s2 = "lecabee"     Output: true
Input:  s1 = "abc", s2 = "lecaabee"    Output: false
Input:  s1 = "ab",  s2 = "eidbaooo"    Output: true
Input:  s1 = "ab",  s2 = "eidboaoo"    Output: false
```

Edge cases to think about: `s1` longer than `s2`, both strings the same length, a match at
the very start or the very end of `s2`, letters repeated in `s1`, the right letters with the
wrong counts, the right letters split apart by other letters, and how much work each
position of `s2` costs you.
