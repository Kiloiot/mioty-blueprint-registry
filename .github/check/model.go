package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// model is one {manufacturer}/{model} directory with all its versions.
type model struct {
	dir        string
	files      int
	records    []*record
	incomplete bool
}

func (m *model) add(r *record) {
	m.files++
	if r == nil {
		m.incomplete = true
		return
	}
	m.records = append(m.records, r)
}

func checkModels(models map[string]*model) []Problem {
	var problems []Problem
	owners := map[string][]string{}
	for _, dir := range slices.Sorted(maps.Keys(models)) {
		m := models[dir]
		problems = append(problems, m.check()...)
		for _, eui := range m.typeEUIs() {
			owners[eui] = append(owners[eui], dir)
		}
	}
	return append(problems, sharedTypeEUIs(owners)...)
}

func (m *model) check() []Problem {
	rep := &reporter{file: m.dir}
	if m.files > maxVersions {
		rep.add("version", "%d versions, at most %d", m.files, maxVersions)
	}
	total := 0
	for _, r := range m.records {
		total += r.specBytes
	}
	if total > maxModelSpecBytes {
		rep.add("spec", "specs of the model add up to more than %d MB", maxModelSpecBytes>>20)
	}
	if euis := m.typeEUIs(); len(euis) > 1 {
		rep.add("type_eui", "versions have different Type EUIs: %s", strings.Join(euis, ", "))
	}
	m.checkVersions(rep)
	m.checkDefault(rep)
	return rep.problems
}

func (m *model) typeEUIs() []string {
	set := map[string]bool{}
	for _, r := range m.records {
		if r.typeEUI != "" {
			set[strings.ToUpper(r.typeEUI)] = true
		}
	}
	return slices.Sorted(maps.Keys(set))
}

func (m *model) checkVersions(rep *reporter) {
	counts := map[string]int{}
	for _, r := range m.records {
		if r.version != "" {
			counts[r.version]++
		}
	}
	for _, v := range slices.Sorted(maps.Keys(counts)) {
		if counts[v] > 1 {
			rep.add("version", "%q is used by %d files", v, counts[v])
		}
	}
}

// checkDefault skips a model with an unreadable file: its default flag is unknown.
func (m *model) checkDefault(rep *reporter) {
	if m.incomplete {
		return
	}
	defaults := 0
	for _, r := range m.records {
		if r.isDefault {
			defaults++
		}
	}
	if defaults != 1 {
		rep.add("is_default", "%d default versions, want exactly one", defaults)
	}
}

func sharedTypeEUIs(owners map[string][]string) []Problem {
	var problems []Problem
	for _, eui := range slices.Sorted(maps.Keys(owners)) {
		dirs := owners[eui]
		if len(dirs) < 2 {
			continue
		}
		for i, dir := range dirs {
			others := slices.Concat(dirs[:i], dirs[i+1:])
			problems = append(problems, Problem{dir, "type_eui", fmt.Sprintf("%s is also used by %s", eui, strings.Join(others, ", "))})
		}
	}
	return problems
}
