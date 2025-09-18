
# array2d

**Warning: This package is currently under development and the API is subject to change. It should be considered unstable.**

Fork of [github.com/fufuok/utils/generic/array2d](https://github.com/fufuok/utils/commit/382fc5c9e91e33694350885335a3155e8f787959).

```go
import "github.com/xll-gen/array2d"
```

Package array2d contains an implementation of a 2D array.

## Index

- [type Array2D](<#type-array2d>)
  - [func FromJagged[J ~[]S, S ~[]E, E any](height, width int, jagged J) (Array2D[E], error)](<#func-fromjagged>)
  - [func FromSlice[T any](height, width int, slice []T) (Array2D[T], error)](<#func-fromslice>)
  - [func New[T any](height, width int) Array2D[T]](<#func-new>)
  - [func NewFilled[T any](height, width int, value T) Array2D[T]](<#func-newfilled>)
  - [func (a Array2D[T]) Copy() Array2D[T]](<#func-array2dt-copy>)
  - [func (a Array2D[T]) Fill(row1, col1, row2, col2 int, value T)](<#func-array2dt-fill>)
  - [func (a Array2D[T]) Get(row, col int) T](<#func-array2dt-get>)
  - [func (a *Array2D[T]) Iterator() *Iter[T]](<#func-array2dt-iterator>)
  - [func (a *Array2D[T]) RowIterator(row int) *RowIter[T]](<#func-array2dt-rowiterator>)
  - [func (a *Array2D[T]) ColIterator(col int) *ColIter[T]](<#func-array2dt-coliterator>)
  - [func (a Array2D[T]) Height() int](<#func-array2dt-height>)
  - [func (a Array2D[T]) Row(row int) []T](<#func-array2dt-row>)
  - [func (a Array2D[T]) RowSpan(row, col1, col2 int) []T](<#func-array2dt-rowspan>)
  - [func (a Array2D[T]) Set(row, col int, value T)](<#func-array2dt-set>)
  - [func (a Array2D[T]) String() string](<#func-array2dt-string>)
  - [func (a Array2D[T]) Width() int](<#func-array2dt-width>)


## type Array2D

Array2D is a 2\-dimensional array.

```go
type Array2D[T any] struct {
    // contains filtered or unexported fields
}
```

<details><summary>Example</summary>
<p>

```go package
import (
	"fmt"
	"strings"

	"github.com/xll-gen/array2d"
)

type Sudoku struct {
	arr array2d.Array2D[byte]
}

func (s Sudoku) PrintBoard() {
	var sb strings.Builder
	for y := 0; y < s.arr.Height(); y++ {
		if y%3 == 0 {
			sb.WriteString("+-------+-------+-------+\n")
		}
		for x := 0; x < s.arr.Width(); x++ {
			if x%3 == 0 {
				sb.WriteString("| ")
			}
			val := s.arr.Get(y, x)
			if val == 0 {
				sb.WriteByte(' ')
			} else {
				fmt.Fprint(&sb, val)
			}
			sb.WriteByte(' ')
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("+-------+-------+-------+\n")
	fmt.Print(sb.String())
}

func ExampleArray2D() {
	arr, err := array2d.FromJagged(9, 9, [][]byte{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	})
	if err != nil {
		panic(err)
	}

	s := Sudoku{
		arr: arr,
	}

	s.arr.Set(5, 2, 3)

	s.PrintBoard()
}

```

#### Output

```
+-------+-------+-------+
| 5 3   |   7   |       |
| 6     | 1 9 5 |       |
|   9 8 |       |   6   |
+-------+-------+-------+
| 8     |   6   |     3 |
| 4     | 8   3 |     1 |
| 7   3 |   2   |     6 |
+-------+-------+-------+
|   6   |       | 2 8   |
|       | 4 1 9 |     5 |
|       |   8   |   7 9 |
+-------+-------+-------+
```

</p>
</details>

### func FromJagged

```go
func FromJagged[J ~[]S, S ~[]E, E any](height, width int, jagged J) (Array2D[E], error)
```

FromJagged creates a 2-dimensional array from a jagged slice. It returns an error if the dimensions of the jagged slice exceed the specified height or width.

### func FromSlice

```go
func FromSlice[T any](height, width int, slice []T) (Array2D[T], error)
```

FromSlice creates a 2-dimensional array from the given slice. The length of the slice must be equal to height * width.

Note: This function does not create a copy of the provided slice. Modifications to the original slice will affect the new Array2D instance.

### func New

```go
func New[T any](width, height int) Array2D[T]
```

New initializes a 2\-dimensional array with all zero values.

### func NewFilled

```go
func NewFilled[T any](width, height int, value T) Array2D[T]
```

NewFilled initializes a 2\-dimensional array with a value.

### func \(Array2D\[T\]\) Copy

```go
func (a Array2D[T]) Copy() Array2D[T]
```

Copy returns a shallow copy of this array.

### func \(Array2D\[T\]\) Fill

```go
func (a Array2D[T]) Fill(row1, col1, row2, col2 int, value T)
```

Fill will assign all values inside the region to the specified value. The coordinates are inclusive, meaning all values from [row1,col1] including [row1,col1] to [row2,col2] including [row2,col2] are set.

The method sorts the arguments, so col2 may be lower than col1 and row2 may be lower than row1.

### func \(Array2D\[T\]\) Get

```go
func (a Array2D[T]) Get(row, col int) T
```

Get returns a value from the array.

The function will panic on out\-of\-bounds access.

### func \(\*Array2D\[T\]\) Iterator

```go
func (a *Array2D[T]) Iterator() *Iter[T]
```

Iterator returns an iterator for the 2D array. The iterator allows to range over all elements of the array, similar to sql.Rows.

Example:
```go
it := arr.Iterator()
for it.Next() {
    row, col, val := it.Value()
    // ...
}
```

### func \(\*Array2D\[T\]\) RowIterator

```go
func (a *Array2D[T]) RowIterator(row int) *RowIter[T]
```

RowIterator returns an iterator for the rows of the 2D array.

### func \(Array2D\[T\]\) Height

```go
func (a Array2D[T]) Height() int
```

Height returns the height of this array. The maximum y value is Height\(\)\-1.

### func \(Array2D\[T\]\) Row

```go
func (a Array2D[T]) Row(y int) []T
```

Row returns a mutable slice for an entire row. Changing values in this slice will affect the array.

### func \(Array2D\[T\]\) RowSpan

```go
func (a Array2D[T]) RowSpan(row, col1, col2 int) []T
```

RowSpan returns a mutable slice for part of a row. Changing values in this slice will affect the array.

### func \(Array2D\[T\]\) Set

```go
func (a Array2D[T]) Set(row, col int, value T)
```

Set sets a value in the array.

The function will panic on out\-of\-bounds access.

### func \(Array2D\[T\]\) String

```go
func (a Array2D[T]) String() string
```

String returns a string representation of this array.

### func \(Array2D\[T\]\) Width

```go
func (a Array2D[T]) Width() int
```

Width returns the width of this array. The maximum x value is Width\(\)\-1.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 xll-gen
Copyright (c) 2021-2024 fufuok
