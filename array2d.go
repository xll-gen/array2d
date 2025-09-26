//go:build go1.18
// +build go1.18

// Package array2d contains an implementation of a 2D array.
package array2d

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrOutOfBounds is returned when an index is outside the array's bounds.
	ErrOutOfBounds = errors.New("array2d: index out of bounds")

	// ErrShape is returned when the dimensions of a source do not match the
	// specified height and width during array creation.
	ErrShape = errors.New("array2d: invalid shape for creation")

	// ErrNilDest is returned by Scan when the destination pointer is nil.
	ErrNilDest = errors.New("array2d: destination for Scan cannot be nil")

	// ErrDestLength is returned by Scan when the destination slice has an incorrect length.
	ErrDestLength = errors.New("array2d: destination slice has incorrect length")
)

const (
	// printThreshold is the size at which array printing is summarized.
	printThreshold = 10
	// edgeItems is the number of array elements to show for each edge.
	edgeItems = 5
)

// New initializes a 2-dimensional array with all zero values.
// By default, it creates a row-major array.
// To create a column-major array, pass true as the optional colMajor argument.
func New[T any](height, width int, colMajor ...bool) Array2D[T] {
	isColMajor := false
	if len(colMajor) > 0 {
		isColMajor = colMajor[0]
	}
	return Array2D[T]{
		height:   height,
		width:    width,
		slice:    make([]T, width*height),
		colMajor: isColMajor,
	}
}

// NewFilled initializes a 2-dimensional array with a value.
// By default, it creates a row-major array.
// To create a column-major array, pass true as the optional colMajor argument.
func NewFilled[T any](height, width int, value T, colMajor ...bool) Array2D[T] {
	isColMajor := false
	if len(colMajor) > 0 {
		isColMajor = colMajor[0]
	}
	slice := make([]T, width*height)
	fill(slice, value)
	return Array2D[T]{
		height:   height,
		width:    width,
		slice:    slice,
		colMajor: isColMajor,
	}
}

// FromSlice creates a 2-dimensional array from the given slice. The length of
// the slice must be equal to height * width.
//
// Note: This function does not create a copy of the provided slice.
// Modifications to the original slice will affect the new Array2D instance.
//
// By default, it creates a row-major array.
// To create a column-major array, pass true as the optional colMajor argument.
func FromSlice[T any](height, width int, slice []T, colMajor ...bool) (Array2D[T], error) {
	isColMajor := false
	if len(colMajor) > 0 {
		isColMajor = colMajor[0]
	}
	if len(slice) != width*height {
		return Array2D[T]{}, fmt.Errorf("%w: slice length %d does not match height*width %d", ErrShape, len(slice), width*height)
	}
	return Array2D[T]{
		height:   height,
		width:    width,
		slice:    slice,
		colMajor: isColMajor,
	}, nil
}

// FromJagged creates a 2-dimensional array from a jagged slice.
// It returns an error if the dimensions of the jagged slice exceed the specified
// height or width.
//
// By default, it creates a row-major array.
// To create a column-major array, pass true as the optional colMajor argument.
func FromJagged[J ~[]S, S ~[]E, E any](height, width int, jagged J, colMajor ...bool) (Array2D[E], error) {
	isColMajor := false
	if len(colMajor) > 0 {
		isColMajor = colMajor[0]
	}
	if len(jagged) > height {
		return Array2D[E]{}, fmt.Errorf("%w: jagged slice height %d exceeds specified height %d", ErrShape, len(jagged), height)
	}
	arr := New[E](height, width, isColMajor)
	for y, row := range jagged {
		if len(row) > width {
			return Array2D[E]{}, fmt.Errorf("%w: row %d width %d exceeds specified width %d", ErrShape, y, len(row), width)
		}
		if isColMajor {
			for x, val := range row {
				arr.setUnchecked(y, x, val)
			}
		} else {
			r, _ := arr.Row(y)
			copy(r, row)
		}
	}
	return arr, nil
}

// Array2D is a 2-dimensional array.
type Array2D[T any] struct {
	height, width int
	slice         []T
	colMajor      bool
}

