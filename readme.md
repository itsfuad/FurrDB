# 🦊 FurrDB

**FurrDB** is a minimal, Redis-inspired in-memory key-value store written in pure Go.

- **Persistent disk storage**
- **Built-in script registration and execution (by hash)**
- **Simple custom command protocol over TCP**
- **No external libraries or dependencies**

> Designed to be lightweight, scriptable, and fast — without the bulk of Redis.

---

## 📌 Features

- [x] In-memory key-value store
- [x] TCP server with custom text protocol
- [x] Command set: `SET`, `GET`, `DEL`, `EXISTS`, `PING`, `EVAL`, `REGSCRIPT`, `RUNSCRIPT`
- [x] Append-only persistence log
- [x] Script registration and hash-based invocation
- [x] Basic CLI client
- [x] REPL (local interactive shell)
- [ ] TTL expiration (planned)
- [ ] Snapshot-based persistence (planned)
- [ ] Scripting sandbox (planned with embedded DSL)

---

## 🏗️ Project Structure

```
furrdb/
├── cmd/
│   └── furrdb/      # Main server entrypoint
├── client/          # CLI client (minidb-cli)
├── internal/
│   ├── db/          # In-memory data store and command handlers
│   ├── engine/      # Persistence engine (AOF-based)
│   ├── server/      # TCP listener and protocol parser
│   ├── script/      # Script registration, hashing, execution
│   ├── repl/        # Optional local REPL shell
│   └── utils/       # Logging, hashing, and helper functions
├── scripts/         # Sample scripts for testing
├── testdata/        # Persistence and input test files
├── go.mod
└── README.md
```

---

## 🧠 Architecture Overview

### 1. Core Components

#### 📦 `db/` - In-Memory Store
- Simple `map[string]string` store
- Thread-safe operations
- Dispatches commands based on input tokens

#### 💾 `engine/` - Persistence Engine
- Append-Only File (AOF) log of all write commands
- Loads AOF on startup to recover state
- Periodic flush support

#### 🌐 `server/` - TCP Protocol Handler
- Accepts client connections on a configurable port
- Parses line-based commands, e.g.:
  ```
  SET key value
  GET key
  DEL key
  ```
- Routes commands to `db` package

#### 📜 `script/` - Script Manager
- Accepts scripts (as text)
- Stores in-memory with a SHA256 hash key
- Can run stored scripts by hash with arguments
- Scripts can access DB via pre-defined keywords and syntax

#### 👨‍💻 `client/` - CLI Tool
- Connects to the server via TCP
- Allows running commands from terminal or scripts

#### 💬 `repl/` - Local Shell (Optional)
- Runs DB commands directly against in-process store
- Debugging/testing without networking

---

## 📚 Commands

| Command         | Description                                 |
|-----------------|---------------------------------------------|
| `SET k v`       | Set key `k` to value `v`                    |
| `GET k`         | Get value of key `k`                        |
| `DEL k`         | Delete key `k`                              |
| `EXISTS k`      | Check if key exists                         |
| `PING`          | Responds with `PONG`                        |
| `REGSCRIPT s`   | Register script `s`, returns hash           |
| `RUNSCRIPT h`   | Run registered script by hash               |
| `EVAL s`        | Evaluate script string without storing it   |
| `SAVE`          | Force persistence flush                     |
| `EXIT`          | Close the connection                        |

---

## 🔐 Scripts

FurrDB supports a **basic script runner**.

- Register a script:
  ```
  REGSCRIPT SET foo bar; GET foo
  ```
- Run it by hash:
  ```
  RUNSCRIPT <hash>
  ```

Script engine features:
- Line-by-line command interpretation
- Sequential execution
- Shared context with DB store

---

## 💾 Persistence

- All write commands (`SET`, `DEL`, etc.) are logged to an `aof.log`
- On startup, this file is replayed to reconstruct the state
- Scripts and their hashes are also persisted

---

## 🚀 Getting Started

### Build Server

```bash
go build -o furrdb ./cmd/furrdb
```

### Run Server

```bash
./furrdb
```
Server runs on `localhost:7070` by default.

### Use Client

```bash
go run ./client
```
Or connect manually:
```bash
telnet localhost 7070
```

---

## ⚙️ Configuration

| Config      | Default      |
|-------------|--------------|
| Host        | localhost    |
| Port        | 7070         |
| AOF Path    | aof.log      |
| Script File | scripts.db   |

Defaults are hardcoded for simplicity.  
Future config may support `.env` or flags.

---

## 🧪 Testing

- Unit tests can be added to each module (`*_test.go`)
- Example scripts: `scripts/`
- Test data: `testdata/`

---

## 🧱 Implementation Details

- Uses only standard library packages
- No goroutines in DB logic (all connections are handled via one goroutine per client)
- Simple TCP text protocol (space-delimited tokens)
- Commands are dispatched via a map of handlers

---

## 🔄 Planned Extensions

- TTL/Expiration logic for keys
- Lua or custom mini-scripting language with memory sandbox
- Pub/Sub channels
- Clustering (gossip or Raft)
- Optional binary protocol

---

## 🧑‍💻 Author

Made by [@Fuad Hasan](https://github.com/fuadhasan) for educational and performance-tuning fun.

---

## 📝 License

MIT License — Use freely, modify, distribute.
