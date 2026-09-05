package pathlib_test

import (
	"fmt"

	"github.com/skalt/pathlib.go"
)

// Syntactic sugar.
func enforce(err error) {
	if err != nil {
		panic(err)
	}
}

// More syntactic sugar.
func expect[T any](val T, err error) T {
	enforce(err)
	return val
}

func Example() {
	tmp := pathlib.TempDir()
	dir := expect(tmp.Join("pathlib-example").As[pathlib.Dir]().Make(0o777))
	defer func() { expect(dir.RemoveAll()) }()

	onDisk := expect(dir.Stat())
	fmt.Printf("Created %s with mode %s\n", expect(dir.Rel(tmp)), onDisk.Mode())

	for i, subPath := range []string{"a.txt", "b.txt", "c/d.txt"} {
		file := dir.Join(subPath).As[pathlib.File]()
		handle := expect(file.MakeAll(0o666, 0o777))
		expect(fmt.Fprintf(handle, "%d", i))
		enforce(handle.Close())
		fmt.Printf(
			"contents of %s: %q\n",
			expect(file.Rel(tmp)),
			string(expect(file.Read())),
		)
	}

	fmt.Printf("contents of %s:\n", expect(dir.Rel(tmp)))
	for _, entry := range expect(dir.Read()) {
		fmt.Println("  - " + entry.Name())
	}
	// Output:
	// Created pathlib-example with mode drwxr-xr-x
	// contents of pathlib-example/a.txt: "0"
	// contents of pathlib-example/b.txt: "1"
	// contents of pathlib-example/c/d.txt: "2"
	// contents of pathlib-example:
	//   - a.txt
	//   - b.txt
	//   - c
}
