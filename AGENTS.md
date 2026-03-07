# Agent Guidelines: Engine Test Coverage

This document provides specialized instructions for AI agents tasked with increasing test coverage for the `internal/engine` module.

## 🎯 Goal

Achieve **80%+ test coverage** for all packages within `internal/engine`, prioritizing the most critical systems first.

## 🔝 Priorities

1. **Physics (`physics/body`, `physics/movement`, `physics/space`)**: Currently low coverage (9.7% - 57.8%). Critical for game stability.
2. **Scene Management (`scene`)**: Low coverage (14.5%). Essential for game flow and state transitions.
3. **Audio (`audio`)**: Low coverage (17.1%).
4. **Rendering (`render/*`)**: Core visual components.

## 🛠 Testing Strategy & Patterns

### 1. Table-Driven Tests

Prefer table-driven tests for logic with multiple input/output scenarios (e.g., movement, collisions, math).

```go
tests := []struct {
    name    string
    input   int
    want    int
}{
    {"Case A", 1, 2},
    {"Case B", 2, 4},
}
```

### 2. Mocking & Contracts

- Use the interfaces in `internal/engine/contracts/` to create mock implementations for testing.
- **Reusable Mocks**: If a mock is used in more than one test file across different packages, place it in `internal/engine/mocks`. This prevents code duplication.
- **Package-Specific Mocks**: If a mock is only relevant to a single package, define it within the `_test.go` file of that package or a `mocks_test.go` file in the same directory.
- Avoid using actual Ebitengine windows or GPU-dependent code in unit tests.
- Mock `BodiesSpace` to test `Actor` or `Item` updates in isolation.

### 3. Physics & Fixed-Point Arithmetic

- Always validate positions using `fp16.From16()` and `fp16.To16()` when checking `x16` and `y16` values.
- Test edge cases for collisions:
  - One pixel before collision.
  - Partial overlap.
  - Full overlap.
  - Multiple collidables in one space.
  - Fast movement (skipping over thin walls).

### 4. Scene Lifecycle

- Test that `OnStart()`, `Update()`, `Draw()`, and `OnFinish()` are called in the correct order.
- Validate `NavigateTo` and `NavigateBack` logic using a mock `SceneManager`.

### 5. Headless Ebitengine

- For tests that require an `ebiten.Image`, use `ebiten.NewImage(w, h)` in a headless environment.
- Avoid tests that depend on human interaction or specific frame timings (use `timing` package mocks).

## 📋 Standard Workflow for Agents

1. **Analyze Coverage**: Run `go test ./internal/engine/[package] -coverprofile=coverage.out && go tool cover -func=coverage.out`.
2. **Identify Gaps**: Read the source file and identify functions or branches with 0% coverage.
3. **Create Test File**: If it doesn't exist, create `[filename]_test.go`.
4. **Write Tests**: Follow the patterns above. Ensure you test both "happy paths" and error/edge cases.
5. **Verify**: Run the test and check the new coverage percentage.

## ⚠️ Precautions

- **Do not modify production code** unless you find a bug that makes it untestable (e.g., global state that needs to be injected).
- **Keep tests fast**. Avoid long `time.Sleep` calls; use virtual time or frame counters.
- **No Flaky Tests**. Ensure tests are deterministic.
- **No `_ = variable` Pattern**. Do not use `_ = variable` to silence unused variable warnings. Use blank identifier in parameter lists instead: `func (t *T) Method(_ Type) {}`

## Code Style: Avoid `_ = variable` Pattern

**Do NOT do this in production code:**

```go
func (t *Transition) Update() {
    _ = t.active  // Bad: clutters code
}
```

**Do this instead:**

```go
func (t *Transition) Update() {}  // Clean: just remove unused field reference
// or
func (t *Transition) Draw(_ *ebiten.Image) {}  // Use blank in param list
```

**Acceptable in tests:** Using `_ = funcCall()` to verify a function doesn't panic without checking return value.

## 🔍 Key Packages to Target

| Package          | Current Coverage | Focus Area                                          |
| :--------------- | :--------------- | :-------------------------------------------------- |
| `physics/body`   | 9.7%             | `ApplyValidPosition`, `CollisionPosition`           |
| `scene`          | 14.5%            | `SceneManager` state machine, `Transition` triggers |
| `audio`          | 17.1%            | `Loader` error handling, Volume state               |
| `physics/skill`  | 32.7%            | Specific skill activations and cooldowns            |
| `render/sprites` | 40.0%            | Layer sorting, frame calculation                    |
