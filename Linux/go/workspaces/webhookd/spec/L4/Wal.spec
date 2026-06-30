name: Wal
role: component
dependsOn: Job
intent: A write-ahead log for durability - every accepted webhook is appended to an on-disk log BEFORE it is enqueued, so a crash between accept and deliver does not lose the job. On startup the log is replayed to re-enqueue jobs that were accepted but may not have been delivered. One JSON-encoded Job per line (append-only).
api:
  - type Wal struct { ... }
  - func OpenWal(path string) (*Wal, error)
  - func (w *Wal) Append(j *Job) error
  - func (w *Wal) Replay() ([]*Job, error)
  - func (w *Wal) Close() error
behavior:
  - "Wal holds the log file path, an *os.File open for appending, and a sync.Mutex (Append is called concurrently by HTTP handlers)."
  - "OpenWal opens path with os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644) for appending and stores the path. Return any error."
  - "Append marshals the job with encoding/json (json.Marshal), then under the mutex writes the bytes followed by a single '\\n' newline to the file. Return any write error. This is the durability point: it must complete before the caller enqueues."
  - "Replay reads the whole log file fresh (os.ReadFile(path)); if the file does not exist yet return (nil, nil). Split on '\\n', skip blank lines, json.Unmarshal each line into a *Job, and return the slice of all recovered jobs. A single malformed line should be skipped (continue), not abort the whole replay."
  - "Close closes the underlying file."
constraints: package main; standard library only (os, encoding/json, sync, bytes or strings); uses Job
