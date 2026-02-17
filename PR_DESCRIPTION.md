## chains: fix SequentialChain to preserve intermediate outputs

Fixes https://github.com/tmc/langchaingo/issues/1095

### Description

`SequentialChain.Call()` was replacing the inputs map with each chain's output (`inputs = outputs`), losing all intermediate values. This meant intermediate outputs like `synopsis` were `nil` when accessed from the final result, even if declared in `outputKeys`.

The fix accumulates all known values (original inputs + each chain's outputs) through the execution loop, then returns only the declared `outputKeys` from the accumulated map. This is consistent with how Python LangChain's `SequentialChain` works.

### Root Cause

In `chains/sequential.go`, the loop body:

```go
inputs = outputs  // replaces ALL accumulated values
```

caused each iteration to discard previous inputs and intermediate outputs. A chain later in the sequence could not reference an earlier (non-adjacent) chain's output, and the final result map only contained the last chain's output keys.

### Changes

- `chains/sequential.go` — Accumulate all known values across chains in `SequentialChain.Call()` and return only declared `outputKeys`
- `chains/sequential_test.go` — Added `TestSequentialChainIntermediateOutputs` reproducing the exact scenario from issue #1095

### How to Test

```bash
# Build
go build ./chains/...

# Run sequential chain tests (including the new one)
go test ./chains/ -v -count=1 -run "TestSequential" -timeout 30s

# Full chains test suite
go test ./chains/ -count=1 -timeout 60s
```

### PR Checklist

- [x] Read the [Contributing documentation](https://github.com/tmc/langchaingo/blob/main/CONTRIBUTING.md).
- [x] Read the [Code of conduct documentation](https://github.com/tmc/langchaingo/blob/main/CODE_OF_CONDUCT.md).
- [x] Name your Pull Request title clearly, concisely, and prefixed with the name of the primarily affected package you changed according to [Good commit messages](https://go.dev/doc/contribute#commit_messages) (such as `memory: add interfaces for X, Y` or `util: add whizzbang helpers`).
- [x] Check that there isn't already a PR that solves the problem the same way to avoid creating a duplicate.
- [x] Provide a description in this PR that addresses **what** the PR is solving, or reference the issue that it solves (e.g. `Fixes #123`).
- [x] Describes the source of new concepts.
- [x] References existing implementations as appropriate.
- [x] Contains test coverage for new functions.
- [x] Passes all [`golangci-lint`](https://golangci-lint.run/) checks.

Co-Authored-By: ljluestc <ljluestc@users.noreply.github.com>
