package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	phase int
)

func init() {
	// needed because usage was output-ed with note "flag provided but not defined: -phase"
	flag.Usage = func() {}

	flag.IntVar(&phase, "phase", 3, "Which program phaze to test")

	flag.Parse()

	fmt.Println(`Custom flags for 'go test':
  -phase int
    	Which program phase to test (default 3)
`)
}

func TestFlags(t *testing.T) {

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	_ = parseFlags()

	if !flag.Parsed() {
		t.Error("Expected cli flags to be parsed!")
	}

	if flag.Lookup("list") == nil {
		t.Errorf("Expected cli flag %q to be readed", "list")
	}

	if flag.Lookup("out") == nil {
		t.Errorf("Expected cli flag %q to be readed", "out")
	}
}

func TestList(t *testing.T) {
	expected := []string{
		"test_resources/first.css",
		"test_resources/second.css",
		"test_resources/third.css",
	}

	listPath := "test_resources" + string(os.PathSeparator) + "list.js"

	result, err := list(listPath)
	if err != nil {
		t.Fatalf("List loading failed with error %s", err)
	}

	if len(result) != len(expected) {
		t.Fatal("Paths are not readed from json file.")
	}

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Errorf("For path %d, expeced %s but recoved %s", i, expected[i], result[i])
		}
	}
}

func TestMarge(t *testing.T) {
	paths := []string{
		"test_resources" + string(os.PathSeparator) + "first.css",
		"test_resources" + string(os.PathSeparator) + "second.css",
		"test_resources" + string(os.PathSeparator) + "third.css",
	}

	tempFile, err := ioutil.TempFile("test_resources", "cssMergeTest")
	if err != nil {
		t.Fatal(err)
	}

	stat, err := tempFile.Stat()
	if err != nil {
		tempFile.Close()
		t.Fatal(err)
	}

	out := "test_resources" + string(os.PathSeparator) + stat.Name()

	tempFile.Close()

	defer func() {
		os.Remove(out)
	}()

	err = merge(paths, out)
	if err != nil {
		t.Errorf("Merge failed with error %q", err)
	}

	outContent, err := ioutil.ReadFile(out)
	if err != nil {
		t.Errorf("Reading of out file failed with error %q", err)
	}

	var expected []byte
	for _, path := range paths {
		cssContent, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("Test setup failed while processing path %q with error %q", path, err)
		}
		expected = append(expected, cssContent...)
	}

	if string(outContent) != string(expected) {
		t.Errorf("Expected %q recived %q", expected, outContent)
	}
}

func TestWatch(t *testing.T) {
	if phase < 2 {
		return
	}

	paths := []string{
		"test_resources" + string(os.PathSeparator) + "first.css",
	}

	originalContent, err := ioutil.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}

	defer func(filename string, content []byte) {
		err := ioutil.WriteFile(filename, content, 0666)
		if err != nil {
			t.Fatal(err)
		}
	}(paths[0], originalContent)

	tempFile, err := ioutil.TempFile("test_resources", "cssWatchTest")
	if err != nil {
		t.Fatal(err)
	}

	stat, err := tempFile.Stat()
	if err != nil {
		tempFile.Close()
		t.Fatal(err)
	}

	out := "test_resources" + string(os.PathSeparator) + stat.Name()

	tempFile.Close()

	defer func() {
		os.Remove(out)
	}()

	go func() {
		err := watch(paths, out)
		if err != nil {
			t.Errorf("Error while watching %q", err)
		}
	}()

	initStats, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	cssFile, err := os.OpenFile(paths[0], os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := cssFile.WriteString("/* --- */"); err != nil {
		cssFile.Close()
		t.Fatal(err)
	}

	cssFile.Close()

	time.Sleep(1 * time.Second)

	endStats, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}

	if endStats.Size() == initStats.Size() {
		t.Errorf("Watch process is not working %d %d", endStats.Size(), initStats.Size())
	}
}

func TestWatchFlag(t *testing.T) {
	if phase < 2 {
		return
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	_ = parseFlags()

	if !flag.Parsed() {
		t.Error("Expected cli flags to be parsed!")
	}

	if flag.Lookup("watch") == nil {
		t.Errorf("Expected cli flag %q to be readed", "watch")
	}
}

func TestServe(t *testing.T) {
	if phase < 3 {
		return
	}

	paths := []string{
		"test_resources" + string(os.PathSeparator) + "first.css",
		"test_resources" + string(os.PathSeparator) + "second.css",
		"test_resources" + string(os.PathSeparator) + "third.css",
	}

	tempFile, err := ioutil.TempFile("test_resources", "cssMergeTest")
	if err != nil {
		t.Fatal(err)
	}

	stat, err := tempFile.Stat()
	if err != nil {
		tempFile.Close()
		t.Fatal(err)
	}

	out := "test_resources" + string(os.PathSeparator) + stat.Name()

	tempFile.Close()

	defer func() {
		os.Remove(out)
	}()

	err = merge(paths, out)
	if err != nil {
		t.Errorf("Merge failed with error %q", err)
	}

	mergedFile := MergedFile(out)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	mergedFile.ServeHTTP(w, req)

	var expected []byte
	for _, path := range paths {
		cssContent, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("Test setup failed while processing path %q with error %q", path, err)
		}
		expected = append(expected, cssContent...)
	}

	if string(expected) != w.Body.String() {
		t.Errorf("Expected serve to output %q, but recived %q", expected, w.Body.String())
	}
}

func TestServeFlag(t *testing.T) {
	if phase < 3 {
		return
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	_ = parseFlags()

	if !flag.Parsed() {
		t.Error("Expected cli flags to be parsed!")
	}

	if flag.Lookup("serve") == nil {
		t.Errorf("Expected cli flag %q to be readed", "serve")
	}
}
