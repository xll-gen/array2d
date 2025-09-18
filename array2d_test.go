//go:build go1.18
// +build go1.18

package array2d

import (
	"testing"
)

func TestArray2D_stringEmpty(t *testing.T) {
	arr := New[int](3, 3)
	got := arr.String()
	want := "[[0 0 0] [0 0 0] [0 0 0]]"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestArray2D_stringValues(t *testing.T) {
	arr := New[int](3, 3)
	arr.Set(0, 0, 1)
	arr.Set(1, 0, 2)
	arr.Set(2, 0, 3)
	arr.Set(0, 1, 4)
	arr.Set(1, 1, 5)
	arr.Set(2, 1, 6)
	arr.Set(0, 2, 7)
	arr.Set(1, 2, 8)
	arr.Set(2, 2, 9)
	got := arr.String()
	want := "[[1 2 3] [4 5 6] [7 8 9]]"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestArray2D_fill(t *testing.T) {
	arr := New[int](64, 64)
	val := 42
	arr.Fill(20, 25, 40, 38, val)
	for x := 0; x < arr.Width(); x++ {
		for y := 0; y < arr.Height(); y++ {
			want := 0
			if x >= 20 && x <= 40 && y >= 25 && y <= 38 {
				want = val
			}
			got := arr.Get(x, y)
			if got != want {
				t.Errorf("x=%d, y=%d: want %d, got %d", x, y, want, got)
			}
		}
	}
}

func TestArray2D_rowSpan(t *testing.T) {
	arr := New[int](5, 5)
	span := arr.RowSpan(1, 3, 2)
	assertLen(t, 3, span)
	copy(span, []int{1, 2, 3})
	for x := 0; x < arr.Width(); x++ {
		for y := 0; y < arr.Height(); y++ {
			want := 0
			if x >= 1 && x <= 3 && y == 2 {
				want = x
			}
			got := arr.Get(x, y)
			if got != want {
				t.Errorf("x=%d, y=%d: want %d, got %d", x, y, want, got)
			}
		}
	}
}

func TestArray2D_row(t *testing.T) {
	arr := New[int](5, 5)
	span := arr.Row(2)
	assertLen(t, 5, span)
	copy(span, []int{1, 2, 3, 4, 5})
	for x := 0; x < arr.Width(); x++ {
		for y := 0; y < arr.Height(); y++ {
			want := 0
			if y == 2 {
				want = x + 1
			}
			got := arr.Get(x, y)
			if got != want {
				t.Errorf("x=%d, y=%d: want %d, got %d", x, y, want, got)
			}
		}
	}
}

func TestFromSlice(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6}
		arr, err := FromSlice(2, 3, slice) // height, width
		if err != nil {
			t.Fatalf("FromSlice() returned an unexpected error: %v", err)
		}

		if arr.Width() != 3 || arr.Height() != 2 {
			t.Errorf("want width=3, height=2, got width=%d, height=%d", arr.Width(), arr.Height())
		}

		want := "[[1 2 3] [4 5 6]]"
		if got := arr.String(); got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		// Test that modifying the original slice affects the array
		slice[0] = 99
		if got := arr.Get(0, 0); got != 99 {
			t.Errorf("modifying original slice did not affect array, want 99, got %d", got)
		}
	})

	t.Run("length mismatch", func(t *testing.T) {
		slice := []int{1, 2, 3}
		_, err := FromSlice(2, 2, slice)
		if err == nil {
			t.Fatal("FromSlice() did not return an error for mismatched length")
		}
		want := "array2d: slice length 3 does not match height*width 4"
		if err.Error() != want {
			t.Errorf("want error %q, got %q", want, err.Error())
		}
	})
}

func assertLen[E any](t *testing.T, want int, slice []E) {
	t.Helper()
	if len(slice) != want {
		t.Errorf("want len %d, got len %d", want, len(slice))
	}
}

func TestFromJagged(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		jagged := [][]int{
			{1, 2},
			{3, 4, 5},
		}
		arr, err := FromJagged(2, 3, jagged)
		if err != nil {
			t.Fatalf("FromJagged() returned an unexpected error: %v", err)
		}

		if arr.Width() != 3 || arr.Height() != 2 {
			t.Errorf("want width=3, height=2, got width=%d, height=%d", arr.Width(), arr.Height())
		}

		want := "[[1 2 0] [3 4 5]]"
		if got := arr.String(); got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("height exceeds specified height", func(t *testing.T) {
		jagged := [][]int{
			{1},
			{2},
			{3},
		}
		_, err := FromJagged(2, 1, jagged)
		if err == nil {
			t.Fatal("FromJagged() did not return an error for exceeding height")
		}
		want := "array2d: jagged slice height 3 exceeds specified height 2"
		if err.Error() != want {
			t.Errorf("want error %q, got %q", want, err.Error())
		}
	})

	t.Run("width exceeds specified width", func(t *testing.T) {
		jagged := [][]int{
			{1},
			{2, 3},
		}
		_, err := FromJagged(2, 1, jagged)
		if err == nil {
			t.Fatal("FromJagged() did not return an error for exceeding width")
		}
		want := "array2d: row 1 width 2 exceeds specified width 1"
		if err.Error() != want {
			t.Errorf("want error %q, got %q", want, err.Error())
		}
	})
}
