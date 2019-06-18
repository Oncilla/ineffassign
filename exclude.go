package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var _ json.Unmarshaler = (*Exclude)(nil)

// Exclude parses the exclude json configuration. Files and directories that are
// matched by any pattern in the json are excluded from the analysis.
//
// This is a sample json:
//  {
//    "testdata/testdata.go": "Tracked in issue #42"
//  }
//
// The keys of the config are the patterns to match. They can either be relative
// or absolute path patterns. The value serves as comment to indicate why the
// ignore is necessary. It is good practice to attach the github issue number
// that tracks the reason.
type Exclude map[string]string

// LoadFromFile loads the exclude configuration from the specified file.
func (e *Exclude) LoadFromFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, e)
}

// UnmarshalJSON parses the config, checks whether all pattern are
// well-formed, and transforms relative patterns to absolute patterns.
func (e *Exclude) UnmarshalJSON(b []byte) error {
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*e = make(map[string]string)
	for pattern := range m {
		if _, err := filepath.Glob(pattern); err != nil {
			return fmt.Errorf("invalid pattern: pattern=%s err=%s", pattern, err)
		}
		absPattern, err := filepath.Abs(pattern)
		if err != nil {
			return fmt.Errorf("unable to get absolute pattern: pattern=%s err=%err", pattern, err)
		}
		(*e)[absPattern] = m[pattern]
	}
	return nil
}

// Excluded indicates whether this path is excluded by the configuration.
func (e *Exclude) Excluded(path string) bool {
	for pattern := range *e {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			panic(fmt.Sprintf("Uncaught bad pattern"))
		}
		if matched {
			return true
		}
	}
	return false

}
