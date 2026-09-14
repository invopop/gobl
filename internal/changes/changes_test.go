package changes_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/invopop/gobl/internal/changes"
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

	t.Run("without pending files", func(t *testing.T) {
		root := repo(t, nil)
		_, err := changes.Preview(root, "v0.506.0", testDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no change files found")
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

	t.Run("without releases", func(t *testing.T) {
		root := repo(t, nil)
		log, err := changes.Changelog(root)
		require.NoError(t, err)
		assert.Contains(t, log, "# Change Log")
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

	t.Run("leaves canonical notes untouched", func(t *testing.T) {
		root := repo(t, map[string]string{
			"changes/releases/2026-09-09-v0.505.0.md": "# v0.505.0 - 2026-09-09\n\n## Added\n\n- one\n",
		})
		names, err := changes.Normalize(root)
		require.NoError(t, err)
		assert.Empty(t, names)
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
		require.NoError(t, os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), []byte(data), 0o644))
	}
	return root
}

func read(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	require.NoError(t, err)
	return string(data)
}
