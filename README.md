# upver (Upstream Versioning)

`upver` is a small, config-driven CLI that computes and applies a repo version bump based on an “upstream” version reference, and can generate a changelog section from Git history.

## Example

A root level `zarf.yaml` file contains a `metadata.version` field that aims to mirror the semantic version of an upstream container image referenced in the `images` section of the same configuration file.

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

In the root of your repo, create an `upver.yaml` configuration file that tells `upver` how to target the version of your project -- in this example a Zarf package as defined by the `metadata.version` field in the `zarf.yaml` file.

```yaml
# upver.yaml
version: # The version of your project
  file: zarf.yaml
  pattern: 'metadata:\s*\n(?:[ \t].*\n)*?[ \t]*version:\s*"?([^"\n]+)"?' # Regex describing 
  group: 1 # Capture group for the version string
```

An upstream version reference is also defined by the `upver.yaml` file. The idea is that elsewhere in your project there is a literal reference to an upstream version that you want to mirror in your own project version. In this example, the upstream version is the tag of a container image referenced in the `images` section of the same `zarf.yaml` file.

```yaml
upstream:
  file: zarf.yaml
  pattern: 'registry\.jam-demo\.net/jam/ansys_license_manager/ansys_license_manager:([^\s"]+)'
  group: 1
```

This works well for projects following trunk-based development. A developer of this example Zarf package merges a working branch into a protected main branch, triggering a CI job that runs `upver` to compute the next version based on the upstream reference and update the `zarf.yaml` file's version field as defined by the `version.pattern` field. The default behavior is to bump the `.seq-x` portion of the version. However, in the event that the `upstream.pattern` tag changes, the base version is updated to match the upstream tag and the sequence resets to `.0`.

The goal is to mirror the upstream version of your project while still being able to indicate your own project-specific iterations on top of that upstream version. Another use case could be a Helm chart that references an application version in its `values.yaml` and/or `Chart.yaml` file.

### Changelog Generation

You can also define a changelog section in the `upver.yaml` file that generates a changelog based on Git history between the previous and new version tags. The project relies on Git being available in the environment where `upver` is run.
