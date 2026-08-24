# Height of Binary Search Tree

**Source:** HackerRank · **Pattern:** Tree DFS by recursion (post-order combine)

## Problem

Return the height of the tree — **the number of nodes** along the longest root-to-leaf path.

The tree comes as three parallel arrays, not pointers. Node `i` has value `values[i]`, left
child at index `leftChild[i]`, right child at index `rightChild[i]`. `-1` means no child.
**Index 0 is the root.**

```
values     = [4, 2, 6, 1, 3, 5, 7]
leftChild  = [1, 3, 5, -1, -1, -1, -1]
rightChild = [2, 4, 6, -1, -1, -1, -1]
-> 3
```

### Constraints

- `0 <= n <= 1000`, all `values[i]` unique
- children are `-1` or a valid index
- BST ordering holds (left < node < right)
- the edges form a single connected acyclic graph

### Samples

| n | tree | Output |
|---|---|---|
| 1 | single node | 1 |
| 2 | root + one left child | 2 |
| 0 | empty | 0 |

## Recognition signals

- **Tree + "longest path" / "depth" / "height"** → DFS recursion, combining results from the
  children.
- The shape is always the same: **base case, recurse on both children, combine.** Height uses
  `max`; count uses `+`; "does X exist below here" uses `||`.
- **The BST property is irrelevant here.** Height doesn't care about ordering. Noticing which
  constraints you don't need is its own skill — say it out loud in an interview.

## Idea

The recursive contract: `height(i)` returns the number of nodes on the longest path from node
`i` down to a leaf.

```
height(i):
    if i == -1: return 0                                  # no node, no path
    return 1 + max(height(left[i]), height(right[i]))     # me, plus the taller side

return height(0)
```

That's the whole algorithm. `n = 0` returns 0 because there's no node 0 to start at — guard
it before recursing.

Complexity: **O(n) time** (every node visited once), **O(h) space** for the call stack, which
is O(n) in the worst case (a degenerate tree that's really a linked list). Go grows goroutine
stacks dynamically, so 1000 frames is a non-issue.

## The iterative version (do this too)

Level-order BFS with a queue: push the root, then repeatedly process the whole current level
and count how many levels you got through. That count *is* the height. O(n) time, O(w) space
where `w` is the widest level.

Worth writing both — Nebius's own prep guidance names tree traversals **and** queues, and
interviewers often ask for the iterative form after you give the recursive one.

## Traps

- **Height is counted in NODES, not edges.** A single node is height 1, not 0. Plenty of
  textbooks use the edge definition; the statement here is explicit.
- `n == 0` → 0. Guard before recursing on index 0.
- `-1` is the "no child" sentinel, and it doubles as the base case. Don't index with it.
- Both children, every time. Recursing on one and forgetting the other passes balanced trees
  and fails everything else.
- Root is index 0 by convention here. If a problem doesn't say, the root is the index that
  never appears in `leftChild` or `rightChild`.

## Notes

2026-08-20. **Did not solve this one — needed the solution.** Two attempts, both stuck on the
same thing. Redo cold on Aug 21, not Aug 23.

**The blocker: I kept recursing on `index + 1`.**

```go
height(..., leftChild[index])   // right — jump to where the child lives
height(..., index+1)            // wrong — step to the next slot in storage
```

`leftChild[index]` *is* the child's address. You look it up and go there. `index+1` is the
next box in memory, which has nothing to do with the tree's shape. The tree lives in the
**contents** of the arrays, not their order. The `scrambled indices` test case exists to make
that fail loudly: the real path is `0 → 2 → 4`.

Second blocker: **base case on "am I a leaf?" instead of "is there a node here?"** Checking
`leftChild[i] == -1 && rightChild[i] == -1` reads the arrays before knowing `i` is valid, so
the empty tree panics, and it forces the `+1` outside the recursion. `index == -1 → 0` fixes
both — a leaf then falls out as `1 + max(0, 0) = 1` with no special case. Same lesson as
`lastEnd = -1` in `intervals/greedy_by_end/` and `found = -1` in `arrays/lower_bound/`:
**pick the sentinel so the special cases stop being special.**

Also: `values` is never needed. Height depends on shape only, not on what's stored.

**Contract vs mechanism, 4th time.** Asked what the recursive call promises, I answered
"returns 1 if the node is present" — that's what *one node contributes*, not what the *call
returns*. The test for a contract: is it strong enough that `1 + max(left, right)` works if
the children keep it? "Returns the subtree's height" passes; "returns 1 if present" gives at
most 2. Start the sentence with **"returns…"**, never with "if".
