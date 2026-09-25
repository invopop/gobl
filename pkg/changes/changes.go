// Package changes assembles release notes from the individual change files
// contributors add to a project's changes directory, so that pull requests
// never edit the same lines of the same changelog.
package changes

import (
	"fmt"
	"io/fs"
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
	path, err := resolve(root, HeaderFile)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path) //nolint:gosec // paths come from the repository
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
	// Strays holds any heading that is neither the title nor a section.
	Strays []string
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
	var fence string
	var blank bool
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if m := fenceMarker(line); m != "" {
			if fence == "" {
				fence = m
			} else if m[0] == fence[0] && len(m) >= len(fence) {
				fence = ""
			}
		}
		if fence == "" {
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
			if isHeading(line) {
				f.Strays = append(f.Strays, strings.TrimSpace(line))
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
	if len(f.Strays) > 0 {
		return fmt.Errorf("%s: unexpected heading %q, group entries under a section like '## Added'", f.Name, f.Strays[0])
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

// LoadUnreleased reads and validates the pending change files, sorted by path.
func LoadUnreleased(root string) ([]*Fragment, error) {
	paths, err := unreleasedFiles(root)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, UnreleasedDir)
	out := make([]*Fragment, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path) //nolint:gosec // paths come from the repository
		if err != nil {
			return nil, err
		}
		name := path
		if rel, err := filepath.Rel(dir, path); err == nil {
			name = filepath.ToSlash(rel)
		}
		f := Parse(name, data)
		if err := f.Validate(); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

// unreleasedFiles returns the pending change files, sorted by path. Branch
// names containing a slash nest the files in subdirectories, so the whole
// directory is walked.
func unreleasedFiles(root string) ([]string, error) {
	dir, err := resolve(root, UnreleasedDir)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0)
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}
		if strings.EqualFold(d.Name(), "README.md") {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return fmt.Errorf("%s: symbolic links are not supported", d.Name())
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

// Notes merges the change files into the release note for a version.
func Notes(version string, date time.Time, frags []*Fragment) string {
	merged := &Fragment{Sections: mergeSections(frags)}
	for _, f := range frags {
		if len(f.Preamble) == 0 {
			continue
		}
		if len(merged.Preamble) > 0 {
			merged.Preamble = append(merged.Preamble, "")
		}
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
	stamp := date.Format(DateFormat)
	base := fmt.Sprintf("%s-%s.md", stamp, version)
	if strings.ContainsAny(version, `/\`) {
		return "", fmt.Errorf("version %q cannot contain a path separator", version)
	}
	if d, tag, ok := splitName(base); !ok || tag != version || d != stamp {
		return "", fmt.Errorf("version %q cannot be read back from a release note name, expected something like v1.2.3", version)
	}
	notes, err := Preview(root, version, date)
	if err != nil {
		return "", err
	}
	if prev, err := findRelease(root, version); err != nil {
		return "", err
	} else if prev != "" {
		return "", fmt.Errorf("release notes for %s already exist at %s, bump the version in version.go first", version, prev)
	}
	name := filepath.Join(ReleasesDir, base)
	path, err := resolve(root, ReleasesDir+"/"+base)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:gosec // release notes are public
		return "", err
	}
	if err := writeFile(path, notes); err != nil {
		return "", err
	}
	// The change files are removed only once everything that can fail has
	// succeeded, so that a failed release leaves the repository as it was.
	log, err := Changelog(root)
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	logPath, err := resolve(root, ChangelogFile)
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	if err := writeFile(logPath, log); err != nil {
		_ = os.Remove(path)
		return "", err
	}
	if err := removeUnreleased(root); err != nil {
		return "", err
	}
	return name, nil
}

// Normalize rewrites any release note that differs from its canonical form,
// returning the names of the files changed.
func Normalize(root string) ([]string, error) {
	paths, err := releaseFiles(root)
	if err != nil {
		return nil, err
	}
	changed := make([]string, 0)
	for _, path := range paths {
		r, err := loadRelease(path)
		if err != nil {
			return nil, err
		}
		notes := r.canonical()
		if notes == r.raw {
			continue
		}
		if err := writeFile(path, notes); err != nil {
			return nil, err
		}
		changed = append(changed, filepath.Base(path))
	}
	return changed, nil
}

// resolve returns the full path of a location inside the project, rejecting a
// symbolic link at every step below the root so that nothing outside the
// project is read or written.
func resolve(root, rel string) (string, error) {
	path := root
	for _, elem := range strings.Split(rel, "/") {
		path = filepath.Join(path, elem)
		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("%s: symbolic links are not supported", elem)
		}
	}
	return path, nil
}

// writeFile replaces the file at the path with the data, leaving the original
// in place if the write fails.
func writeFile(path, data string) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".changes-*")
	if err != nil {
		return err
	}
	name := f.Name()
	_, err = f.WriteString(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(name, 0644) //nolint:gosec // release notes are public
	}
	if err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return err
	}
	return nil
}

// Changelog assembles the complete changelog from the release notes, newest
// release first.
func Changelog(root string) (string, error) {
	paths, err := releaseFiles(root)
	if err != nil {
		return "", err
	}
	releases := make([]*release, 0, len(paths))
	for _, path := range paths {
		r, err := loadRelease(path)
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

// releaseFiles returns the release notes, sorted by name. Unlike a glob, a
// directory that cannot be read is reported instead of read as empty.
func releaseFiles(root string) ([]string, error) {
	dir, err := resolve(root, ReleasesDir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	paths := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			continue
		}
		if e.Type()&fs.ModeSymlink != 0 {
			return nil, fmt.Errorf("%s: symbolic links are not supported", e.Name())
		}
		paths = append(paths, filepath.Join(dir, e.Name()))
	}
	sort.Strings(paths)
	return paths, nil
}

// findRelease returns the release note file already published for the version,
// or an empty string.
func findRelease(root, version string) (string, error) {
	paths, err := releaseFiles(root)
	if err != nil {
		return "", err
	}
	for _, path := range paths {
		if _, tag, ok := splitName(filepath.Base(path)); ok && tag == version {
			return filepath.Join(ReleasesDir, filepath.Base(path)), nil
		}
	}
	return "", nil
}

func removeUnreleased(root string) error {
	paths, err := unreleasedFiles(root)
	if err != nil {
		return err
	}
	removed := make(map[string]*removedFile, len(paths))
	for _, path := range paths {
		file, err := readRemovable(path)
		if err == nil {
			err = os.Remove(path)
		}
		if err != nil {
			// A release is already written by this point, so the change files
			// go back rather than some of them being lost to a failed cleanup.
			restoreUnreleased(removed)
			return err
		}
		removed[path] = file
	}
	return pruneDirs(filepath.Join(root, UnreleasedDir))
}

// removedFile holds a change file so it can be put back.
type removedFile struct {
	data []byte
	mode fs.FileMode
}

func readRemovable(path string) (*removedFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path) //nolint:gosec // paths come from the repository
	if err != nil {
		return nil, err
	}
	return &removedFile{data: data, mode: info.Mode().Perm()}, nil
}

// restoreUnreleased puts back the change files removed before a failure.
func restoreUnreleased(removed map[string]*removedFile) {
	for path, file := range removed {
		_ = os.WriteFile(path, file.data, file.mode)
	}
}

// pruneDirs removes the subdirectories left empty by the change files.
func pruneDirs(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sub := filepath.Join(dir, e.Name())
		if err := pruneDirs(sub); err != nil {
			return err
		}
		_ = os.Remove(sub) // only succeeds once empty
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

// fenceMarker returns the run of backticks or tildes opening or closing a
// fenced code block.
func fenceMarker(line string) string {
	trimmed := strings.TrimLeft(line, " ")
	for _, c := range []byte{'`', '~'} {
		if len(trimmed) < 3 || trimmed[0] != c {
			continue
		}
		n := 0
		for n < len(trimmed) && trimmed[n] == c {
			n++
		}
		if n >= 3 {
			return trimmed[:n]
		}
	}
	return ""
}

// isHeading reports whether the line opens any heading level.
func isHeading(line string) bool {
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	return n > 0 && n <= 6 && strings.HasPrefix(line[n:], " ")
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
