//go:build go1.18
// +build go1.18

package array2d

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestArray2D_stringEmpty(t *testing.T) {
	arr := New[int](3, 3)
	got := arr.String()
	want := "Array2d[int] 3x3 [[0 0 0] [0 0 0] [0 0 0]]"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestArray2D_stringValues(t *testing.T) {
	arr := New[int](3, 3)
	_ = arr.Set(0, 0, 1)
	_ = arr.Set(0, 1, 2)
	_ = arr.Set(0, 2, 3)
	_ = arr.Set(1, 0, 4)
	_ = arr.Set(1, 1, 5)
	_ = arr.Set(1, 2, 6)
	_ = arr.Set(2, 0, 7)
	_ = arr.Set(2, 1, 8)
	_ = arr.Set(2, 2, 9)
	got := arr.String()
	want := "Array2d[int] 3x3 [[1 2 3] [4 5 6] [7 8 9]]"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestArray2D_stringSummarized(t *testing.T) {
	t.Run("summarize rows", func(t *testing.T) {
		arr := New[int](12, 3)
		for i := 0; i < arr.Height(); i++ {
			for j := 0; j < arr.Width(); j++ {
				_ = arr.Set(i, j, i*10+j)
			}
		}
		got := arr.String()
		want := "Array2d[int] 12x3 [[0 1 2] [10 11 12] [20 21 22] [30 31 32] [40 41 42] ... [70 71 72] [80 81 82] [90 91 92] [100 101 102] [110 111 112]]"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("summarize cols", func(t *testing.T) {
		arr := New[int](3, 12)
		for i := 0; i < arr.Height(); i++ {
			for j := 0; j < arr.Width(); j++ {
				_ = arr.Set(i, j, i*100+j)
			}
		}
		got := arr.String()
		want := "Array2d[int] 3x12 [[0 1 2 3 4 ... 7 8 9 10 11] [100 101 102 103 104 ... 107 108 109 110 111] [200 201 202 203 204 ... 207 208 209 210 211]]"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("summarize both", func(t *testing.T) {
		arr := New[int](11, 11)
		got := arr.String()
		want := "Array2d[int] 11x11 [[0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] ... [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0] [0 0 0 0 0 ... 0 0 0 0 0]]"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})
}

func TestArray2D_fill(t *testing.T) {
	arr := New[int](64, 64)
	val := 42
	if err := arr.Fill(25, 20, 38, 40, val); err != nil {
		t.Fatalf("Fill returned an unexpected error: %v", err)
	}
	for x := 0; x < arr.Width(); x++ {
		for y := 0; y < arr.Height(); y++ {
			want := 0
			if x >= 20 && x <= 40 && y >= 25 && y <= 38 {
				want = val
			}
			got, _ := arr.Get(y, x)
			if got != want {
				t.Errorf("x=%d, y=%d: want %d, got %d", x, y, want, got)
			}
		}
	}
}

func TestArray2D_row(t *testing.T) {
	arr := New[int](5, 5)
	span, ok := arr.Row(2)
	if !ok {
		t.Fatal("Row(2) returned ok=false unexpectedly")
	}
	assertLen(t, 5, span)
	copy(span, []int{1, 2, 3, 4, 5})
	for x := 0; x < arr.Width(); x++ {
		for y := 0; y < arr.Height(); y++ {
			want := 0
			if y == 2 {
				want = x + 1
			}
			got, _ := arr.Get(y, x)
			if got != want {
				t.Errorf("x=%d, y=%d: want %d, got %d", x, y, want, got)
			}
		}
	}
}

func TestArray2D_col(t *testing.T) {
	t.Run("row-major", func(t *testing.T) {
		arr := New[int](3, 4)
		// [[0 1 2 3]
		//  [4 5 6 7]
		//  [8 9 10 11]]
		for i := 0; i < arr.Height(); i++ {
			for j := 0; j < arr.Width(); j++ {
				_ = arr.Set(i, j, i*arr.Width()+j)
			}
		}

		col, ok := arr.Col(1)
		if !ok {
			t.Fatal("Col(1) returned ok=false unexpectedly")
		}
		assertLen(t, arr.Height(), col)

		want := []int{1, 5, 9}
		if !reflect.DeepEqual(col, want) {
			t.Errorf("unexpected column values: want %v, got %v", want, col)
		}

		// Modify the returned slice and ensure the original array is not affected
		for i := range col {
			col[i] *= 2
		}
		wantOriginal := []int{1, 5, 9}
		if reflect.DeepEqual(col, wantOriginal) {
			t.Errorf("the original array was affected by modifying the column")
		}
	})

	t.Run("column-major", func(t *testing.T) {
		arr := New[int](3, 4, true) // colMajor = true
		// [[0  3  6  9]
		//  [1  4  7 10]
		//  [2  5  8 11]]
		for i := 0; i < arr.Height(); i++ {
			for j := 0; j < arr.Width(); j++ {
				_ = arr.Set(i, j, i+j*arr.Height())
			}
		}

		col, ok := arr.Col(1)
		if !ok {
			t.Fatal("Col(1) returned ok=false unexpectedly")
		}
		assertLen(t, arr.Height(), col)

		want := []int{3, 4, 5}
		if !reflect.DeepEqual(col, want) {
			t.Errorf("unexpected column values: want %v, got %v", want, col)
		}

		// Modify the returned slice and ensure the original array IS affected
		col[0] = 42
		col[1] = 43
		col[2] = 44

		wantModified := []int{42, 43, 44}
		if !reflect.DeepEqual(col, wantModified) {
			t.Errorf("col slice was not correctly modified")
		}

		v1, _ := arr.Get(0, 1)
		v2, _ := arr.Get(1, 1)
		v3, _ := arr.Get(2, 1)
		if v1 != 42 || v2 != 43 || v3 != 44 {
			t.Errorf("original array was not modified through the col slice")
		}
	})
}

func TestArray2D_rows(t *testing.T) {
	arr := New[int](3, 4)
	// [[0 1 2 3]
	//  [4 5 6 7]
	//  [8 9 10 11]]
	for i := 0; i < arr.Height(); i++ {
		for j := 0; j < arr.Width(); j++ {
			_ = arr.Set(i, j, i*arr.Width()+j)
		}
	}

	rows := arr.Rows()
	if rows.Index() != -1 {
		t.Errorf("initial rows.Index() want -1, got %d", rows.Index())
	}
	defer func() {
		if err := rows.Err(); err != nil {
			t.Errorf("unexpected error during rows iteration: %v", err)
		}
	}()

	rowIndex := 0
	for rows.Next() {
		if rows.Index() != rowIndex {
			t.Errorf("rows.Index() want %d, got %d", rowIndex, rows.Index())
		}
		row := make([]int, arr.Width())
		if err := rows.Scan(&row); err != nil {
			t.Fatalf("error scanning row: %v", err)
		}

		want := []int{rowIndex * 4, rowIndex*4 + 1, rowIndex*4 + 2, rowIndex*4 + 3}
		if !reflect.DeepEqual(row, want) {
			t.Errorf("row %d: want %v, got %v", rowIndex, want, row)
		}
		rowIndex++
	}

	if rowIndex != arr.Height() {
		t.Errorf("rows iterator should have visited %d rows, but visited %d", arr.Height(), rowIndex)
	}
}

func TestArray2D_cols(t *testing.T) {
	arr := New[int](3, 4)
	// [[0 1 2 3]
	//  [4 5 6 7]
	//  [8 9 10 11]]
	for i := 0; i < arr.Height(); i++ {
		for j := 0; j < arr.Width(); j++ {
			_ = arr.Set(i, j, i*arr.Width()+j)
		}
	}

	cols := arr.Cols()
	if cols.Index() != -1 {
		t.Errorf("initial cols.Index() want -1, got %d", cols.Index())
	}
	colIndex := 0

	for cols.Next() {
		if cols.Index() != colIndex {
			t.Errorf("cols.Index() want %d, got %d", colIndex, cols.Index())
		}
		col := make([]int, arr.Height())
		if err := cols.Scan(&col); err != nil {
			t.Fatalf("error scanning col: %v", err)
		}

		want := []int{colIndex, colIndex + 4, colIndex + 8}
		if !reflect.DeepEqual(col, want) {
			t.Errorf("col %d: want %v, got %v", colIndex, want, col)
		}
		colIndex++
	}

	if colIndex != arr.Width() {
		t.Errorf("cols iterator should have visited %d cols, but visited %d", arr.Width(), colIndex)
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

		want := "Array2d[int] 2x3 [[1 2 3] [4 5 6]]"
		if got := arr.String(); got != want {
			t.Errorf("want %q, got %q", want, got)
		}

		// Test that modifying the original slice affects the array
		slice[0] = 99
		got, _ := arr.Get(0, 0)
		if got != 99 {
			t.Errorf("modifying original slice did not affect array, want 99, got %d", got)
		}
	})

	t.Run("length mismatch", func(t *testing.T) {
		slice := []int{1, 2, 3}
		_, err := FromSlice(2, 2, slice)
		if err == nil {
			t.Fatal("FromSlice() did not return an error for mismatched length")
		}
		if !errors.Is(err, ErrShape) {
			t.Errorf("want error to be ErrShape, but it was not. got: %v", err)
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

		want := "Array2d[int] 2x3 [[1 2 0] [3 4 5]]"
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
		if !errors.Is(err, ErrShape) {
			t.Errorf("want error to be ErrShape, but it was not. got: %v", err)
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
		if !errors.Is(err, ErrShape) {
			t.Errorf("want error to be ErrShape, but it was not. got: %v", err)
		}
	})
}

func TestArray2D_ToSlices(t *testing.T) {
	t.Run("row-major zero-copy", func(t *testing.T) {
		arr := New[int](2, 3)
		_ = arr.Set(0, 0, 1)
		_ = arr.Set(0, 1, 2)
		_ = arr.Set(1, 2, 6)

		slices := arr.ToSlices()
		want := [][]int{{1, 2, 0}, {0, 0, 6}}
		if !reflect.DeepEqual(slices, want) {
			t.Errorf("ToSlices() got = %v, want %v", slices, want)
		}

		// Modify slice and check if original array is affected
		slices[0][2] = 3
		got, _ := arr.Get(0, 2)
		if got != 3 {
			t.Errorf("modification on slice did not affect original array. got %d, want 3", got)
		}
	})

	t.Run("column-major copy", func(t *testing.T) {
		arr := New[int](2, 3, true) // colMajor
		_ = arr.Set(0, 0, 1)
		_ = arr.Set(0, 1, 2)
		_ = arr.Set(1, 2, 6)

		slices := arr.ToSlices()
		want := [][]int{{1, 2, 0}, {0, 0, 6}}
		if !reflect.DeepEqual(slices, want) {
			t.Errorf("ToSlices() got = %v, want %v", slices, want)
		}

		// Modify slice and check if original array is NOT affected
		slices[0][2] = 3
		got, _ := arr.Get(0, 2)
		if got != 0 {
			t.Errorf("modification on slice affected original array. got %d, want 0", got)
		}
	})
}

func TestArray2D_ToSlicesByCol(t *testing.T) {
	t.Run("row-major copy", func(t *testing.T) {
		arr := New[int](3, 2)
		_ = arr.Set(0, 0, 1)
		_ = arr.Set(1, 0, 2)
		_ = arr.Set(2, 1, 6)

		slices := arr.ToSlicesByCol()
		want := [][]int{{1, 2, 0}, {0, 0, 6}}
		if !reflect.DeepEqual(slices, want) {
			t.Errorf("ToSlicesByCol() got = %v, want %v", slices, want)
		}

		// Modify slice and check if original array is NOT affected
		slices[0][2] = 3
		got, _ := arr.Get(2, 0)
		if got != 0 {
			t.Errorf("modification on slice affected original array. got %d, want 0", got)
		}
	})

	t.Run("column-major zero-copy", func(t *testing.T) {
		arr := New[int](3, 2, true) // colMajor
		_ = arr.Set(0, 0, 1)
		_ = arr.Set(1, 0, 2)
		_ = arr.Set(2, 1, 6)

		slices := arr.ToSlicesByCol()
		want := [][]int{{1, 2, 0}, {0, 0, 6}}
		if !reflect.DeepEqual(slices, want) {
			t.Errorf("ToSlicesByCol() got = %v, want %v", slices, want)
		}

		// Modify slice and check if original array is affected
		slices[0][2] = 3
		got, _ := arr.Get(2, 0)
		if got != 3 {
			t.Errorf("modification on slice did not affect original array. got %d, want 3", got)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		arr := New[int](2, 3)
		// [[0 1 2]
		//  [3 4 5]]
		for i := 0; i < arr.Height(); i++ {
			for j := 0; j < arr.Width(); j++ {
				_ = arr.Set(i, j, i*arr.Width()+j)
			}
		}

		mappedArr := Map(arr, func(v int) string {
			return fmt.Sprintf("v%d", v)
		})

		want := "Array2d[string] 2x3 [[v0 v1 v2] [v3 v4 v5]]"
		if got := mappedArr.String(); got != want {
			t.Errorf("Map() result incorrect.\nwant: %s\ngot:  %s", want, got)
		}
	})
}
