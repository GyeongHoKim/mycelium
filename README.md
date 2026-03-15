# 🍄 Mycelium

[![codecov](https://codecov.io/gh/GyeongHoKim/mycelium/branch/main/graph/badge.svg)](https://codecov.io/gh/GyeongHoKim/mycelium)

> Automatically discover and link related notes across your knowledge base — powered by multilingual embeddings, running entirely on your device.

Mycelium is an open-source daemon that watches your local markdown vault, finds semantically related notes, and writes those connections back into your files automatically. Just like the underground fungal network that silently connects trees in a forest, Mycelium works quietly in the background to surface the hidden structure of your knowledge.

---

## How It Works

```mermaid
flowchart TB
  Vault["Your Vault (.md files)"]
  Daemon["Mycelium Daemon"]
  Ollama["Ollama (local embeddings)"]
  VectorDB["Vector DB<br/>(embedding vectors, similarity search)"]
  SQLite["SQLite DB<br/>(note metadata, paths, timestamps)"]
  Plugin["Plugin (Obsidian / Logseq / Foam)"]
  Related["Auto-updated Related section in your notes"]

  Vault --> Daemon
  Daemon --> Ollama
  Daemon --> VectorDB
  Daemon --> SQLite
  Daemon --> Plugin
  Plugin --> Related
```

1. The **daemon** watches your vault folder for changes.
2. When a note is created or updated, the daemon generates an embedding (Ollama), stores it in the **embedded vector DB (chromem-go)**, and keeps note metadata (including content hashes) in SQLite.
3. Similarity is computed via the vector DB. The daemon exposes "related notes" over IPC.
4. The **plugin** (running inside your editor) requests related notes from the daemon and writes them into the file using the editor's internal API. This ensures **"Plugin-as-Writer"** safety, avoiding file system conflicts and preserving undo history.

Everything runs on your device. No cloud. No API keys.

---

## Features

- **Multilingual** — Supports Korean, English, Japanese, Chinese, and 100+ languages out of the box via `qwen3-embedding`
- **Local-first** — Embeddings (via `chromem-go`) and note metadata (SQLite) stay on your machine
- **Intelligent Updates** — Uses a combination of **Debouncing** (waits for typing to stop) and **Content Hashing** (ignores changes within the `## Related` section) to minimize Ollama API calls and prevent update loops.
- **Plugin-as-Writer** — The daemon never modifies your files directly. The plugin handles all writes via editor APIs, ensuring safety and compatibility with editor features like Undo.
- **Single-instance daemon** — Uses a PID file and socket handshake to guarantee only one daemon process runs per vault. A second invocation connects to the existing instance instead of spawning a duplicate.
- **Non-destructive** — Only the content between `<!-- mycelium:start -->` and `<!-- mycelium:end -->` is updated; the rest of the note is left untouched.
- **Plugin architecture** — Core logic is editor-agnostic; thin plugins handle each tool's format

---

## Supported Editors

| Editor        | Status       | Related Notes Format                  | Installation constraint          |
| ------------- | ------------ | ------------------------------------- | -------------------------------- |
| Obsidian      | ✅ Available | `## Related` section (bottom of note) | Native installer only (see note) |
| Logseq        | 🚧 Planned   | `related::` frontmatter property      | —                                |
| Foam (VSCode) | 🚧 Planned   | `## Related` section (bottom of note) | —                                |

> **Obsidian sandbox note.** The Obsidian community plugin store only allows `manifest.json`,
> `styles.css`, and `main.js` in each release — binaries cannot be bundled. Additionally,
> sandboxed Obsidian distributions (Snap, Flatpak, AppImage) block child-process execution.
> The plugin can only download and launch the Mycelium daemon when Obsidian is installed via
> the **native desktop installer** (`.dmg`, `.exe`, `apt` repo). Users must also disable
> **Restricted Mode** in Obsidian settings before enabling the plugin.

---

## Requirements

- [Ollama](https://ollama.com) — for local embedding generation
- Ollama model: `qwen3-embedding` (~270MB)

```bash
ollama pull qwen3-embedding
```

---

## Platform Support

| OS      | Architecture     | IPC transport                      | Binary name                            |
| ------- | ---------------- | ---------------------------------- | -------------------------------------- |
| macOS   | amd64            | Unix Domain Socket                 | `mycelium_{version}_darwin_amd64`      |
| macOS   | arm64 (M-series) | Unix Domain Socket                 | `mycelium_{version}_darwin_arm64`      |
| Linux   | amd64            | Unix Domain Socket                 | `mycelium_{version}_linux_amd64`       |
| Linux   | arm64            | Unix Domain Socket                 | `mycelium_{version}_linux_arm64`       |
| Windows | amd64            | Named Pipe (`\\.\pipe\mycelium`)   | `mycelium_{version}_windows_amd64.exe` |

**Platform detection in the plugin (Node.js/TypeScript):**

```typescript
function resolveBinaryName(version: string): string {
  const osMap:   Record<string, string> = { darwin: "darwin", linux: "linux", win32: "windows" };
  const archMap: Record<string, string> = { x64: "amd64", arm64: "arm64" };
  const ext = process.platform === "win32" ? ".exe" : "";
  return `mycelium_${version}_${osMap[process.platform]}_${archMap[process.arch]}${ext}`;
}
```

> **WSL note.** `process.platform` reports `"linux"` inside WSL — use the Linux binary and
> Unix Domain Socket. Named Pipes are only used on native Windows.

---

## Installation

### Via Obsidian Plugin (recommended)

> **Prerequisite:** Obsidian must be installed via the **native desktop installer**
> (`.dmg`, `.exe`, `apt` repo). Sandboxed distributions (Snap, Flatpak, AppImage) block
> child-process execution. See [Platform Support](#platform-support).

1. Open Settings → Community Plugins → turn off **Restricted Mode**.
2. Browse → search `Mycelium` → Install → Enable.
3. On first activation the plugin:
   - Detects OS and architecture
   - Downloads `checksums.txt` from the matching GitHub release
   - Downloads `mycelium_{version}_{os}_{arch}[.exe]`
   - **Verifies SHA-256 checksum** against `checksums.txt` — refuses to run on mismatch
   - Places binary at `~/.mycelium/mycelium` (or `%APPDATA%\mycelium\mycelium.exe`)
   - Starts the daemon (see [Daemon Lifecycle](#daemon-lifecycle))

### Manual Installation

```bash
VERSION="0.1.0"
OS="darwin"      # darwin | linux | windows
ARCH="arm64"     # amd64 | arm64

curl -LO "https://github.com/GyeongHoKim/mycelium/releases/download/v${VERSION}/checksums.txt"
curl -LO "https://github.com/GyeongHoKim/mycelium/releases/download/v${VERSION}/mycelium_${VERSION}_${OS}_${ARCH}"

# Verify (macOS/Linux)
sha256sum --check --ignore-missing checksums.txt

chmod +x mycelium_${VERSION}_${OS}_${ARCH}
mv mycelium_${VERSION}_${OS}_${ARCH} /usr/local/bin/mycelium
```

### Build from Source

```bash
git clone https://github.com/GyeongHoKim/mycelium
cd mycelium
go build -o mycelium ./cmd/mycelium
```

No cgo required — SQLite driver (`modernc.org/sqlite`) is pure Go.

---

## Output Format

The **plugin** adds a managed section at the bottom of each note (or in frontmatter for editors like Logseq). It only replaces content between the markers — everything else in your note is untouched.

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

Configuration lives in `~/.mycelium/config.toml` (one vault path per config; one daemon instance typically serves that vault, and multiple plugins can connect to the same daemon):

```toml
[vault]
path = "/Users/you/Documents/vault"

[embedding]
model   = "qwen3-embedding"
ollama  = "http://localhost:11434"

[similarity]
top_k     = 5        # how many related notes to show per file
threshold = 0.75     # minimum similarity score (0.0 ~ 1.0)

[output]
format = "section"   # "section" | "frontmatter" (default per editor: Obsidian → section, Logseq → frontmatter)
```

---

## Architecture

The repo layout below is the target design. The `vectordb/` package (vector store for embeddings and ANN search) is planned; until then, similarity may be implemented with an in-memory or external vector store.

```mermaid
flowchart TB
  subgraph mycelium["mycelium/"]
    subgraph cmd["cmd/"]
      main["mycelium/ (daemon entry point)"]
    end
    subgraph internal["internal/ (private Go packages)"]
      watcher["watcher/ (vault scan, fsnotify)"]
      embedder["embedder/ (Ollama API)"]
      similarity["similarity/ (backed by vector DB)"]
      db["db/ (SQLite, note metadata)"]
      vectordb["vectordb/ (embeddings, ANN) — planned"]
      ipc["ipc/ (UDS / Named Pipe)"]
    end
    subgraph plugins["plugins/"]
      obsidian["obsidian/ (TypeScript)"]
      logseq["logseq/ (planned)"]
      foam["foam/ (planned)"]
    end
    scripts["scripts/"]
    workflows[".github/workflows/"]
  end
```

### IPC

The daemon and plugin communicate over a **Unix Domain Socket** at `~/.mycelium/mycelium.sock`
(macOS/Linux) or a **Named Pipe** (`\\.\pipe\mycelium`) on Windows. The protocol is
**JSON-RPC 2.0** — the same framing as LSP and MCP — so any language with a JSON library
can implement a client.

```mermaid
flowchart LR
  Plugin["Plugin\n(Node.js / TypeScript)"]
  Daemon["Mycelium Daemon\n(Go)"]
  SQLite["SQLite\n(note metadata)"]
  VectorDB["Vector DB\n(embeddings)"]
  Ollama["Ollama"]

  Plugin -->|"JSON-RPC 2.0\nUDS / Named Pipe"| Daemon
  Daemon --> SQLite
  Daemon --> VectorDB
  Daemon --> Ollama
```

#### Version Handshake

Every connection begins with an `initialize` exchange before any other method.
Major version mismatch → plugin must abort and prompt user to upgrade.

```jsonc
// Plugin → Daemon
{ "jsonrpc": "2.0", "id": 1, "method": "initialize",
  "params": { "protocolVersion": "1.0",
               "clientInfo": { "name": "mycelium-obsidian", "version": "0.3.0" } } }

// Daemon → Plugin
{ "jsonrpc": "2.0", "id": 1,
  "result": { "protocolVersion": "1.0",
               "serverInfo": { "name": "mycelium", "version": "0.1.0" } } }
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

## Daemon Lifecycle

### State Files

| File                        | Windows equivalent                   | Purpose                     |
| --------------------------- | ------------------------------------ | --------------------------- |
| `~/.mycelium/mycelium.sock` | `\\.\pipe\mycelium`                  | IPC endpoint                |
| `~/.mycelium/daemon.pid`    | `%APPDATA%\mycelium\daemon.pid`      | Running PID for dedup check |
| `~/.mycelium/index.db`      | `%APPDATA%\mycelium\index.db`        | SQLite note metadata        |
| `~/.mycelium/config.toml`   | `%APPDATA%\mycelium\config.toml`     | User configuration          |

### Startup Flow (Plugin perspective)

```mermaid
flowchart TD
  A([Plugin activates]) --> B{socket file exists?}
  B -- No --> F[Spawn daemon]
  B -- Yes --> C[Try connect + initialize handshake]
  C -- Success --> D([Use existing daemon])
  C -- Fail: stale socket --> E[Delete .sock and .pid files]
  E --> F
  F --> G[Wait for socket — max 5 s / 100 ms poll]
  G --> H[Connect + initialize handshake]
  H --> D
```

### Shutdown Flow (Daemon perspective)

```mermaid
flowchart TD
  A([SIGTERM / SIGINT]) --> B[Stop accepting new connections]
  B --> C[Drain in-flight requests — max 10 s]
  C --> D[Close IPC listener]
  D --> E[Delete .sock file]
  E --> F[Delete .pid file]
  F --> G([Exit 0])
```

If the daemon crashes without cleanup, the stale socket causes the next `connect` to fail
(connection refused), which triggers the cleanup path in the startup flow above.

---

## Privacy

- All processing happens locally on your machine
- Notes are never sent to any external server
- Note metadata is stored in SQLite at `~/.mycelium/index.db`; embedding vectors are stored in a local vector DB (e.g. under `~/.mycelium/`)
- Ollama runs entirely offline after the initial model download

---

## Roadmap

- [x] Core daemon (Go)
- [x] SQLite schema with content hashing
- [x] Ollama embedder (`qwen3-embedding`)
- [ ] `chromem-go` integration (Vector DB)
- [ ] Intelligent update logic (Debounce + Hashing)
- [ ] IPC: JSON-RPC 2.0 over UDS / Named Pipe
- [ ] Daemon lifecycle (PID file, socket, graceful shutdown)
- [ ] Binary download + SHA-256 verification in plugin
- [ ] Obsidian plugin (v2 with IPC fetch)
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
