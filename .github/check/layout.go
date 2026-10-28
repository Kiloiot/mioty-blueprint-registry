package main

import (
	"io/fs"
	"path"
	"strings"
)

const recordDepth = 3

// discover applies the record-file rule shared with the RIS registry import, so both see the same records.
func discover(fsys fs.FS) ([]string, []Problem) {
	var w walker
	w.walk(fsys, ".", 1)
	return w.records, w.problems
}

type walker struct {
	records  []string
	problems []Problem
}

func (w *walker) walk(fsys fs.FS, dir string, depth int) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		w.problems = append(w.problems, Problem{dir, "path", "cannot read directory: " + err.Error()})
		return
	}
	for _, e := range entries {
		name := path.Join(dir, e.Name())
		switch {
		case e.IsDir():
			if !strings.HasPrefix(e.Name(), ".") {
				w.walk(fsys, name, depth+1)
			}
		case !e.Type().IsRegular():
			w.problems = append(w.problems, Problem{name, "path", "not a regular file"})
		case path.Ext(e.Name()) != ".json":
		case depth == recordDepth:
			w.records = append(w.records, name)
		default:
			w.problems = append(w.problems, Problem{name, "path", ".json outside the {manufacturer}/{model}/v{version}.json layout"})
		}
	}
}
