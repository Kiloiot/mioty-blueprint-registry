package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/Kiloiot/kilo-service-center/KC-Core/pkg/blueprint"
)

const oversizeCase = "oversize-file"

func TestCheck_Cases(t *testing.T) {
	wantCases := []string{
		"bad-layout", "duplicate-version", "lowercase-spec-type-eui", "missing-envelope-field",
		"mixed-type-eui-in-model-dir", "oversize-file", "shared-type-eui", "stray-json", "symlink",
		"tools-manufacturer", "two-defaults", "unparsable-spec", "version-mismatch",
	}
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	var cases []string
	for _, e := range entries {
		if e.IsDir() {
			cases = append(cases, e.Name())
		}
	}
	if !slices.Equal(cases, wantCases) {
		t.Fatalf("testdata cases %v, want the shared list %v", cases, wantCases)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			dir := filepath.Join("testdata", e.Name())
			want, err := os.ReadFile(filepath.Join(dir, "expected.txt"))
			if err != nil {
				t.Fatal(err)
			}
			if e.Name() == oversizeCase {
				dir = withOversizeFile(t, dir)
			}
			got := format(Check(os.DirFS(dir)))
			if got != string(want) {
				t.Errorf("problems mismatch\n--- got\n%s--- want\n%s", got, want)
			}
		})
	}
}

func TestCheck_RepositoryRoot(t *testing.T) {
	// Only the count: the records are pull request data, and only the framed registry check step may print them.
	if problems := Check(os.DirFS(filepath.Join("..", ".."))); len(problems) != 0 {
		t.Errorf("repository root must pass, got %d problems; see the registry check step", len(problems))
	}
}

