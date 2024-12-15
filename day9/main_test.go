package main

import "testing"

func TestDiscMapFromString(t *testing.T) {
	tests := []struct {
		s       string
		discMap []int
	}{
		{"5724", []int{5, 7, 2, 4}},
	}

testsLoop:
	for _, test := range tests {
		discMap := DiscMapFromString(test.s)
		if len(discMap) != len(test.discMap) {
			t.Errorf("expected disc map of length %d, got %d", len(test.discMap), len(discMap))
			continue
		}
		for i, got := range discMap {
			expected := test.discMap[i]
			if got != expected {
				t.Errorf("expected %v, got %v", expected, got)
				continue testsLoop
			}
		}
	}
}

func TestReadFromBack(t *testing.T) {
	tests := []struct {
		discMap          []int
		id, length, read int
	}{
		{[]int{3, 4, 5}, 1, 5, 1},
		{[]int{3, 4, 5, 6, 7, 8}, 2, 7, 2},
	}

	for _, test := range tests {
		id, length, read := ReadFromBack(test.discMap)
		if id != test.id {
			t.Errorf("expected id %d, got %d", test.id, id)
		}
		if length != test.length {
			t.Errorf("expected length %d, got %d", test.length, length)
		}
		if read != test.read {
			t.Errorf("expected to read %d, got %d", test.read, read)
		}
	}
}

func TestComputeChecksum(t *testing.T) {
	tests := []struct {
		discMap []int
		expect  int
	}{
		{[]int{1, 2, 3}, 1 * (1 + 2 + 3)},
		{[]int{1, 2, 3, 5, 4}, 0*0 + (1+2)*2 + (3+4+5)*1 + (6+7)*2},
	}

	for _, test := range tests {
		got := ComputeChecksum(test.discMap)
		if got != test.expect {
			t.Errorf("expected %d, got %d", test.expect, got)
		}
	}
}

func TestMemorySpanStringConversions(t *testing.T) {
	tests := []struct {
		s        string
		expanded string
	}{
		{"123", "0..111"},
	}

	for _, test := range tests {
		spans := MemorySpansFromString(test.s)
		got := MemorySpansToString(spans)

		if got != test.expanded {
			t.Errorf("expected %s, got %s", test.expanded, got)
		}
	}
}

func TestReorderSpans(t *testing.T) {
	tests := []struct {
		s      string
		expect string
	}{
		{"123", "0..111"},
		{"133", "0111..."},
		{"13324", "0111.....2222"},
		{"654321", "00000022...1111......"},
	}

	for _, test := range tests {
		got := MemorySpansToString(ReorderSpans(MemorySpansFromString(test.s)))
		if got != test.expect {
			t.Errorf("expected %s (length=%d), got %s (length=%d)", test.expect, len(test.expect), got, len(got))
		}
	}
}

func TestComputeChecksumFromMemorySpans(t *testing.T) {
	tests := []struct {
		s        string
		checksum int
	}{
		{"123", (3 + 4 + 5) * 1},
		{"12324", (3+4+5)*1 + (8+9+10+11)*2},
	}

	for _, test := range tests {
		got := ComputeChecksumFromMemorySpans(MemorySpansFromString(test.s))
		if got != test.checksum {
			t.Errorf("expected %d, got %d", test.checksum, got)
		}
	}
}
