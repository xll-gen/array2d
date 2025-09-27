# array2d

**Warning: This package is currently under development and the API is subject to change. It should be considered unstable.**

Fork of [github.com/fufuok/utils/generic/array2d](https://github.com/fufuok/utils/commit/382fc5c9e91e33694350885335a3155e8f787959).

```go
import "github.com/xll-gen/array2d"
```

Package array2d contains an implementation of a 2D array.

## Index

- [array2d](#array2d)
	- [Index](#index)
	- [type Array2D](#type-array2d)
		- [func New](#func-new)
		- [func NewFilled](#func-newfilled)
		- [func FromSlice](#func-fromslice)
		- [func FromJagged](#func-fromjagged)
		- [func (Array2D\[T\]) Row](#func-array2dt-row)
		- [func (Array2D\[T\]) Col](#func-array2dt-col)
		- [func (Array2D\[T\]) Fill](#func-array2dt-fill)
		- [func (Array2D\[T\]) Get](#func-array2dt-get)
		- [func (Array2D\[T\]) Set](#func-array2dt-set)
		- [func (Array2D\[T\]) Copy](#func-array2dt-copy)
		- [func (Array2D\[T\]) String](#func-array2dt-string)
		- [func (Array2D\[T\]) Height](#func-array2dt-height)
		- [func (Array2D\[T\]) Width](#func-array2dt-width)
		- [func (Array2D\[T\]) ToSlices](#func-array2dt-toslices)
		- [func (Array2D\[T\]) ToSlicesByCol](#func-array2dt-toslicesbycol)
		- [func (\*Array2D\[T\]) Rows](#func-array2dt-rows)
		- [func (\*Rows\[T\]) Index](#func-rowst-index)
		- [func (\*Array2D\[T\]) Cols](#func-array2dt-cols)
		- [func (\*Cols\[T\]) Index](#func-colst-index)
		- [func Map](#func-map)
	- [License](#license)

## type Array2D

Array2D is a 2-dimensional array.

```go
type Array2D[T any] struct {
    // contains filtered or unexported fields
}
```

### func New

```go
func New[T any](height, width int, colMajor ...bool) Array2D[T]
```

New initializes a 2-dimensional array with all zero values.  
By default, it creates a row-major array.  
To create a column-major array, pass `true` as the optional `colMajor` argument.

### func NewFilled

```go
func NewFilled[T any](height, width int, value T, colMajor ...bool) Array2D[T]
```

NewFilled initializes a 2-dimensional array with a value.  
By default, it creates a row-major array.  
To create a column-major array, pass `true` as the optional `colMajor` argument.

### func FromSlice

```go
func FromSlice[T any](height, width int, slice []T, colMajor ...bool) (Array2D[T], error)
```

FromSlice creates a 2-dimensional array from the given slice. The length of the slice must be equal to height * width.

Note: This function does not create a copy of the provided slice. Modifications to the original slice will affect the new Array2D instance.

By default, it creates a row-major array.  
To create a column-major array, pass `true` as the optional `colMajor` argument.

### func FromJagged

```go
func FromJagged[J ~[]S, S ~[]E, E any](height, width int, jagged J, colMajor ...bool) (Array2D[E], error)
```

FromJagged creates a 2-dimensional array from a jagged slice. It returns an error if the dimensions of the jagged slice exceed the specified height or width.

By default, it creates a row-major array.  
To create a column-major array, pass `true` as the optional `colMajor` argument.

### func (Array2D[T]) Row

```go
func (a Array2D[T]) Row(row int) ([]T, bool)
```

Row returns a slice for an entire row.

- For row-major arrays, this function returns a mutable slice. Changing values in this slice will affect the array.
- For column-major arrays, this function returns a new slice containing a copy of the data, so modifications to it will not affect the original array.
- It returns `false` if the row index is out of bounds.

### func (Array2D[T]) Col

```go
func (a Array2D[T]) Col(col int) ([]T, bool)
```

Col returns a slice for an entire column.

- For column-major arrays, this function returns a mutable slice. Changing values in this slice will affect the array.
- For row-major arrays, this function returns a new slice containing a copy of the data, so modifications to it will not affect the original array.
- It returns `false` if the column index is out of bounds.

### func (Array2D[T]) Fill

```go
func (a Array2D[T]) Fill(row1, col1, row2, col2 int, value T) error
```

Fill will assign all values inside the region to the specified value. The coordinates are inclusive, meaning all values from [row1,col1] including [row1,col1] to [row2,col2] including [row2,col2] are set.

It returns an error if any of the coordinates are out of bounds.

### func (Array2D[T]) Get

```go
func (a Array2D[T]) Get(row, col int) (T, bool)
```

Get returns a value from the array.

It returns the zero value for T and `false` if the access is out-of-bounds.

### func (Array2D[T]) Set

```go
func (a Array2D[T]) Set(row, col int, value T) error
```

Set sets a value in the array.

It returns an error on out-of-bounds access.

### func (Array2D[T]) Copy

```go
func (a Array2D[T]) Copy() Array2D[T]
```

Copy returns a shallow copy of this array.

### func (Array2D[T]) String

```go
func (a Array2D[T]) String() string
```

String returns a string representation of this array.

### func (Array2D[T]) Height

```go
func (a Array2D[T]) Height() int
```

Height returns the height of this array. The maximum y value is Height()-1.

### func (Array2D[T]) Width

```go
func (a Array2D[T]) Width() int
```

Width returns the width of this array. The maximum x value is Width()-1.

### func (Array2D[T]) ToSlices

```go
func (a Array2D[T]) ToSlices() [][]T
```

ToSlices returns a slice of slices representation of the array, organized by rows.

- For row-major arrays, this is a zero-copy operation (sub-slices of the underlying array).
- For column-major arrays, this returns copies of each row (modifying the result does **not** affect the original array).

### func (Array2D[T]) ToSlicesByCol

```go
func (a Array2D[T]) ToSlicesByCol() [][]T
```

ToSlicesByCol returns a slice of slices representation of the array, organized by columns.

- For column-major arrays, this is a zero-copy operation (sub-slices of the underlying array).
- For row-major arrays, this returns copies of each column (modifying the result does **not** affect the original array).

### func (*Array2D[T]) Rows

```go
func (a *Array2D[T]) Rows() *Rows[T]
```

Rows returns an iterator over the rows of the array, similar to sql.Rows.

**Example:**
```go
rows := arr.Rows()
rowData := make([]int, arr.Width())
for rows.Next() {
    if err := rows.Scan(&rowData); err != nil {
        // handle error
    }
    // use rowData
}
```

### func (*Rows[T]) Index

```go
func (r *Rows[T]) Index() int
```

Index returns the current row index. It returns -1 if Next has not been called yet.

### func (*Array2D[T]) Cols

```go
func (a *Array2D[T]) Cols() *Cols[T]
```

Cols returns an iterator over the columns of the array, similar to sql.Rows.

**Example:**
```go
cols := arr.Cols()
colData := make([]int, arr.Height())
for cols.Next() {
    if err := cols.Scan(&colData); err != nil {
        // handle error
    }
    // use colData
}
```

### func (*Cols[T]) Index

```go
func (c *Cols[T]) Index() int
```

Index returns the current column index. It returns -1 if Next has not been called yet.

### func Map

```go
func Map[T any, U any](a Array2D[T], f func(v T) U) Array2D[U]
```

Map creates a new Array2D by applying a function to each element of the input array.  
The new array will have the same dimensions and memory layout (row/column-major) as the original.  
The mapping function `f` is applied to each element of type `T` to produce an element of type `U`.

**Example:**
```go
arr := array2d.NewFilled[int](2, 3, 1)
mapped := array2d.Map(arr, func(v int) string {
    return fmt.Sprintf("val:%d", v)
})
fmt.Println(mapped)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 xll-gen  
Copyright (c) 2021-2024 fufuok
