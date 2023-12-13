# Codeowners

Tool to generate a
[GitHub CODEOWNERS file](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
from multiple CODEOWNERS files throughout the repo. This makes it easier to
manage code ownership in large repos and thereby reduces the number of
irrelevant review requests and blocked PRs.

## Example

By default, GitHub expects one `CODEOWNERS` file in the repos `.github` dir like
this:

```gitignore
# File: .github/CODEOWNERS

* @org/admin-user
src/go @org/go-developer
src/go/lib.go @org/lib-specialist
```

This file tends to get messy and outdated in large repos with many contributors,
leading to lots of unnecessary approval requests in pull requests.

With this tool these files can instead be placed into the directories to which
they refer:

```gitignore
# File: CODEOWNERS
# Root CODEOWNERS file that sets the default owner of everything in this repo

@org/admin-user  # No glob required here, target is taken from the location of this CODEOWNERS file
```

```gitignore
# File: src/go/CODEOWNERS
# Tiny nested file that sets the owner of everything under src/go

@org/go-developer
lib.go @org/lib-specialist
```

Note that in the second file, ownership for an individual file can still be
assigned as expected. Patterns like `*.go` can be used as well, though they
should only refer to the subdirectory they are located in.

`codeowners` will traverse all CODEOWNERS files from the provided `path`
(respecting `.gitignore` files). If no path is provided, the current working
directory is used.

The resulting root CODEOWNERS file is written to the path specified by
`--output` (default: `.github/CODEOWNERS`). To write to stdout, use `-`.

## Installation

Install as a Go tool via `go get github.com/wlynch/codeowners`.

## Use as GitHub Action

For maximum convenience it is recommended to run this tool automatically in a
GitHub Action like this:

```yaml
# File: .github/workflows/codeowners.yaml

name: Update CODEOWNERS
on:
  pull_request:
    paths:
      - "**/CODEOWNERS" # Trigger for every CODEOWNERS file in the repo
      - "!.github/CODEOWNERS" # except for the generated file itself
jobs:
  check-codeowners:
    runs-on: ubuntu
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4
      - name: Update CODEOWNERS file
        uses: wlynch/codeowners@main
      # Check that there's no diff
      - uses: chainguard-dev/actions/nodiff@main
        with:
          fixup-command: "codeowners"
```

This workflow runs on every change to a `CODEOWNERS` file and regenerates and
commits the root CODEOWNERS file.
