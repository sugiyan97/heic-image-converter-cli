# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

HEIC Image Converter CLI (`heic-convert`) — converts HEIC/HEIF images to JPEG, written in Go. Single-command tool: no subcommands, only flags. All user-facing strings (CLI messages, error text) are in Japanese — preserve this when adding new output.

## Commands

```bash
make deps           # go mod download && go mod tidy
make build           # build for current platform -> bin/convert (alias for build-go)
make build-release   # build with release ldflags (embeds Version, matches CI)
make test            # go test -v ./... (parallelized across CPU count)
make test-coverage   # test + coverage.html
make lint            # golangci-lint run --config golangci.yml
```

Run a single test: `go test -v ./internal/cli -run TestName`.

CGO must be enabled (default) — the `goheif` HEVC decoder is a cgo binding, so a C compiler is required (Xcode CLT / build-essential / MinGW-w64) but no system `libheif` is needed.

Test fixtures live in `test_images/` (`test.HEIC` — has EXIF+GPS, `test_no_exif.HEIC` — no EXIF at all, from goheif's own test suite). Tests copy these into a temp dir before operating on them; never edit them in place. Referenced as `filepath.Join("..", "..", "test_images", ...)` from `internal/*/`.

## Architecture

Three packages, each with a single responsibility, imported only by `cmd/convert/main.go` → `internal/cli`:

- **`internal/cli`** (`root.go`): Cobra root command, flag parsing, and orchestration. `runConvert` dispatches to one of four mutually exclusive modes based on flags — version display, `--uninstall`, `--check-exif`, `--show-exif` (without `--remove-exif`), or the default convert mode — then each mode resolves its target path (arg or cwd) and either finds files recursively in a directory or validates a single file, via the shared `findFilesByType` helper.
- **`internal/converter`** (`converter.go`): HEIC→JPEG pixel conversion via `goheif.Decode`. Notable subtlety: `goheif.SafeEncoding = true` is set in `init()` because goheif's default YCbCr buffers alias C memory that gets freed before Go can read it — removing this reintroduces an intermittent segfault. `jpeg.Encode` has a fast path for `*image.YCbCr`/`*image.Gray` (goheif's native output) that's used directly; other color models (RGBA/NRGBA/generic, e.g. images with alpha) go through `convertToRGBA`, which composites alpha onto a white background.
- **`internal/exif`** (`exif.go`): EXIF extraction from HEIC (`goheif.ExtractExif`) and read/write/strip on JPEG (via `go-exif/v3` + `go-jpeg-image-structure/v2`). `buildIfdChain` deliberately skips any tag whose raw byte length doesn't match its declared type/count — some cameras (e.g. Apple's padded SceneType tag) write non-conforming values that would otherwise corrupt every subsequent tag's offset when re-encoded.

EXIF handling has two independent code paths depending on the flag combination in `root.go`'s `runConvertMode`: EXIF is either stripped (`RemoveEXIFFromJPEG`) or copied over from the source HEIC (`CopyEXIFFromHEICToJPEG`) — never both attempted. `--show-exif` reads straight from the HEIC file and is independent of what happens to the output JPEG's EXIF.

Errors from per-file operations in directory/batch mode are non-fatal by design: a failing file is logged and skipped so the rest of the batch still completes (see `runConvertMode`, `runCheckEXIF`). Don't change this to fail-fast without discussion.

## Docs of record

- `docs/requirements.md` — REQ-NNN numbered requirements; source of truth for intended behavior.
- `docs/test-cases.md` — TC-NNN-NN test cases mapped to REQs, plus TD-NNN test data requirements. Keep in sync when adding tests or fixtures.
- `docs/specification.md`, `docs/development.md`, `docs/troubleshooting.md` — implementation spec, dev setup, known issues.

When changing behavior that a REQ/TC describes, update the corresponding doc in the same change.

## CI

`.github/workflows/` pins third-party Actions by commit SHA (not tag), with the version in a trailing comment (e.g. `actions/checkout@<SHA> # v6`) — a supply-chain safeguard. Dependabot tracks and bumps both the SHA and the comment; if updating manually, resolve the target tag to its commit SHA yourself.
