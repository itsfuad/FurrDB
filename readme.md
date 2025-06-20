# 🦊 FurrDB

**FurrDB** is a minimal, Redis-inspired in-memory key-value store written in pure Go.

![Cover](./cover.png)

- **Persistent disk storage**
- **Built-in script registration and execution (by hash)**
- **Simple custom command protocol over TCP**
- **No external libraries or dependencies**

> Designed to be lightweight, scriptable, and fast — without the bulk of Redis.

---

## 📌 Features

- [x] In-memory key-value store
- [x] TCP server with custom text protocol
- [x] Command set: `SET`, `GET`, `DEL`, `EXISTS`, `PING`, `EVAL`, `REGSCRIPT`, `RUNSCRIPT`, `EXPIRE`, `TTL`, `SNAPSHOT`
- [x] Append-only persistence log
- [x] Script registration and hash-based invocation
- [x] Basic CLI client
- [x] REPL (local interactive shell)
- [x] TTL expiration
- [x] Snapshot-based persistence
- [x] Scripting sandbox (embedded DSL: LET, IF, END)

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
| `LPUSH k v [v..]` | Push value(s) to head of list             |
| `RPUSH k v [v..]` | Push value(s) to tail of list             |
| `LPOP k`        | Pop value from head of list                 |
| `RPOP k`        | Pop value from tail of list                 |
| `LRANGE k s e`  | Get list elements from s to e               |
| `SADD k v [v..]`| Add value(s) to set                         |
| `SREM k v [v..]`| Remove value(s) from set                    |
| `SMEMBERS k`    | List all set members                        |
| `KEYS`          | List all keys                               |
| `FLUSHDB`       | Clear the database                          |
| `INFO`          | Show server info/stats                      |
| `PING`          | Responds with `PONG`                        |
| `REGSCRIPT s`   | Register script `s`, returns hash           |
| `RUNSCRIPT h`   | Run registered script by hash               |
| `EVAL s`        | Evaluate script string without storing it   |
| `SAVE`          | Force persistence flush                     |
| `EXIT`          | Close the connection                        |

### 📝 Command Examples

#### String
```
SET foo bar
GET foo
DEL foo
EXISTS foo
```

#### List
```
LPUSH mylist a b
RPUSH mylist c
t# mylist is now [b, a, c]
LPOP mylist   # returns b
RPOP mylist   # returns c
LRANGE mylist 0 1  # returns a
```

#### Set
```
SADD myset x y z
SMEMBERS myset   # returns x,y,z
SREM myset y
SMEMBERS myset   # returns x,z
```

#### Meta
```
KEYS        # returns all keys
FLUSHDB     # clears the database
INFO        # returns keys:<count>
```

---

## 🔐 Scripts

FurrDB supports a **basic script runner** and a minimal embedded DSL for scripting.

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
- **Embedded DSL:**
  - Variable assignment: `LET x = GET foo`
  - Conditionals: `IF x == bar ... END`
  - Only whitelisted commands allowed in scripts (sandboxed)
  - Script length limit for safety

**DSL Example:**
```
LET x = GET foo;
IF x == bar;
  SET foo baz;
END;
GET foo
```

---

## 💾 Persistence

- All write commands (`SET`, `DEL`, etc.) are logged to an `aof.log`
- On startup, this file is replayed to reconstruct the state
- Scripts and their hashes are also persisted

---

## 🚀 Getting Started

### 1. **Clone the Repository**
```bash
git clone https://github.com/itsfuad/FurrDB.git
cd FurrDB
```

### 2. **Build the Server**
```bash
go build -o furrdb ./cmd/furrdb
```

### 3. **Run the Server**
```bash
./furrdb
```
The server will start on `localhost:7070` by default.

### 4. **Use the REPL (Interactive Shell)**
```bash
go run ./cmd/furrdb --repl
```
Type commands directly, or use `HELP` for a list of commands.

### 5. **Use the CLI Client to Run Scripts**
```bash
go run ./cmd/furrdbcli scripts/all_commands.txt
```
Or connect manually with telnet:
```bash
telnet localhost 7070
```

### 6. **Using Release Binaries**
- Download the latest release from the [GitHub Releases page](https://github.com/itsfuad/FurrDB/releases).
- Extract and run the binary for your OS:
  - `furrdb` (server)
  - `furrdbcli` (CLI client)

### 7. **Run Tests**
```bash
go test ./...
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

## 📝 License
GNU GENERAL PUBLIC LICENSE