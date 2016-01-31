package main

import "net/http"

// main is the function that will be called when program starts, when main function exists program exists
func main() {

	// Phase 1 help:
	// call parseFlags
	// call list
	// call merge

	// Phase 2 help:
	// call watch

	// Phase 3 help:
	// casting to MergedFile
	// https://golang.org/pkg/net/http/#Handle
	// https://golang.org/pkg/net/http/#ListenAndServe
	// https://golang.org/pkg/sync/#WaitGroup

}

// cliParams is structure contains cli params flags
type cliParams struct {
	list string // phase 1
	out  string // phase 1
	// watch bool   // phase 2
	// serve int    // phase 3
}

// parseFlags will parse cli program arguments into internal structure for later use
func parseFlags() (params cliParams) {

	// Phase 1 help:
	// https://golang.org/pkg/flag/#StringVar

	// Phase 2 help:
	// https://golang.org/pkg/flag/#BoolVar

	// Phase 3 help:
	// https://golang.org/pkg/flag/#IntVar

	// help:
	// https://golang.org/pkg/flag/#Parse

	return
}

// list will read list of css files from list json file
func list(listFile string) (cssFilePaths []string, err error) {

	// help:
	// https://golang.org/pkg/io/ioutil/#ReadFile
	// https://golang.org/pkg/encoding/json/#Unmarshal

	return
}

// merge will merge css files into one big new file, if merged file exists it will be overwritten
func merge(cssFilePaths []string, mergedFile string) (err error) {

	// help:
	// https://golang.org/pkg/os/#Create
	// https://golang.org/pkg/os/#File.Close
	// https://golang.org/pkg/os/#Open
	// https://golang.org/pkg/io/#Copy

	return nil
}

// watch will watch changes in cssFilePaths files and rebuild mergedFile
func watch(cssFilePaths []string, mergedFile string) (err error) {

	// help: https://golang.org/pkg/os/#Stat

	return
}

// MergedFile is a custom type representing merged css file
type MergedFile string

// ServeHTTP of MergedFile type satisfy http.Handler interface making it accessible via http protocol
func (mf MergedFile) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// help:
	// https://golang.org/pkg/os/#Stat
	// https://golang.org/pkg/os/#Open
	// https://golang.org/pkg/net/http/#ServeContent
	// https://golang.org/pkg/net/http/#Error

}