// String returns a string representation of this array.
func (a Array2D[T]) String() string {
	var t T
	typeName := reflect.TypeOf(t).Name()
	if typeName == "" {
		typeName = reflect.TypeOf(t).String()
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Array2d[%s] %dx%d ", typeName, a.height, a.width)

	if a.height == 0 || a.width == 0 {
		sb.WriteString("[]")
		return sb.String()
	}

	summarizeRows := a.height > printThreshold
	summarizeCols := a.width > printThreshold

	sb.WriteByte('[')

	for y := 0; y < a.height; y++ {
		if summarizeRows && y == edgeItems {
			if y > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString("...")
			y = a.height - edgeItems - 1 // The loop will increment to a.height - edgeItems
			continue
		}

		if y > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('[')

		for x := 0; x < a.width; x++ {
			if summarizeCols && x == edgeItems {
				if x > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString("...")
				x = a.width - edgeItems - 1 // The loop will increment to a.width - edgeItems
				continue
			}

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
// It returns the zero value for T and false if the access is out-of-bounds.
func (a Array2D[T]) Get(row, col int) (T, bool) {
	if col < 0 || col >= a.width || row < 0 || row >= a.height {
		var zero T
		return zero, false
	}
	return a.getUnchecked(row, col), true
}

func (a Array2D[T]) getUnchecked(row, col int) T {
	if a.colMajor {
		return a.slice[row+col*a.height]
	}
	return a.slice[col+row*a.width]
}

// Set sets a value in the array.
// It returns an error on out-of-bounds access.
func (a Array2D[T]) Set(row, col int, value T) error {
	if col < 0 || col >= a.width {
		return fmt.Errorf("%w: col index %d out of range for width %d", ErrOutOfBounds, col, a.width)
	}
	if row < 0 || row >= a.height {
		return fmt.Errorf("%w: row index %d out of range for height %d", ErrOutOfBounds, row, a.height)
	}
	a.setUnchecked(row, col, value)
	return nil
}

func (a Array2D[T]) setUnchecked(row, col int, value T) {
	if a.colMajor {
		a.slice[row+col*a.height] = value
	} else {
		a.slice[col+row*a.width] = value
	}
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
		height:   a.height,
		width:    a.width,
		slice:    slice,
		colMajor: a.colMajor,
	}
}

// Row returns a mutable slice for an entire row. Changing values in this slice
// will affect the array.
//
// For column-major arrays, this function returns a new slice containing a copy
// of the data, so modifications to it will not affect the original array.
func (a Array2D[T]) Row(row int) ([]T, bool) {
	if row < 0 || row >= a.height {
		return nil, false
	}
	if a.colMajor {
		r := make([]T, a.width)
		for c := 0; c < a.width; c++ {
			r[c] = a.getUnchecked(row, c)
		}
		return r, true
	}
	return a.slice[row*a.width : a.width+row*a.width], true
}

// Col returns a slice for an entire column.
//
// For row-major arrays, this function returns a new slice containing a copy
// of the data, so modifications to it will not affect the original array.
//
// For column-major arrays, this function returns a mutable slice. Changing
// values in this slice will affect the array.
func (a Array2D[T]) Col(col int) ([]T, bool) {
	if col < 0 || col >= a.width {
		return nil, false
	}
	if a.colMajor {
		start := col * a.height
		return a.slice[start : start+a.height], true
	}
	c := make([]T, a.height)
	for r := 0; r < a.height; r++ {
		c[r] = a.getUnchecked(r, col)
	}
	return c, true
}

// Fill will assign all values inside the region to the specified value.
// The coordinates are inclusive, meaning all values from [x1,y1] including
// [x1,y1] to [x2,y2] including [x2,y2] are set.
//
// The method sorts the arguments, so col2 may be lower than col1 and row2 may be
// lower than row1.
func (a Array2D[T]) Fill(row1, col1, row2, col2 int, value T) error {
	if col1 < 0 || col1 >= a.width {
		return fmt.Errorf("%w: col1 index %d out of range for width %d", ErrOutOfBounds, col1, a.width)
	}
	if row1 < 0 || row1 >= a.height {
		return fmt.Errorf("%w: row1 index %d out of range for height %d", ErrOutOfBounds, row1, a.height)
	}
	if col2 < 0 || col2 >= a.width {
		return fmt.Errorf("%w: col2 index %d out of range for width %d", ErrOutOfBounds, col2, a.width)
	}
	if row2 < 0 || row2 >= a.height {
		return fmt.Errorf("%w: row2 index %d out of range for height %d", ErrOutOfBounds, row2, a.height)
	}

	if a.colMajor {
		// For simplicity, fill cell by cell for column-major.
		// This can be optimized if needed.
		for r := row1; r <= row2; r++ {
			for c := col1; c <= col2; c++ {
				a.setUnchecked(r, c, value)
			}
		}
		return nil
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
	return nil
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

// Rows returns an iterator over the rows of the array, similar to sql.Rows.
func (a *Array2D[T]) Rows() *Rows[T] {
	return &Rows[T]{
		arr: a,
		row: -1,
	}
}

// Rows is an iterator over the rows of an Array2D.
type Rows[T any] struct {
	arr *Array2D[T]
	row int
	err error
}

// Next advances the iterator to the next row.
// It returns false when the iteration is complete.
func (r *Rows[T]) Next() bool {
	if r.row+1 >= r.arr.height {
		return false
	}
	r.row++
	return true
}

// Scan copies the current row's data into the provided destination slice.
// The destination slice must have a length equal to the array's width.
func (r *Rows[T]) Scan(dest *[]T) error {
	if r.err != nil {
		return r.err
	}
	if dest == nil {
		r.err = ErrNilDest
		return r.err
	}
	if len(*dest) != r.arr.width {
		r.err = fmt.Errorf("%w: destination slice has length %d, but array width is %d", ErrDestLength, len(*dest), r.arr.width)
		return r.err
	}

	// Optimization: avoid intermediate slice allocation by copying directly.
	if r.arr.colMajor {
		// For column-major, row elements are not contiguous. Copy element by element.
		for c := 0; c < r.arr.width; c++ {
			(*dest)[c] = r.arr.getUnchecked(r.row, c)
		}
	} else {
		// For row-major, row elements are contiguous. A single copy is efficient.
		sourceRow, _ := r.arr.Row(r.row)
		copy(*dest, sourceRow)
	}

	return nil
}

// Err returns the error, if any, that was encountered during iteration.
func (r *Rows[T]) Err() error {
	return r.err
}

// Cols returns an iterator over the columns of the array, similar to sql.Rows.
func (a *Array2D[T]) Cols() *Cols[T] {
	return &Cols[T]{
		arr: a,
		col: -1,
	}
}

// Cols is an iterator over the columns of an Array2D.
type Cols[T any] struct {
	arr *Array2D[T]
	col int
	err error
}

// Next advances the iterator to the next column.
// It returns false when the iteration is complete.
func (c *Cols[T]) Next() bool {
	if c.col+1 >= c.arr.width {
		return false
	}
	c.col++
	return true
}

// Scan copies the current column's data into the provided destination slice.
// The destination slice must have a length equal to the array's height.
func (c *Cols[T]) Scan(dest *[]T) error {
	if c.err != nil {
		return c.err
	}
	if dest == nil {
		c.err = ErrNilDest
		return c.err
	}
	if len(*dest) != c.arr.height {
		c.err = fmt.Errorf("%w: destination slice has length %d, but array height is %d", ErrDestLength, len(*dest), c.arr.height)
		return c.err
	}

	// Optimization: avoid intermediate slice allocation by copying directly.
	if !c.arr.colMajor {
		// For row-major, column elements are not contiguous. Copy element by element.
		for r := 0; r < c.arr.height; r++ {
			(*dest)[r] = c.arr.getUnchecked(r, c.col)
		}
	} else {
		// For column-major, column elements are contiguous. A single copy is efficient.
		sourceCol, _ := c.arr.Col(c.col)
		copy(*dest, sourceCol)
	}

	return nil
}

// Err returns the error, if any, that was encountered during iteration.
func (c *Cols[T]) Err() error {
	return c.err
}
