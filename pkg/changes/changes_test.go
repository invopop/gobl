package changes_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/invopop/gobl/pkg/changes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDate is the release date used by the tests.
var testDate = time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)

func TestParse(t *testing.T) {
	t.Run("sections and entries", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one\n- two\n\n## Fixed\n\n- three\n"))
		require.Len(t, f.Sections, 2)
		assert.Equal(t, "Added", f.Sections[0].Title)
		assert.Equal(t, []string{"- one", "- two"}, f.Sections[0].Lines)
		assert.Equal(t, "Fixed", f.Sections[1].Title)
		assert.Equal(t, []string{"- three"}, f.Sections[1].Lines)
		assert.Empty(t, f.Title)
		assert.Empty(t, f.Preamble)
	})

	t.Run("third level headings", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("### Added\n\n- one\n"))
		require.Len(t, f.Sections, 1)
		assert.Equal(t, "Added", f.Sections[0].Title)
	})

	t.Run("title", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("# v0.506.0 - 2026-09-14\n\n## Added\n\n- one\n"))
		assert.Equal(t, "v0.506.0 - 2026-09-14", f.Title)
		assert.Empty(t, f.Preamble)
		require.Len(t, f.Sections, 1)
	})

	t.Run("repeated heading", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one\n\n## added\n\n- two\n"))
		require.Len(t, f.Sections, 1)
		assert.Equal(t, []string{"- one", "", "- two"}, f.Sections[0].Lines)
	})

	t.Run("preamble", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("\nA summary.\n\n## Added\n\n- one\n"))
		assert.Equal(t, []string{"A summary."}, f.Preamble)
		require.Len(t, f.Sections, 1)
	})

	t.Run("multi-line entries", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one\n  continued\n"))
		assert.Equal(t, []string{"- one", "  continued"}, f.Sections[0].Lines)
	})

	t.Run("headings inside a code block", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one:\n\n```yaml\n# a comment\n## not a heading\n```\n"))
		require.Len(t, f.Sections, 1)
		assert.Equal(t, []string{"- one:", "", "```yaml", "# a comment", "## not a heading", "```"}, f.Sections[0].Lines)
		assert.Empty(t, f.Strays)
	})

	t.Run("headings inside a tilde fenced block", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one:\n\n~~~md\n## not a heading\n~~~\n"))
		require.Len(t, f.Sections, 1)
		assert.Equal(t, []string{"- one:", "", "~~~md", "## not a heading", "~~~"}, f.Sections[0].Lines)
		assert.Empty(t, f.Strays)
	})

	t.Run("fence closed only by a matching run", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n````\n```\n## still fenced\n````\n\n## Fixed\n\n- two\n"))
		require.Len(t, f.Sections, 2)
		assert.Contains(t, f.Sections[0].Lines, "## still fenced")
		assert.Equal(t, "Fixed", f.Sections[1].Title)
	})

	t.Run("stray headings", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\n\n- one\n\n#### Details\n\n- two\n"))
		assert.Equal(t, []string{"#### Details"}, f.Strays)
	})

	t.Run("carriage returns", func(t *testing.T) {
		f := changes.Parse("example.md", []byte("## Added\r\n\r\n- one\r\n"))
		require.Len(t, f.Sections, 1)
		assert.Equal(t, []string{"- one"}, f.Sections[0].Lines)
	})
}

func TestFragmentValidate(t *testing.T) {
	tests := []struct {
		name string
		data string
		err  string
	}{
		{name: "valid", data: "## Added\n\n- one\n"},
		{
			name: "unknown heading",
			data: "## Fixes\n\n- one\n",
			err:  `example.md: unknown heading "Fixes"`,
		},
		{
			name: "with a title",
			data: "# My change\n\n## Added\n\n- one\n",
			err:  `example.md: unexpected title "My change"`,
		},
		{
			name: "heading at an unsupported level",
			data: "## Added\n\n- one\n\n#### Details\n",
			err:  `example.md: unexpected heading "#### Details"`,
		},
		{
			name: "no heading",
			data: "- one\n",
			err:  "example.md: no entries found",
		},
		{
			name: "empty section",
			data: "## Added\n",
			err:  "example.md: no entries found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := changes.Parse("example.md", []byte(test.data)).Validate()
			if test.err == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.err)
		})
	}
}

func TestNotes(t *testing.T) {
	frags := []*changes.Fragment{
		changes.Parse("a.md", []byte("## Fixed\n\n- fix from a\n\n## Added\n\n- add from a\n")),
		changes.Parse("b.md", []byte("Summary of b.\n\n## Added\n\n- add from b\n")),
	}
	out := changes.Notes("v0.506.0", testDate, frags)
	assert.Equal(t, `# v0.506.0 - 2026-09-14

Summary of b.

## Added

- add from a
- add from b

## Fixed

- fix from a
`, out)
}

func TestNotesPreambles(t *testing.T) {
	frags := []*changes.Fragment{
		changes.Parse("a.md", []byte("Summary of a.\n\n## Added\n\n- one\n")),
		changes.Parse("b.md", []byte("Summary of b.\n\n## Added\n\n- two\n")),
	}
	out := changes.Notes("v0.506.0", testDate, frags)
	assert.Equal(t, `# v0.506.0 - 2026-09-14

Summary of a.

Summary of b.

## Added

- one
- two
`, out)
}

func TestPreview(t *testing.T) {
	t.Run("merges pending files", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/b.md": "## Added\n\n- second\n",
			"changes/unreleased/a.md": "## Added\n\n- first\n",
		})
		out, err := changes.Preview(root, "v0.506.0", testDate)
		require.NoError(t, err)
		assert.Equal(t, "# v0.506.0 - 2026-09-14\n\n## Added\n\n- first\n- second\n", out)
	})

	t.Run("ignores the readme", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/README.md": "# Instructions\n\n## Nonsense\n\n- instructions\n",
			"changes/unreleased/a.md":      "## Added\n\n- first\n",
		})
		out, err := changes.Preview(root, "v0.506.0", testDate)
		require.NoError(t, err)
		assert.Contains(t, out, "- first")
		assert.NotContains(t, out, "instructions")
	})

	t.Run("inside a branch subdirectory", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/claude/nested-work.md": "## Added\n\n- nested\n",
			"changes/unreleased/a.md":                  "## Added\n\n- flat\n",
		})
		out, err := changes.Preview(root, "v0.506.0", testDate)
		require.NoError(t, err)
		assert.Equal(t, "# v0.506.0 - 2026-09-14\n\n## Added\n\n- flat\n- nested\n", out)
	})

	t.Run("reports the path of an invalid nested file", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/claude/nested-work.md": "## Nope\n\n- nested\n",
		})
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "claude/nested-work.md: unknown heading")
	})

	t.Run("without pending files", func(t *testing.T) {
		root := repo(t, nil)
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no change files found")
	})

	t.Run("without an unreleased directory", func(t *testing.T) {
		_, err := changes.Preview(t.TempDir(), "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no change files found")
	})

	t.Run("with an unreadable subdirectory", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable directory does not stop root from walking it")
		}
		root := repo(t, map[string]string{"changes/unreleased/claude/nested.md": "## Added\n\n- one\n"})
		dir := filepath.Join(root, "changes/unreleased/claude")
		require.NoError(t, os.Chmod(dir, 0o000))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
	})

	t.Run("with a symlinked unreleased directory", func(t *testing.T) {
		root := repo(t, nil)
		outside := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(outside, "a.md"), []byte("## Added\n\n- outside\n"), 0o644))
		require.NoError(t, os.RemoveAll(filepath.Join(root, changes.UnreleasedDir)))
		require.NoError(t, os.Symlink(outside, filepath.Join(root, changes.UnreleasedDir)))
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unreleased: symbolic links are not supported")
	})

	t.Run("with a symlinked change file", func(t *testing.T) {
		root := repo(t, nil)
		require.NoError(t, os.WriteFile(filepath.Join(root, "secret.md"), []byte("## Added\n\n- secret\n"), 0o644))
		require.NoError(t, os.Symlink(
			filepath.Join(root, "secret.md"),
			filepath.Join(root, "changes/unreleased/link.md"),
		))
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "link.md: symbolic links are not supported")
	})

	t.Run("with an unreadable change file", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable file does not stop root from reading it")
		}
		root := repo(t, map[string]string{"changes/unreleased/a.md": "## Added\n\n- one\n"})
		path := filepath.Join(root, "changes/unreleased/a.md")
		require.NoError(t, os.Chmod(path, 0o000))
		t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "a.md")
	})

	t.Run("with an invalid file", func(t *testing.T) {
		root := repo(t, map[string]string{"changes/unreleased/a.md": "## Nope\n\n- first\n"})
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown heading")
	})
}

func TestRelease(t *testing.T) {
	t.Run("writes notes and regenerates the changelog", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/a.md":                 "## Added\n\n- first\n",
			"changes/unreleased/README.md":            "instructions\n",
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- previous\n",
		})

		name, err := changes.Release(root, "v0.506.0", testDate)
		require.NoError(t, err)
		assert.Equal(t, "changes/releases/2026-09-14-v0.506.0.md", name)

		assert.Equal(t,
			"# v0.506.0 - 2026-09-14\n\n## Added\n\n- first\n",
			read(t, root, name),
		)

		assert.NoFileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/README.md"))

		log := read(t, root, changes.ChangelogFile)
		assert.Contains(t, log, "# Change Log")
		assert.Contains(t, log, "## [v0.506.0] - 2026-09-14\n\n### Added\n\n- first\n\n## [v0.505.0] - 2026-09-09")
	})

	t.Run("clears branch subdirectories", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/claude/nested-work.md": "## Added\n\n- nested\n",
		})
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.NoError(t, err)
		assert.NoDirExists(t, filepath.Join(root, "changes/unreleased/claude"))
		assert.DirExists(t, filepath.Join(root, "changes/unreleased"))
	})

	t.Run("keeps the change files when the changelog cannot be built", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/a.md":                 "## Added\n\n- first\n",
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-01-01\n",
		})
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected")
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
		assert.NoFileExists(t, filepath.Join(root, "changes/releases/2026-09-14-v0.506.0.md"))
		assert.NoFileExists(t, filepath.Join(root, changes.ChangelogFile))
	})

	t.Run("keeps the changelog when the history cannot be read", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable directory does not stop root from reading it")
		}
		root := repo(t, map[string]string{
			"changes/unreleased/a.md":                 "## Added\n\n- first\n",
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- previous\n",
			"CHANGELOG.md":                            "# Change Log\n\n## [v0.505.0] - 2026-09-09\n",
		})
		dir := filepath.Join(root, changes.ReleasesDir)
		require.NoError(t, os.Chmod(dir, 0o333))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, read(t, root, changes.ChangelogFile), "v0.505.0")
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
	})

	t.Run("creates the releases directory for a first release", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(root, changes.UnreleasedDir), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(root, changes.UnreleasedDir, "a.md"), []byte("## Added\n\n- first\n"), 0o644))

		name, err := changes.Release(root, "v0.1.0", testDate)
		require.NoError(t, err)
		assert.Equal(t, "changes/releases/2026-09-14-v0.1.0.md", name)
		assert.Equal(t, "# v0.1.0 - 2026-09-14\n\n## Added\n\n- first\n", read(t, root, name))
		assert.Contains(t, read(t, root, changes.ChangelogFile), "## [v0.1.0] - 2026-09-14")
	})

	t.Run("with a version that cannot be stored in a file name", func(t *testing.T) {
		versions := map[string]string{
			"v0.506.0/../../evil": "cannot contain a path separator",
			`v0.506.0\evil`:       "cannot contain a path separator",
			"0.506.0":             "cannot be read back",
			"":                    "cannot be read back",
		}
		for version, msg := range versions {
			root := repo(t, map[string]string{"changes/unreleased/a.md": "## Added\n\n- first\n"})
			_, err := changes.Release(root, version, testDate)
			require.Error(t, err, version)
			assert.Contains(t, err.Error(), msg, version)
			assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
			assert.NoFileExists(t, filepath.Join(root, changes.ChangelogFile))
		}
	})

	t.Run("puts the change files back when the cleanup fails", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("a read-only directory does not stop root from removing files")
		}
		root := repo(t, map[string]string{
			"changes/unreleased/a.md":    "## Added\n\n- first\n",
			"changes/unreleased/zz/b.md": "## Added\n\n- second\n",
		})
		dir := filepath.Join(root, "changes/unreleased/zz")
		require.NoError(t, os.Chmod(dir, 0o555))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Equal(t, "## Added\n\n- first\n", read(t, root, "changes/unreleased/a.md"))
		assert.Equal(t, "## Added\n\n- second\n", read(t, root, "changes/unreleased/zz/b.md"))
	})

	t.Run("with the version already released", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/unreleased/a.md":                 "## Added\n\n- first\n",
			"changes/releases/2026-09-09-v0.506.0.md": "# v0.506.0 - 2026-09-09\n",
		})
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exist")
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
	})

	t.Run("without pending files", func(t *testing.T) {
		root := repo(t, nil)
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.NoFileExists(t, filepath.Join(root, changes.ChangelogFile))
	})
}

func TestChangelog(t *testing.T) {
	t.Run("nests each release a level deeper", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\nA summary.\n\n## Added\n\n- one\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Contains(t, log, "## [v0.505.0] - 2026-09-09\n\nA summary.\n\n### Added\n\n- one\n")
	})

	t.Run("orders by version", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2024-04-08-v0.71.0.md":      "# v0.71.0 - 2024-04-08\n\n## Added\n\n- a\n",
			"changes/releases/2026-06-10-v0.500.0.md":     "# v0.500.0 - 2026-06-10\n\n## Added\n\n- b\n",
			"changes/releases/2026-04-13-v0.400.0-rc1.md": "# v0.400.0-rc1 - 2026-04-13\n\n## Added\n\n- c\n",
			"changes/releases/2026-04-15-v0.400.0.md":     "# v0.400.0 - 2026-04-15\n\n## Added\n\n- d\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Regexp(t, `(?s)v0\.500\.0.*v0\.400\.0\].*v0\.400\.0-rc1.*v0\.71\.0`, log)
	})

	t.Run("separates releases with a blank line", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
			"changes/releases/2026-08-05-v0.504.0.md": "# v0.504.0 - 2026-08-05\n\n## Added\n\n- two\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Contains(t, log, "- one\n\n## [v0.504.0]")
		assert.True(t, len(log) > 0 && log[len(log)-1] == '\n')
	})

	t.Run("with a mismatched title", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-01-01\n\n## Added\n\n- one\n",
		})
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `title is "v0.505.0 - 2026-01-01", expected "v0.505.0 - 2026-09-09"`)
	})

	t.Run("with an unexpected file name", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/v0.505.0.md": "# v0.505.0 - 2026-09-09\n",
		})
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected a name like")
	})

	t.Run("with an invalid date", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2024-20-21-v0.203.0.md": "# v0.203.0 - 2024-20-21\n",
		})
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected a name like")
	})

	t.Run("with a project header", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/HEADER.md":                       "# GOBL Changes\n\nEverything worth knowing.\n",
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Equal(t, "# GOBL Changes\n\nEverything worth knowing.\n\n## [v0.505.0] - 2026-09-09\n\n### Added\n\n- one\n", log)
	})

	t.Run("orders releases of the same version by file name", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- first\n",
			"changes/releases/2026-09-10-v0.505.0.md": "# v0.505.0 - 2026-09-10\n\n## Added\n\n- second\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Regexp(t, `(?s)2026-09-10.*2026-09-09`, log)
	})

	t.Run("orders releases without a usable version by file name", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-vNOPE.md": "# vNOPE - 2026-09-09\n\n## Added\n\n- first\n",
			"changes/releases/2026-08-05-vALSO.md": "# vALSO - 2026-08-05\n\n## Added\n\n- second\n",
		})
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Regexp(t, `(?s)vNOPE.*vALSO`, log)
	})

	t.Run("with a symlinked release note", func(t *testing.T) {
		root := repo(t, nil)
		require.NoError(t, os.WriteFile(filepath.Join(root, "secret.md"), []byte("anything at all\n"), 0o644))
		require.NoError(t, os.Symlink(
			filepath.Join(root, "secret.md"),
			filepath.Join(root, "changes/releases/2026-09-09-v0.505.0.md"),
		))
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "symbolic links are not supported")
	})

	t.Run("with an unreadable release note", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable file does not stop root from reading it")
		}
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		path := filepath.Join(root, "changes/releases/2026-09-09-v0.505.0.md")
		require.NoError(t, os.Chmod(path, 0o000))
		t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
		_, err := changes.Changelog(root)
		require.Error(t, err)
	})

	t.Run("with an unreadable releases directory", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable directory does not stop root from reading it")
		}
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		dir := filepath.Join(root, changes.ReleasesDir)
		require.NoError(t, os.Chmod(dir, 0o000))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		_, err := changes.Changelog(root)
		require.Error(t, err)
	})

	t.Run("with a symlinked header", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		require.NoError(t, os.WriteFile(filepath.Join(root, "secret.md"), []byte("# Secrets\n"), 0o644))
		require.NoError(t, os.Symlink(filepath.Join(root, "secret.md"), filepath.Join(root, changes.HeaderFile)))
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HEADER.md: symbolic links are not supported")
	})

	t.Run("with a symlinked changes directory", func(t *testing.T) {
		root := t.TempDir()
		outside := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(outside, "releases"), 0o755))
		require.NoError(t, os.WriteFile(
			filepath.Join(outside, "releases", "2026-09-09-v9.9.9.md"),
			[]byte("# v9.9.9 - 2026-09-09\n"), 0o644,
		))
		require.NoError(t, os.Symlink(outside, filepath.Join(root, changes.Dir)))
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "changes: symbolic links are not supported")
	})

	t.Run("with a symlinked releases directory", func(t *testing.T) {
		root := repo(t, nil)
		outside := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(outside, "2026-09-09-v9.9.9.md"), []byte("# v9.9.9 - 2026-09-09\n"), 0o644))
		require.NoError(t, os.RemoveAll(filepath.Join(root, changes.ReleasesDir)))
		require.NoError(t, os.Symlink(outside, filepath.Join(root, changes.ReleasesDir)))
		_, err := changes.Changelog(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "releases: symbolic links are not supported")
	})

	t.Run("with an unreadable header", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		require.NoError(t, os.Mkdir(filepath.Join(root, changes.HeaderFile), 0o755))
		_, err := changes.Changelog(root)
		require.Error(t, err)
	})

	t.Run("ignores anything that is not a release note", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
			"changes/releases/notes.txt":              "scratch\n",
		})
		require.NoError(t, os.Mkdir(filepath.Join(root, "changes/releases/archive"), 0o755))
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Contains(t, log, "## [v0.505.0] - 2026-09-09")
	})

	t.Run("without a releases directory", func(t *testing.T) {
		log, err := changes.Changelog(t.TempDir())
		require.NoError(t, err)
		assert.Contains(t, log, "# Change Log")
	})

	t.Run("without a project header", func(t *testing.T) {
		root := repo(t, nil)
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Contains(t, log, "# Change Log")
		assert.Contains(t, log, "All notable changes to this project")
	})
}

func TestNormalize(t *testing.T) {
	t.Run("merges repeated headings and drops empty ones", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n\n## Removed\n\n## Added\n\n- two\n",
		})
		names, err := changes.Normalize(root)
		require.NoError(t, err)
		assert.Equal(t, []string{"2026-09-09-v0.505.0.md"}, names)
		assert.Equal(t,
			"# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n\n- two\n",
			read(t, root, "changes/releases/2026-09-09-v0.505.0.md"),
		)
	})

	t.Run("reports a release note it cannot read", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-01-01\n\n## Added\n\n- one\n",
		})
		_, err := changes.Normalize(root)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title is")
	})

	t.Run("cannot rewrite a release note", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("a read-only directory does not stop root from writing")
		}
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n\n## Added\n\n- two\n",
		})
		require.NoError(t, os.Chmod(filepath.Join(root, changes.ReleasesDir), 0o555))
		t.Cleanup(func() { _ = os.Chmod(filepath.Join(root, changes.ReleasesDir), 0o755) })
		_, err := changes.Normalize(root)
		require.Error(t, err)
	})

	t.Run("reports an unreadable releases directory", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("an unreadable directory does not stop root from reading it")
		}
		root := repo(t, nil)
		dir := filepath.Join(root, changes.ReleasesDir)
		require.NoError(t, os.Chmod(dir, 0o000))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
		_, err := changes.Normalize(root)
		require.Error(t, err)
	})

	t.Run("leaves canonical notes untouched", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		names, err := changes.Normalize(root)
		require.NoError(t, err)
		assert.Empty(t, names)
	})
}

func TestReleaseWriteFailures(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("a read-only directory does not stop root from writing")
	}
	pending := map[string]string{"changes/unreleased/a.md": "## Added\n\n- first\n"}

	// readOnly makes the directory reject new files for the rest of the test.
	readOnly := func(t *testing.T, dir string) {
		t.Helper()
		require.NoError(t, os.Chmod(dir, 0o555))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	}

	t.Run("cannot write the release note", func(t *testing.T) {
		root := repo(t, pending)
		readOnly(t, filepath.Join(root, changes.ReleasesDir))
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
		assert.NoFileExists(t, filepath.Join(root, changes.ChangelogFile))
	})

	t.Run("cannot write the changelog", func(t *testing.T) {
		root := repo(t, pending)
		readOnly(t, root)
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.FileExists(t, filepath.Join(root, "changes/unreleased/a.md"))
		assert.NoFileExists(t, filepath.Join(root, "changes/releases/2026-09-14-v0.506.0.md"))
		assert.NoFileExists(t, filepath.Join(root, changes.ChangelogFile))
	})

	t.Run("cannot remove the change files", func(t *testing.T) {
		root := repo(t, pending)
		readOnly(t, filepath.Join(root, changes.UnreleasedDir))
		_, err := changes.Release(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.FileExists(t, filepath.Join(root, "changes/releases/2026-09-14-v0.506.0.md"))
		assert.FileExists(t, filepath.Join(root, changes.ChangelogFile))
	})
}

// repo prepares a temporary repository root containing the given files.
func repo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{changes.UnreleasedDir, changes.ReleasesDir} {
		require.NoError(t, os.MkdirAll(filepath.Join(root, dir), 0o755))
	}
	for name, data := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(data), 0o644))
	}
	return root
}

func read(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	require.NoError(t, err)
	return string(data)
}
