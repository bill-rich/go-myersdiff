# Go-MyersDiff
This implementation of the Myers Diff Algorythm was written for when the diff of many files is needed fast in Go. It
includes options for what output is needed (additions, deletions, or similarities).

## Usage
```go
import "github.com/bill-rich/go-myersdiff"

func main() {
	a := []string {
		"line 1",
		"line 2a",
		"line 3",
	}
	a := []string {
		"line 1", 
		"line 2b", 
		"line 3",
	}
	fmt.Println(myersdiff.GenerateDiff(a, b, myersdiff.NewOptions()))
}
```