// Package changes assembles release notes from the individual change files
// contributors add to a project's changes directory, so that pull requests
// never edit the same lines of the same changelog.
package changes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

// Paths used to store the change files, relative to the repository root.
const (
	Dir           = "changes"
	UnreleasedDir = Dir + "/unreleased"
	ReleasesDir   = Dir + "/releases"
	HeaderFile    = Dir + "/HEADER.md"
	ChangelogFile = "CHANGELOG.md"
)

// DateFormat is used in release note titles and file names.
const DateFormat = "2006-01-02"

// Heading levels: change files and release notes stand alone, so their
// sections start at the second level. The changelog nests each release under
// its own heading, so its sections start at the third.
const (
	notesSection     = "## "
	changelogSection = "### "
)

// defaultHeader introduces the generated changelog of a project that does not
// provide its own HEADER.md.
const defaultHeader = `# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

This file is generated from the change files in the [changes](./changes) directory. Add new entries to ` + "`changes/unreleased`" + ` instead of editing this file; see [changes/README.md](./changes/README.md).`

// header returns the introduction to the generated changelog, taken from the
// project's HEADER.md when it provides one.
func header(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, HeaderFile)) //nolint:gosec // paths come from the repository
	if err != nil {
		if os.IsNotExist(err) {
			return defaultHeader + "\n", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)) + "\n", nil
}

// Sections are the headings a change file may use, in the order they are
// rendered.
var Sections = []string{
	"Added",
	"Changed",
	"Deprecated",
	"Removed",
	"Fixed",
	"Security",
}

// Fragment contains the notes read from a change file or release note.
type Fragment struct {
	// Name of the file the fragment was read from, without the directory.
	Name string
	// Title holds the first level heading, used by release notes.
	Title string
	// Preamble holds any text appearing before the first section.
	Preamble []string
	// Sections holds the entries grouped by heading, in the order found.
	Sections []*Section
}

// Section groups the entries under a single heading.
type Section struct {
	Title string
	Lines []string
}

// Parse reads the title and sections out of a change file or release note.
func Parse(name string, data []byte) *Fragment {
	f := &Fragment{Name: name}
	var cur *Section
	var fenced, blank bool
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
		}
		if !fenced {
			if title, ok := strings.CutPrefix(line, "# "); ok && f.Title == "" && cur == nil {
				f.Title = strings.TrimSpace(title)
				blank = true
				continue
			}
			if title, ok := sectionTitle(line); ok {
				cur = f.section(title)
				blank = true
				continue
			}
		}
		if blank {
			if strings.TrimSpace(line) == "" {
				continue
			}
			blank = false
		}
		if cur == nil {
			f.Preamble = append(f.Preamble, line)
			continue
		}
		cur.Lines = append(cur.Lines, line)
	}
	f.Preamble = trimBlank(f.Preamble)
	for _, s := range f.Sections {
		s.Lines = trimBlank(s.Lines)
	}
	return f
}

// Validate ensures a change file uses known section headings and provides at
// least one entry.
func (f *Fragment) Validate() error {
	if f.Title != "" {
		return fmt.Errorf("%s: unexpected title %q, start with a section heading like '## Added'", f.Name, f.Title)
	}
	empty := true
	for _, s := range f.Sections {
		if !knownSection(s.Title) {
			return fmt.Errorf("%s: unknown heading %q, expected one of: %s", f.Name, s.Title, strings.Join(Sections, ", "))
		}
		if len(s.Lines) > 0 {
			empty = false
		}
	}
	if empty {
		return fmt.Errorf("%s: no entries found, expected bullet points under a heading like '## Added'", f.Name)
	}
	return nil
}

func (f *Fragment) section(title string) *Section {
	for _, s := range f.Sections {
		if strings.EqualFold(s.Title, title) {
			if lines := trimBlank(s.Lines); len(lines) > 0 {
				s.Lines = append(lines, "")
			}
			return s
		}
	}
	s := &Section{Title: title}
	f.Sections = append(f.Sections, s)
	return s
}

