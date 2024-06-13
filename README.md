# Go Debounce

[![Tests on Linux, MacOS and Windows](https://github.com/qpoint-io/debounce/workflows/Test/badge.svg)](https://github.com/qpoint-io/debounce/actions?query=workflow:Test)
[![GoDoc](https://godoc.org/github.com/qpoint-io/debounce?status.svg)](https://godoc.org/github.com/qpoint-io/debounce)

## Example

```go
func ExampleNew() {
	var counter uint64

	f := func() {
		atomic.AddUint64(&counter, 1)
	}

	debounced := debounce.New(100 * time.Millisecond, 50)

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			debounced(f)
		}

		time.Sleep(200 * time.Millisecond)
	}

	c := int(atomic.LoadUint64(&counter))

	fmt.Println("Counter is", c)
	// Output: Counter is 3
}
```

