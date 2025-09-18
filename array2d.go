//go:build go1.18
// +build go1.18

// Package array2d contains an implementation of a 2D array.
package array2d

import (
	"fmt"
	"strings"
)

// New initializes a 2-dimensional array with all zero values.
func New[T any](height, width int) Array2D[T] {
	return Array2D[T]{
		height: height,
		width:  width,
		slice:  make([]T, width*height),
	}
}

// NewFilled initializes a 2-dimensional array with a value.
func NewFilled[T any](height, width int, value T) Array2D[T] {
	slice := make([]T, width*height)
	fill(slice, value)
	return Array2D[T]{
		height: height,
		width:  width,
		slice:  slice,
	}
}

// FromSlice creates a 2-dimensional array from the given slice. The length of
// the slice must be equal to height * width.
//
// Note: This function does not create a copy of the provided slice.
// Modifications to the original slice will affect the new Array2D instance.
func FromSlice[T any](height, width int, slice []T) (Array2D[T], error) {
	if len(slice) != width*height {
		return Array2D[T]{}, fmt.Errorf("array2d: slice length %d does not match height*width %d", len(slice), width*height)
	}
	return Array2D[T]{
		height: height,
		width:  width,
		slice:  slice,
	}, nil
}

// FromJagged creates a 2-dimensional array from a jagged slice.
// It returns an error if the dimensions of the jagged slice exceed the specified
// height or width.
func FromJagged[J ~[]S, S ~[]E, E any](height, width int, jagged J) (Array2D[E], error) {
	if len(jagged) > height {
		return Array2D[E]{}, fmt.Errorf("array2d: jagged slice height %d exceeds specified height %d", len(jagged), height)
	}
	arr := New[E](height, width)
	for y, row := range jagged {
		if len(row) > width {
			return Array2D[E]{}, fmt.Errorf("array2d: row %d width %d exceeds specified width %d", y, len(row), width)
		}
		copy(arr.Row(y), row)
	}
	return arr, nil
}

// Array2D is a 2-dimensional array.
type Array2D[T any] struct {
	height, width int
	slice         []T
}

// Iterator returns an iterator for the 2D array.
// The iterator allows to range over all elements of the array, similar to sql.Rows.
//
// Example:
//
//	it := arr.Iterator()
//	for it.Next() {
//	    row, col, val := it.Value()
//	    // ...
//	}
func (a *Array2D[T]) Iterator() *Iter[T] {
	return &Iter[T]{
		arr: a,
		idx: -1,
	}
}

// Iter is an iterator for the Array2D.
type Iter[T any] struct {
	arr *Array2D[T]
	idx int
}

// Next advances the iterator to the next element.
// It returns false when the iteration is complete.
func (it *Iter[T]) Next() bool {
	it.idx++
	return it.idx < len(it.arr.slice)
}

// Value returns the current row, column, and value.
func (it *Iter[T]) Value() (row, col int, value T) {
	row = it.idx / it.arr.width
	col = it.idx % it.arr.width
	value = it.arr.slice[it.idx]
	return
}

// RowIterator returns an iterator for a specific row of the 2D array.
func (a *Array2D[T]) RowIterator(row int) *RowIter[T] {
	if row < 0 || row >= a.height {
		panic(fmt.Sprintf("array2d: row index out of range [%d] with height %d", row, a.height))
	}
	return &RowIter[T]{
		arr: a,
		row: row,
		col: -1,
	}
}

// RowIter is an iterator for the rows of an Array2D.
type RowIter[T any] struct {
	arr      *Array2D[T]
	row, col int
}

// Next advances the iterator to the next row.
// It returns false when the iteration is complete.
func (it *RowIter[T]) Next() bool {
	it.col++
	return it.col < it.arr.width
}

// Value returns the current row index and a slice representing the row.
func (it *RowIter[T]) Value() (col int, value T) {
	return it.col, it.arr.getUnchecked(it.row, it.col)
}

// ColIterator returns an iterator for the columns of the 2D array.
func (a *Array2D[T]) ColIterator(col int) *ColIter[T] {
	if col < 0 || col >= a.width {
		panic(fmt.Sprintf("array2d: col index out of range [%d] with width %d", col, a.width))
	}
	return &ColIter[T]{
		arr: a,
		col: col,
		row: -1,
	}
}

// ColIter is an iterator for the columns of an Array2D.
type ColIter[T any] struct {
	arr      *Array2D[T]
	col, row int
}

// Next advances the iterator to the next column.
// It returns false when the iteration is complete.
func (it *ColIter[T]) Next() bool {
	it.row++
	return it.row < it.arr.height
}

// Value returns the current column index and a new slice containing the elements of that column.
func (it *ColIter[T]) Value() (row int, value T) {
	return it.row, it.arr.getUnchecked(it.row, it.col)
}

