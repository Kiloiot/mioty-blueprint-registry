package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"time"

	"github.com/Kiloiot/kilo-service-center/KC-Core/pkg/blueprint"
)

type record struct {
	typeEUI   string
	version   string
	isDefault bool
	specBytes int
}

type reporter struct {
	file     string
	problems []Problem
}

func (r *reporter) add(field, format string, args ...any) {
	r.problems = append(r.problems, Problem{r.file, field, fmt.Sprintf(format, args...)})
}

// readRecord returns a nil record when the file is unreadable or not a JSON object.
func readRecord(fsys fs.FS, file string) (*record, []Problem) {
	rep := &reporter{file: file}
	checkPath(rep, file)
	data, err := readLimited(fsys, file)
	if err != nil {
		rep.add("file", "%s", err)
		return nil, rep.problems
	}
	var env map[string]json.RawMessage
	if err := json.Unmarshal(data, &env); err != nil {
		rep.add("file", "invalid JSON: %s", err)
		return nil, rep.problems
	}
	if env == nil {
		rep.add("file", "must be a JSON object")
		return nil, rep.problems
	}
	r := &record{}
	checkKeyCase(rep, env)
	checkEnvelope(rep, env, r)
	checkFileName(rep, file, r.version)
	checkSpec(rep, env["spec"], r)
	return r, rep.problems
}

func readLimited(fsys fs.FS, file string) ([]byte, error) {
	f, err := fsys.Open(file)
	if err != nil {
		return nil, fmt.Errorf("cannot open: %w", err)
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}
	if len(data) > maxFileBytes {
		return nil, fmt.Errorf("larger than %d KB", maxFileBytes>>10)
	}
	return data, nil
}

// checkPath requires directory names as KiloCenter builds them on submit: the model one from its code or its name.
func checkPath(rep *reporter, file string) {
	manufacturer, model := path.Split(path.Dir(file))
	manufacturer = strings.TrimSuffix(manufacturer, "/")
	checkSegment(rep, manufacturer)
	if code, err := blueprint.CanonicalModelCode(model); err == nil && code == model {
		return
	}
	checkSegment(rep, model)
}

func checkSegment(rep *reporter, dir string) {
	segment, err := blueprint.PathSegment(dir)
	switch {
	case err != nil:
		rep.add("path", "directory %q is not a valid registry path segment: %v", dir, err)
	case segment != dir:
		rep.add("path", "directory %q must be named %q", dir, segment)
	}
}

func checkFileName(rep *reporter, file, version string) {
	if version == "" {
		return
	}
	// Same derivation as the KiloCenter "Submit to Registry" path, so a submitted file always passes.
	segment, err := blueprint.VersionSegment(blueprint.NormalizeVersion(version))
	if err != nil {
		rep.add("version", "%q is not a valid registry version: %v", version, err)
		return
	}
	if want := segment + ".json"; path.Base(file) != want {
		rep.add("path", "file name must be %q for version %q", want, version)
	}
}

// checkKeyCase refuses near-duplicates of envelope keys, which a case-insensitive JSON reader would take instead.
func checkKeyCase(rep *reporter, env map[string]json.RawMessage) {
	envelopeKeys := []string{
		"id", "device_model_id", "manufacturer_name", "device_model_name", "device_model_code",
		"type_eui", "version", "is_default", "spec", "created_at", "updated_at",
	}
	for key := range env {
		for _, known := range envelopeKeys {
			if key != known && strings.EqualFold(key, known) {
				rep.add(key, "differs from %q only in case", known)
			}
		}
	}
}

func checkEnvelope(rep *reporter, env map[string]json.RawMessage, r *record) {
	for _, key := range []string{"id", "device_model_id"} {
		if s, ok := requiredString(rep, env, key); ok && !isUUID(s) {
			rep.add(key, "must be a UUID")
		}
	}
	for _, key := range []string{"manufacturer_name", "device_model_name", "device_model_code"} {
		requiredString(rep, env, key)
	}
	if s, ok := requiredString(rep, env, "type_eui"); ok {
		r.typeEUI = s
		if !isTypeEUI(s) {
			rep.add("type_eui", "must be 16 uppercase hex characters")
		}
	}
	if s, ok := requiredString(rep, env, "version"); ok {
		r.version = s
	}
	r.isDefault = requiredBool(rep, env, "is_default")
	for _, key := range []string{"created_at", "updated_at"} {
		optionalTimestamp(rep, env, key)
	}
}

func optionalTimestamp(rep *reporter, env map[string]json.RawMessage, key string) {
	raw, ok := env[key]
	if !ok {
		return
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		rep.add(key, "must be a string")
		return
	}
	if _, err := time.Parse(time.RFC3339, s); err != nil {
		rep.add(key, "must be an RFC3339 timestamp")
	}
}

func requiredString(rep *reporter, env map[string]json.RawMessage, key string) (string, bool) {
	raw, ok := env[key]
	if !ok {
		rep.add(key, "required")
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		rep.add(key, "must be a string")
		return "", false
	}
	if s == "" {
		rep.add(key, "required")
		return "", false
	}
	return s, true
}

func requiredBool(rep *reporter, env map[string]json.RawMessage, key string) bool {
	raw, ok := env[key]
	if !ok {
		rep.add(key, "required")
		return false
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err != nil {
		rep.add(key, "must be a boolean")
		return false
	}
	return b
}

func checkSpec(rep *reporter, raw json.RawMessage, r *record) {
	if raw == nil {
		rep.add("spec", "required")
		return
	}
	if trimmed := bytes.TrimLeft(raw, " \t\r\n"); len(trimmed) == 0 || trimmed[0] != '{' {
		rep.add("spec", "must be a JSON object")
		return
	}
	r.specBytes = len(raw)
	if len(raw) > maxSpecBytes {
		rep.add("spec", "larger than %d KB", maxSpecBytes>>10)
		return
	}
	// Compare what the decoder reads: it matches keys case-insensitively and the last one wins.
	spec, err := blueprint.ParseAndValidate(raw)
	if err != nil {
		rep.add("spec", "rejected by the KiloCenter decoder: %s", err)
		return
	}
	if r.typeEUI != "" && !strings.EqualFold(spec.TypeEUI, r.typeEUI) {
		rep.add("spec.typeEui", "%q does not match type_eui %q", spec.TypeEUI, r.typeEUI)
	}
	if r.version != "" && spec.Version != r.version {
		rep.add("spec.version", "%q does not equal version %q", spec.Version, r.version)
	}
}

func isTypeEUI(s string) bool {
	if len(s) != 16 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !isHex(c) {
				return false
			}
		}
	}
	return true
}

func isHex(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
