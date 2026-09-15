## Added

- `changes`: release notes are now written as one file per pull request in `changes/unreleased`, merged into a dated note in `changes/releases` when a version is published. `CHANGELOG.md` is generated from those notes and is no longer edited by hand; see [changes/README.md](changes/README.md).
- `pkg/changes`: the tooling behind the changes directory, for any project that lays its release notes out the same way.
