// Command check validates the Mioty blueprint registry: run it with the registry root as the only argument.
package main

import (
	"cmp"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	maxFileBytes      = 512 << 10
	maxSpecBytes      = 256 << 10
	maxModelSpecBytes = 2 << 20
	maxVersions       = 50
)

// Problem is one violation, printed as "<file>: <field>: <message>".
type Problem struct {
	File, Field, Message string
}

func (p Problem) String() string {
	return fmt.Sprintf("%s: %s: %s", printable(p.File), printable(p.Field), printable(p.Message))
}

// printable quotes text that could split the one-line output; workflow commands are disabled by run's frame.
func printable(s string) string {
	unsafe := !utf8.ValidString(s) ||
		strings.IndexFunc(s, func(r rune) bool { return !strconv.IsPrint(r) }) >= 0
	if unsafe {
		return strconv.Quote(s)
	}
	return s
}

// Check returns every violation in the registry tree, sorted by file, field and message.
func Check(fsys fs.FS) []Problem {
	files, problems := discover(fsys)
	models := map[string]*model{}
	for _, file := range files {
		dir := path.Dir(file)
		m, ok := models[dir]
		if !ok {
			m = &model{dir: dir}
			models[dir] = m
		}
		r, ps := readRecord(fsys, file)
		problems = append(problems, ps...)
		m.add(r)
	}
	problems = append(problems, checkModels(models)...)
	slices.SortFunc(problems, func(a, b Problem) int {
		return cmp.Or(cmp.Compare(a.File, b.File), cmp.Compare(a.Field, b.Field), cmp.Compare(a.Message, b.Message))
	})
	return problems
}
