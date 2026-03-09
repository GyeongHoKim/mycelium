# 🍄 Mycelium

> Automatically discover and link related notes across your knowledge base — powered by multilingual embeddings, running entirely on your device.

Mycelium is an open-source daemon that watches your local markdown vault, finds semantically related notes, and writes those connections back into your files automatically. Just like the underground fungal network that silently connects trees in a forest, Mycelium works quietly in the background to surface the hidden structure of your knowledge.

---

## How It Works

```
Your Vault (.md files)
        │
        ▼
  Mycelium Daemon  ──→  Ollama (local embeddings)
        │
        ▼
    SQLite DB  (embeddings + similarity scores)
        │
        ▼
  Plugin (Obsidian / Logseq / Foam)
        │
        ▼
  Auto-updated Related section in your notes
```

1. Mycelium watches your vault folder for changes
2. When a note is created or updated, it generates an embedding via a local Ollama model
3. Cosine similarity is computed against all other notes in the vault
4. The plugin writes the top-K related notes back into each file — either as frontmatter or as a `## Related` section at the bottom

Everything runs on your device. No cloud. No API keys.

---

## Features

- **Multilingual** — Supports Korean, English, Japanese, Chinese, and 100+ languages out of the box via `multilingual-e5-small`
- **Local-first** — Embeddings and similarity data stay on your machine
- **Automatic** — Watches for file changes and updates related notes incrementally
- **Non-destructive** — Manual links you write yourself are never overwritten
- **Plugin architecture** — Core logic is editor-agnostic; thin plugins handle each tool's format

---

## Supported Editors

| Editor        | Status       | Related Notes Format                  |
| ------------- | ------------ | ------------------------------------- |
| Obsidian      | ✅ Available | `## Related` section (bottom of note) |
| Logseq        | 🚧 Planned   | `related::` frontmatter property      |
| Foam (VSCode) | 🚧 Planned   | `## Related` section (bottom of note) |

---

## Requirements

- [Ollama](https://ollama.com) — for local embedding generation
- Ollama model: `multilingual-e5-small` (~270MB)

```bash
ollama pull multilingual-e5-small
```

---

## Installation

**Obsidian**

1. Open Settings → Community Plugins → Browse
2. Search for `Mycelium`
3. Install and enable
4. The plugin will automatically download and start the Mycelium daemon on first run

Or install manually by copying `plugins/obsidian` into your vault's `.obsidian/plugins/` folder.

**Manual daemon installation (advanced)**

If you prefer to manage the daemon yourself:

```bash
# macOS / Linux
brew install mycelium        # coming soon

# or build from source:
git clone https://github.com/your-handle/mycelium
cd mycelium
go build -o mycelium ./cmd/mycelium
```

---

## Output Format

Mycelium adds a managed section at the bottom of each note. It only touches content between the markers — everything else in your note is untouched.

```markdown
# My Note on HLS Streaming

...your content here...

<!-- mycelium:start -->

## Related

- [[go-astiav segment lifecycle]] · 0.94
- [[HLS.js buffer timeout]] · 0.91
- [[WebCodecs VideoFrame memory]] · 0.87
<!-- mycelium:end -->
```

The `<!-- mycelium:start -->` / `<!-- mycelium:end -->` markers tell the plugin exactly which lines to replace on the next update. You can freely add your own links above or below this section.

---

## Configuration

Configuration lives in `~/.mycelium/config.toml`:

```toml
[vault]
path = "/Users/you/Documents/vault"

[embedding]
model   = "multilingual-e5-small"
ollama  = "http://localhost:11434"

[similarity]
top_k     = 5        # how many related notes to show per file
threshold = 0.75     # minimum similarity score (0.0 ~ 1.0)

[output]
format = "section"   # "section" | "frontmatter"
```

---

## Architecture

```
mycelium/
├── cmd/
│   └── mycelium/            # daemon entry point
│
├── internal/                # private Go packages
│   ├── indexer/             # vault scan, file watcher (fsnotify)
│   ├── embedder/            # Ollama API client
│   ├── similarity/          # cosine similarity, ranking
│   ├── db/                  # SQLite (notes + similarity cache)
│   └── ipc/                 # UDS (macOS/Linux) / Named Pipe (Windows)
│
├── plugins/
│   ├── obsidian/            # TypeScript
│   ├── logseq/              # TypeScript (planned)
│   └── foam/                # TypeScript (planned)
│
├── scripts/                 # build & release scripts
└── .github/workflows/       # CI/CD
```

### IPC

The daemon and plugin communicate over a Unix Domain Socket (macOS/Linux) or Named Pipe (Windows) — no network stack, no port conflicts.

```
Plugin (Node.js / TypeScript)
    └─ IPC connection
           └─ Mycelium Daemon (Go)
                  ├─ SQLite
                  └─ Ollama
```

### Scorer interface

The similarity engine is abstracted behind a `Scorer` interface, making it straightforward to swap or combine algorithms:

```go
type Scorer interface {
    Index(notes []Note) error
    Similar(note Note, topK int) ([]ScoredNote, error)
}
```

Current implementation uses multilingual embeddings + cosine similarity. A BM25-based scorer (for offline / English-only use cases) is planned as an alternative.

---

## Privacy

- All processing happens locally on your machine
- Notes are never sent to any external server
- Embeddings are stored in a local SQLite database at `~/.mycelium/index.db`
- Ollama runs entirely offline after the initial model download

---

## Roadmap

- [x] Core daemon (Go)
- [x] Obsidian plugin
- [ ] Logseq plugin
- [ ] Foam plugin
- [ ] BM25 scorer (offline fallback, no Ollama required)
- [ ] Hybrid scoring (BM25 + embedding)
- [ ] CLI mode (`mycelium similar <note-path>`)
- [ ] Tag-based similarity boost

---

## Contributing

Contributions are welcome. The daemon is editor-agnostic — if you want to build a plugin for another editor, the IPC protocol is straightforward to implement in any language.

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup.

---

## License

MIT License — see [LICENSE](./LICENSE)

---

## Name

Mycelium is the underground fungal network that silently connects trees in a forest, passing nutrients and signals between them. Notes in a knowledge base are like trees — individually complete, but richer when connected. The author's family farms mushrooms in Gangwon Province, Korea.
