package main

import "github.com/ofabry/go-callvis/origin"

// Use go-callvis as a package instead of binary on user machine.
// It simplifies go-codevis usage start and allows to manage dependencies by developer,
// therefore makes behaviour more predictable for developer (different version compatibility etc).
func main() {
	origin.Run()
}
