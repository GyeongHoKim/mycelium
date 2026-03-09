# Contributing to Mycelium

Thank you for your interest in contributing.

## Development setup

- **Daemon (Go):** Requires Go 1.21+. From the repo root: `go build -o mycelium ./cmd/mycelium`
- **Plugins (TypeScript):** See each plugin's directory under `plugins/` for editor-specific setup (e.g. Obsidian, Logseq, Foam).
- **Commit messages:** This project uses [Conventional Commits](https://www.conventionalcommits.org/) (enforced via commitlint and lefthook).

## Architecture overview

- The **daemon** (Go) watches the vault, computes embeddings and similarity, and exposes results over IPC.
- **Plugins** (TypeScript) run inside the editor, talk to the daemon via IPC, and write the "Related" section into note files.
- Embeddings live in a **vector DB**; note metadata (path, title, mtime) lives in **SQLite**. See [README.md](./README.md#architecture) for details.

## Questions?

Open an issue or discussion on the repository.
