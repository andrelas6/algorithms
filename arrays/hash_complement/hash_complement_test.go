package hashcomplement

import (
	"math/rand"
	"testing"
)

// anyPairExists is the obviously-correct oracle: does SOME valid pair exist?
func anyPairExists(durations []int32, slot int32) bool {
	for i := 0; i < len(durations); i++ {
		for j := i + 1; j < len(durations); j++ {
			if durations[i]+durations[j] == slot {
				return true
			}
		}
	}
	return false
}

// check validates the CONTRACT rather than one specific answer, since the
// statement says any valid pair is accepted.
func check(t *testing.T, name string, durations []int32, slot int32, got []int32) {
	t.Helper()

	if len(got) != 2 {
		t.Errorf("%s: got %v, want a slice of length 2", name, got)
		return
	}

	want := anyPairExists(durations, slot)

	if !want {
		if got[0] != -1 || got[1] != -1 {
			t.Errorf("%s: got %v, but no valid pair exists — want [-1 -1]", name, got)
		}
		return
	}

	i, j := got[0], got[1]
	n := int32(len(durations))

	switch {
	case i == -1 || j == -1:
		t.Errorf("%s: got %v, but a valid pair DOES exist", name, got)
	case i < 0 || j < 0 || i >= n || j >= n:
		t.Errorf("%s: got %v, indices out of range for n=%d", name, got, n)
	case i >= j:
		t.Errorf("%s: got %v, but the statement requires i < j", name, got)
	case durations[i]+durations[j] != slot:
		t.Errorf("%s: got %v — durations[%d]+durations[%d] = %d, want %d",
			name, got, i, j, durations[i]+durations[j], slot)
	}
}

func TestFindTaskPairForSlot(t *testing.T) {
	cases := []struct {
		name      string
		durations []int32
		slot      int32
	}{
		// examples from the statement
		{"example 1", []int32{2, 7, 11, 15}, 9},
		{"example 2, no pair", []int32{1, 2, 3, 4}, 8},

		// samples
		{"empty array", []int32{}, 10},
		{"single element", []int32{5}, 5},

		// the answer starts at index 0 — breaks a `!= 0` sentinel
		{"pair at 0 and 1", []int32{1, 3}, 4},
		{"pair at 0 and last", []int32{1, 9, 9, 9, 3}, 4},

		// must not pair an element with itself
		{"single element, half the target", []int32{3}, 6},
		{"element is half the target, no partner", []int32{3, 1, 2}, 6},
		{"two equal halves, valid", []int32{3, 3}, 6},

		// pair at the very end
		{"pair at the end", []int32{9, 9, 9, 1, 5}, 6},

		// duplicates everywhere
		{"all identical, valid", []int32{4, 4, 4, 4}, 8},
		{"all identical, invalid", []int32{4, 4, 4, 4}, 9},

		// no pair
		{"target too small", []int32{5, 6, 7}, 3},
		{"target too large", []int32{5, 6, 7}, 100},

		// larger values, per the constraints
		{"max-ish values", []int32{1000000, 1000000}, 2000000},
		{"max-ish values, no pair", []int32{1000000, 999999}, 2000000},
	}

	for _, c := range cases {
		check(t, c.name, c.durations, c.slot, findTaskPairForSlot(c.durations, c.slot))
	}
}

func TestFindTaskPairForSlotRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))

	for trial := 0; trial < 5000; trial++ {
		n := r.Intn(8)
		durations := make([]int32, n)
		for i := range durations {
			durations[i] = int32(1 + r.Intn(10))
		}
		slot := int32(1 + r.Intn(22))

		check(t, "random", durations, slot, findTaskPairForSlot(durations, slot))
	}
}