func TestCheck_Rules(t *testing.T) {
	tests := []struct {
		name string
		fsys fstest.MapFS
		want string
	}{
		{
			name: "version with a leading v is named as KiloCenter submits it",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "v1.0", true, testSpec("70B3D59CD0000101", "v1.0")))},
			want: "",
		},
		{
			name: "version KiloCenter cannot name",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0/x", true, testSpec("70B3D59CD0000101", "1.0/x")))},
			want: `acme/widget/v1.0.json: version: "1.0/x" is not a valid registry version: ` + blueprint.ErrRegistryPathSegment.Error() + "\n",
		},
		{
			name: "spec typeEui differs from envelope",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000999", "1.0")))},
			want: `acme/widget/v1.0.json: spec.typeEui: "70B3D59CD0000999" does not match type_eui "70B3D59CD0000101"` + "\n",
		},
		{
			name: "spec typeEui overridden by a differently cased key the decoder reads last",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true,
				withSpecKey(testSpec("70B3D59CD0000101", "1.0"), `"TypeEUI": "70B3D59CD0000999"`)))},
			want: `acme/widget/v1.0.json: spec.typeEui: "70B3D59CD0000999" does not match type_eui "70B3D59CD0000101"` + "\n",
		},
		{
			name: "spec version overridden by a differently cased key the decoder reads last",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true,
				withSpecKey(testSpec("70B3D59CD0000101", "1.0"), `"VERSION": "2.0"`)))},
			want: `acme/widget/v1.0.json: spec.version: "2.0" does not equal version "1.0"` + "\n",
		},
		{
			name: "envelope key differing from a known key only in case",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.Replace(
				envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")),
				`"type_eui"`, `"TYPE_EUI": "70B3D59CD0000999", "type_eui"`, 1))},
			want: `acme/widget/v1.0.json: TYPE_EUI: differs from "type_eui" only in case` + "\n",
		},
		{
			name: "envelope Type EUI not uppercase hex and ids not UUIDs",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.NewReplacer(
				`"id": "00000000-0000-4000-8000-000000000001"`, `"id": "abc"`,
				`"device_model_id": "00000000-0000-4000-8000-000000000002"`, `"device_model_id": 7`,
			).Replace(envelope("70b3d59cd000010g", "1.0", true, testSpec("70B3D59CD0000101", "1.0"))))},
			want: "acme/widget/v1.0.json: device_model_id: must be a string\n" +
				"acme/widget/v1.0.json: id: must be a UUID\n" +
				`acme/widget/v1.0.json: spec.typeEui: "70B3D59CD0000101" does not match type_eui "70b3d59cd000010g"` + "\n" +
				"acme/widget/v1.0.json: type_eui: must be 16 uppercase hex characters\n",
		},
		{
			name: "envelope is not JSON",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file("{")},
			want: "acme/widget/v1.0.json: file: invalid JSON: unexpected end of JSON input\n",
		},
		{
			name: "spec is not an object",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true, `"text"`))},
			want: "acme/widget/v1.0.json: spec: must be a JSON object\n",
		},
		{
			name: "spec over 256 KB",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true, paddedSpec("70B3D59CD0000101", "1.0", maxSpecBytes+1)))},
			want: "acme/widget/v1.0.json: spec: larger than 256 KB\n",
		},
		{
			name: "too many versions",
			fsys: manyVersions(maxVersions + 1),
			want: fmt.Sprintf("acme/widget: version: %d versions, at most %d\n", maxVersions+1, maxVersions),
		},
		{
			name: "specs of a model over 2 MB",
			fsys: bigModel(9),
			want: "acme/widget: spec: specs of the model add up to more than 2 MB\n",
		},
		{
			name: "layout is reported for an unparsable file",
			fsys: fstest.MapFS{"Acme/widget/v1.0.json": file("{")},
			want: `Acme/widget/v1.0.json: file: invalid JSON: unexpected end of JSON input` + "\n" +
				`Acme/widget/v1.0.json: path: directory "Acme" must be named "acme"` + "\n",
		},
		{
			name: "is_default is not a boolean",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.Replace(
				envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")),
				`"is_default": true`, `"is_default": "true"`, 1))},
			want: "acme/widget: is_default: 0 default versions, want exactly one\n" +
				"acme/widget/v1.0.json: is_default: must be a boolean\n",
		},
		{
			name: "spec is missing",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.Replace(
				envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")),
				`"spec"`, `"payload"`, 1))},
			want: "acme/widget/v1.0.json: spec: required\n",
		},
		{
			name: "non-regular file that is not json",
			fsys: fstest.MapFS{"acme/README.md": &fstest.MapFile{Mode: fs.ModeSymlink}},
			want: "acme/README.md: path: not a regular file\n",
		},
		{
			name: "file name that is not UTF-8 cannot be opened through fs.FS, is refused and is quoted",
			fsys: fstest.MapFS{"acme/widget/v\xff.json": file(envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")))},
			want: `"acme/widget/v\xff.json": file: "cannot open: open acme/widget/v\xff.json: file does not exist"` + "\n",
		},
		{
			name: "file names cannot split output lines",
			fsys: fstest.MapFS{
				"a\nb.json":       file("{}"),
				"::error::x.json": file("{}"),
			},
			want: `::error::x.json: path: .json outside the {manufacturer}/{model}/v{version}.json layout` + "\n" +
				`"a\nb.json": path: .json outside the {manufacturer}/{model}/v{version}.json layout` + "\n",
		},
		{
			name: "envelope is null",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file("null")},
			want: "acme/widget/v1.0.json: file: must be a JSON object\n",
		},
		{
			name: "timestamps must be RFC3339 strings when present",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.Replace(
				envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")),
				`"is_default"`, `"created_at": "2026-04-30 15:43", "updated_at": 7, "is_default"`, 1))},
			want: "acme/widget/v1.0.json: created_at: must be an RFC3339 timestamp\n" +
				"acme/widget/v1.0.json: updated_at: must be a string\n",
		},
		{
			name: "RFC3339 timestamps with fractional seconds pass",
			fsys: fstest.MapFS{"acme/widget/v1.0.json": file(strings.Replace(
				envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0")),
				`"is_default"`, `"created_at": "2026-04-30T15:43:01.680226Z", "updated_at": "2026-04-30T15:43:01+02:00", "is_default"`, 1))},
			want: "",
		},
		{
			name: "directories KiloCenter would not build; a model directory may be a model code with a double hyphen",
			fsys: fstest.MapFS{
				"-acme/widget/v1.0.json": file(envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0"))),
				"acme/a--b/v1.0.json":    file(envelope("70B3D59CD0000102", "1.0", true, testSpec("70B3D59CD0000102", "1.0"))),
				"a--b/widget/v1.0.json":  file(envelope("70B3D59CD0000104", "1.0", true, testSpec("70B3D59CD0000104", "1.0"))),
				"acme/-x/v1.0.json":      file(envelope("70B3D59CD0000105", "1.0", true, testSpec("70B3D59CD0000105", "1.0"))),
				"beta/---/v1.0.json":     file(envelope("70B3D59CD0000103", "1.0", true, testSpec("70B3D59CD0000103", "1.0"))),
			},
			want: `-acme/widget/v1.0.json: path: directory "-acme" must be named "acme"` + "\n" +
				`a--b/widget/v1.0.json: path: directory "a--b" must be named "a-b"` + "\n" +
				`acme/-x/v1.0.json: path: directory "-x" must be named "x"` + "\n" +
				`beta/---/v1.0.json: path: directory "---" is not a valid registry path segment: ` + blueprint.ErrRegistryPathSegment.Error() + "\n",
		},
		{
			name: "dot directories and non-json files are not records",
			fsys: fstest.MapFS{
				".github/check/testdata/x/v1.json": file("{"),
				"README.md":                        file("readme"),
				"acme/README.md":                   file("readme"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(Check(tt.fsys)); got != tt.want {
				t.Errorf("problems mismatch\n--- got\n%s--- want\n%s", got, tt.want)
			}
		})
	}
}