// LoadUnreleased reads and validates the pending change files, sorted by name.
func LoadUnreleased(root string) ([]*Fragment, error) {
	names, err := filepath.Glob(filepath.Join(root, UnreleasedDir, "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	out := make([]*Fragment, 0, len(names))
	for _, name := range names {
		base := filepath.Base(name)
		if strings.EqualFold(base, "README.md") {
			continue
		}
		data, err := os.ReadFile(name) //nolint:gosec // paths come from the repository
		if err != nil {
			return nil, err
		}
		f := Parse(base, data)
		if err := f.Validate(); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

// Notes merges the change files into the release note for a version.
func Notes(version string, date time.Time, frags []*Fragment) string {
	merged := &Fragment{Sections: mergeSections(frags)}
	for _, f := range frags {
		merged.Preamble = append(merged.Preamble, f.Preamble...)
	}
	title := fmt.Sprintf("# %s - %s", version, date.Format(DateFormat))
	return render(title, notesSection, merged)
}

// mergeSections gathers the entries of every fragment under the known
// headings, in the order the headings are declared.
func mergeSections(frags []*Fragment) []*Section {
	out := make([]*Section, 0, len(Sections))
	for _, title := range Sections {
		s := &Section{Title: title}
		for _, f := range frags {
			for _, fs := range f.Sections {
				if strings.EqualFold(fs.Title, title) {
					s.Lines = append(s.Lines, fs.Lines...)
				}
			}
		}
		if len(s.Lines) > 0 {
			out = append(out, s)
		}
	}
	return out
}

// render writes the notes under the given title, with each section heading
// introduced by the given prefix.
func render(title, prefix string, f *Fragment) string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "%s\n", title)
	if len(f.Preamble) > 0 {
		fmt.Fprintf(b, "\n%s\n", strings.Join(f.Preamble, "\n"))
	}
	for _, s := range f.Sections {
		if len(s.Lines) == 0 {
			continue
		}
		fmt.Fprintf(b, "\n%s%s\n\n%s\n", prefix, s.Title, strings.Join(s.Lines, "\n"))
	}
	return b.String()
}

// Preview returns the release note the pending change files would produce.
func Preview(root, version string, date time.Time) (string, error) {
	frags, err := LoadUnreleased(root)
	if err != nil {
		return "", err
	}
	if len(frags) == 0 {
		return "", fmt.Errorf("no change files found in %s", UnreleasedDir)
	}
	return Notes(version, date, frags), nil
}

// Release moves the pending change files into a dated release note and
// regenerates the changelog. It returns the path of the release note.
func Release(root, version string, date time.Time) (string, error) {
	notes, err := Preview(root, version, date)
	if err != nil {
		return "", err
	}
	if prev, err := findRelease(root, version); err != nil {
		return "", err
	} else if prev != "" {
		return "", fmt.Errorf("release notes for %s already exist at %s, bump the version in version.go first", version, prev)
	}
	name := filepath.Join(ReleasesDir, fmt.Sprintf("%s-%s.md", date.Format(DateFormat), version))
	if err := os.WriteFile(filepath.Join(root, name), []byte(notes), 0644); err != nil { //nolint:gosec // release notes are public
		return "", err
	}
	if err := removeUnreleased(root); err != nil {
		return "", err
	}
	log, err := Changelog(root)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(root, ChangelogFile), []byte(log), 0644); err != nil { //nolint:gosec // the changelog is public
		return "", err
	}
	return name, nil
}

// Normalize rewrites any release note that differs from its canonical form,
// returning the names of the files changed.
func Normalize(root string) ([]string, error) {
	names, err := filepath.Glob(filepath.Join(root, ReleasesDir, "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	changed := make([]string, 0)
	for _, path := range names {
		r, err := loadRelease(path)
		if err != nil {
			return nil, err
		}
		notes := r.canonical()
		if notes == r.raw {
			continue
		}
		if err := os.WriteFile(path, []byte(notes), 0644); err != nil { //nolint:gosec // release notes are public
			return nil, err
		}
		changed = append(changed, filepath.Base(path))
	}
	return changed, nil
}

// Changelog assembles the complete changelog from the release notes, newest
// release first.
func Changelog(root string) (string, error) {
	names, err := filepath.Glob(filepath.Join(root, ReleasesDir, "*.md"))
	if err != nil {
		return "", err
	}
	releases := make([]*release, 0, len(names))
	for _, name := range names {
		r, err := loadRelease(name)
		if err != nil {
			return "", err
		}
		releases = append(releases, r)
	}
	sort.Slice(releases, func(i, j int) bool {
		a, b := releases[i], releases[j]
		if a.version != nil && b.version != nil && !a.version.Equal(b.version) {
			return a.version.GreaterThan(b.version)
		}
		return a.name > b.name
	})
	head, err := header(root)
	if err != nil {
		return "", err
	}
	b := new(strings.Builder)
	b.WriteString(head)
	for _, r := range releases {
		title := fmt.Sprintf("## [%s] - %s", r.tag, r.date)
		fmt.Fprintf(b, "\n%s", render(title, changelogSection, r.notes))
	}
	return b.String(), nil
}

type release struct {
	name    string
	tag     string
	date    string
	raw     string
	version *semver.Version
	notes   *Fragment
}

// canonical renders the release note as it should be stored on disk.
func (r *release) canonical() string {
	return render(fmt.Sprintf("# %s - %s", r.tag, r.date), notesSection, r.notes)
}

// loadRelease reads a release note, whose file name provides the version and
// date of the release.
func loadRelease(path string) (*release, error) {
	name := filepath.Base(path)
	date, tag, ok := splitName(name)
	if !ok {
		return nil, fmt.Errorf("%s: expected a name like 2026-09-14-v0.506.0.md", name)
	}
	data, err := os.ReadFile(path) //nolint:gosec // paths come from the repository
	if err != nil {
		return nil, err
	}
	r := &release{name: name, tag: tag, date: date, raw: string(data), notes: Parse(name, data)}
	if title := fmt.Sprintf("%s - %s", tag, date); r.notes.Title != title {
		return nil, fmt.Errorf("%s: title is %q, expected %q", name, r.notes.Title, title)
	}
	if v, err := semver.NewVersion(tag); err == nil {
		r.version = v
	}
	return r, nil
}

// splitName extracts the date and version from a release note file name.
func splitName(name string) (date, tag string, ok bool) {
	date, tag, ok = strings.Cut(strings.TrimSuffix(name, ".md"), "-v")
	if !ok {
		return "", "", false
	}
	if _, err := time.Parse(DateFormat, date); err != nil {
		return "", "", false
	}
	return date, "v" + tag, true
}

// findRelease returns the release note file already published for the version,
// or an empty string.
func findRelease(root, version string) (string, error) {
	names, err := filepath.Glob(filepath.Join(root, ReleasesDir, "*-"+version+".md"))
	if err != nil || len(names) == 0 {
		return "", err
	}
	return filepath.Join(ReleasesDir, filepath.Base(names[0])), nil
}

func removeUnreleased(root string) error {
	names, err := filepath.Glob(filepath.Join(root, UnreleasedDir, "*.md"))
	if err != nil {
		return err
	}
	for _, name := range names {
		if strings.EqualFold(filepath.Base(name), "README.md") {
			continue
		}
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return nil
}

// sectionTitle returns the heading title of a section, accepting both the
// second and third heading levels.
func sectionTitle(line string) (string, bool) {
	for _, prefix := range []string{"### ", "## "} {
		if title, ok := strings.CutPrefix(line, prefix); ok {
			return strings.TrimSpace(title), true
		}
	}
	return "", false
}

func knownSection(title string) bool {
	for _, s := range Sections {
		if strings.EqualFold(s, title) {
			return true
		}
	}
	return false
}

func trimBlank(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