// String returns a string representation of this array.
func (a Array2D[T]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for y := 0; y < a.height; y++ {
		if y > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('[')
		for x := 0; x < a.width; x++ {
			if x > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, a.getUnchecked(y, x))
		}
		sb.WriteByte(']')
	}
	sb.WriteByte(']')
	return sb.String()
}

// Get returns a value from the array.
//
// The function will panic on out-of-bounds access.
func (a Array2D[T]) Get(row, col int) T {
	if col < 0 || col >= a.width {
		panic(fmt.Sprintf("array2d: col index out of range [%d] with width %d", col, a.width))
	}
	if row < 0 || row >= a.height {
		panic(fmt.Sprintf("array2d: row index out of range [%d] with height %d", row, a.height))
	}
	return a.getUnchecked(row, col)
}

func (a Array2D[T]) getUnchecked(row, col int) T {
	return a.slice[col+row*a.width]
}

// Set sets a value in the array.
//
// The function will panic on out-of-bounds access.
func (a Array2D[T]) Set(row, col int, value T) {
	if col < 0 || col >= a.width {
		panic(fmt.Sprintf("array2d: col index out of range [%d] with width %d", col, a.width))
	}
	if row < 0 || row >= a.height {
		panic(fmt.Sprintf("array2d: row index out of range [%d] with height %d", row, a.height))
	}
	a.setUnchecked(row, col, value)
}

func (a Array2D[T]) setUnchecked(row, col int, value T) {
	a.slice[col+row*a.width] = value
}

// Width returns the width of this array. The maximum x value is Width()-1.
func (a Array2D[T]) Width() int {
	return a.width
}

// Height returns the height of this array. The maximum y value is Height()-1.
func (a Array2D[T]) Height() int {
	return a.height
}

// Copy returns a shallow copy of this array.
func (a Array2D[T]) Copy() Array2D[T] {
	slice := make([]T, len(a.slice))
	copy(slice, a.slice)
	return Array2D[T]{
		height: a.height,
		width:  a.width,
		slice:  slice,
	}
}

// RowSpan returns a mutable slice for part of a row. Changing values in this
// slice will affect the array.
func (a Array2D[T]) RowSpan(row, col1, col2 int) []T {
	if row < 0 || row >= a.height {
		panic(fmt.Sprintf("array2d: row index out of range [%d] with height %d", row, a.height))
	}
	if col1 < 0 || col1 >= a.width {
		panic(fmt.Sprintf("array2d: col1 index out of range [%d] with width %d", col1, a.width))
	}
	if col2 < 0 || col2 >= a.width {
		panic(fmt.Sprintf("array2d: col2 index out of range [%d] with width %d", col2, a.width))
	}
	return a.slice[col1+row*a.width : 1+col2+row*a.width]
}

// Row returns a mutable slice for an entire row. Changing values in this slice
// will affect the array.
func (a Array2D[T]) Row(row int) []T {
	if row < 0 || row >= a.height {
		panic(fmt.Sprintf("array2d: row index out of range [%d] with height %d", row, a.height))
	}
	return a.slice[row*a.width : a.width+row*a.width]
}

// Fill will assign all values inside the region to the specified value.
// The coordinates are inclusive, meaning all values from [x1,y1] including
// [x1,y1] to [x2,y2] including [x2,y2] are set.
//
// The method sorts the arguments, so col2 may be lower than col1 and row2 may be
// lower than row1.
func (a Array2D[T]) Fill(row1, col1, row2, col2 int, value T) {
	if col1 < 0 || col1 >= a.width {
		panic(fmt.Sprintf("array2d: col1 index out of range [%d] with width %d", col1, a.width))
	}
	if row1 < 0 || row1 >= a.height {
		panic(fmt.Sprintf("array2d: row1 index out of range [%d] with height %d", row1, a.height))
	}
	if col2 < 0 || col2 >= a.width {
		panic(fmt.Sprintf("array2d: col2 index out of range [%d] with width %d", col2, a.width))
	}
	if row2 < 0 || row2 >= a.height {
		panic(fmt.Sprintf("array2d: row2 index out of range [%d] with height %d", row2, a.height))
	}
	if col2 < col1 {
		col1, col2 = col2, col1
	}
	if row2 < row1 {
		row1, row2 = row2, row1
	}
	firstRow := a.slice[col1+row1*a.width : 1+col2+row1*a.width]
	fill(firstRow, value)
	for row := row1 + 1; row <= row2; row++ {
		copy(a.slice[col1+row*a.width:1+col2+row*a.width], firstRow)
	}
}

func fill[E any](slice []E, value E) {
	if len(slice) == 0 {
		return
	}
	// Exponential copy to fill a slice
	slice[0] = value
	for i := 1; i < len(slice); i += i {
		copy(slice[i:], slice[:i])
	}
}