func format(problems []Problem) string {
	var b strings.Builder
	for _, p := range problems {
		b.WriteString(p.String())
		b.WriteByte('\n')
	}
	return b.String()
}

func withOversizeFile(t *testing.T, caseDir string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(caseDir)); err != nil {
		t.Fatal(err)
	}
	big := filepath.Join(dir, "acme", "big")
	if err := os.MkdirAll(big, 0o755); err != nil {
		t.Fatal(err)
	}
	content := envelope("70B3D59CD0000101", "1.0", true, testSpec("70B3D59CD0000101", "1.0"))
	content += strings.Repeat(" ", maxFileBytes+1-len(content))
	if err := os.WriteFile(filepath.Join(big, "v1.0.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func manyVersions(n int) fstest.MapFS {
	fsys := fstest.MapFS{}
	for i := range n {
		version := fmt.Sprintf("1.%d", i)
		fsys[path.Join("acme/widget", "v"+version+".json")] = file(envelope("70B3D59CD0000101", version, i == 0, testSpec("70B3D59CD0000101", version)))
	}
	return fsys
}

func bigModel(n int) fstest.MapFS {
	fsys := fstest.MapFS{}
	for i := range n {
		version := fmt.Sprintf("1.%d", i)
		spec := paddedSpec("70B3D59CD0000101", version, maxSpecBytes-1)
		fsys[path.Join("acme/widget", "v"+version+".json")] = file(envelope("70B3D59CD0000101", version, i == 0, spec))
	}
	return fsys
}

func file(content string) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(content)}
}

func envelope(typeEUI, version string, isDefault bool, spec string) string {
	return fmt.Sprintf(`{
  "id": "00000000-0000-4000-8000-000000000001",
  "device_model_id": "00000000-0000-4000-8000-000000000002",
  "manufacturer_name": "Acme",
  "device_model_name": "Widget",
  "device_model_code": "widget",
  "type_eui": %q,
  "version": %q,
  "is_default": %t,
  "spec": %s
}`, typeEUI, version, isDefault, spec)
}

func testSpec(typeEUI, version string) string {
	return fmt.Sprintf(`{"meta": {"name": "Widget", "vendor": "Acme"}, "uplink": [{"id": 0, "crypto": 0, "payload": [{"name": "t", "component": "c"}]}], "typeEui": %q, "version": %q, "component": {"c": {"size": 16, "type": "int", "littleEndian": false}}}`, typeEUI, version)
}

func paddedSpec(typeEUI, version string, size int) string {
	s := testSpec(typeEUI, version)
	return s[:len(s)-1] + strings.Repeat(" ", size-len(s)) + "}"
}

func withSpecKey(spec, member string) string {
	return spec[:len(spec)-1] + ", " + member + "}"
}
