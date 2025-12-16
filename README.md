# upver (Upstream Versioning)

`upver` is a small, config-driven CLI that computes and applies a repo version bump based on an “upstream” version reference, and can generate a changelog section from Git history.

## Example

A root-level `zarf.yaml` file contains a `metadata.version` field that mirrors the semantic version of an upstream container image referenced in the `images` section of the same configuration file.

```yaml
kind: ZarfPackageConfig
metadata:
  name: upstream-example
  version: 1.2-seq.0 # <-- 1.2 is based on the upstream reference
components:
  - name: test
    images:
      - registry.com/upstream-example:1.2 # <-- this tag is the upstream reference
```

In the root of your repo, create an upver.yaml configuration file that tells upver how to locate and update the version of your project. In this example, the project version is defined by the metadata.version field in zarf.yaml.

```yaml
# upver.yaml
version: # The version of your project
  file: zarf.yaml
  pattern: metadata:\s*\n(?:[ \t].*\n)*?[ \t]*version:\s*"?([^"\n]+)"? # Regex describing your versioning scheme
  group: 1 # Capture group for the version string
```

An upstream version reference is also defined in upver.yaml. The idea is that somewhere in your project there is a literal reference to an upstream version that you want to mirror in your own versioning scheme. In this example, the upstream version is the tag of a container image referenced in the images section of the same zarf.yaml file.

```yaml
upstream:
  file: zarf.yaml
  pattern: registry\.com/upstream-example:([^\s"]+)
  group: 1
```

This works well for projects following trunk-based development.

In this workflow, a developer of this Zarf Package merges a short-lived working branch into a protected main branch, triggering a CI job that runs upver. The tool computes the next project version based on the upstream reference and updates the zarf.yaml version field as defined by version.pattern.

By default, upver increments the sequence portion of the version (e.g. .seq.0 → .seq.1). However, if the upstream reference changes, the base version is updated to match the upstream tag and the sequence is reset to .0.

The goal is to mirror the upstream version of a dependency while still indicating project-specific iterations layered on top of that upstream release. Another example use case could be a Helm chart that references an application version in values.yaml and/or Chart.yaml.

### Changelog Generation

Optionally, upver can generate or update a changelog section based on Git history between the previous and new version tags.

Changelog generation relies on Git being available in the environment where upver is run and uses first-parent, non-merge commits to reflect what advanced the mainline history between releases.
