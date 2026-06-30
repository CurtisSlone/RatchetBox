package main

// file-kw: wal write ahead log durability accepted webhook appended disk before enqueued crash between accept

import (
	"encoding/json"
	"os"
	"sync"
)

// kw: wal write ahead log
type Wal struct {
	path string
	file *os.File
	mu   sync.Mutex
}

// kw: open wal path write ahead log
func OpenWal(path string) (*Wal, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &Wal{
		path: path,
		file: f,
	}, nil
}

// kw: wal job write ahead log
func (w *Wal) Append(j *Job) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	_, err = w.file.Write(append(b, '\n'))
	return err
}

// kw: replay wal job write ahead log
func (w *Wal) Replay() ([]*Job, error) {
	b, err := os.ReadFile(w.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	lines := splitLines(b)
	var jobs []*Job
	for _, line := range lines {
		if line == "" {
			continue
		}
		var j Job
		if err := json.Unmarshal([]byte(line), &j); err != nil {
			continue
		}
		jobs = append(jobs, &j)
	}
	return jobs, nil
}

// kw: close wal write ahead log
func (w *Wal) Close() error {
	return w.file.Close()
}

// kw: split lines write ahead log
func splitLines(b []byte) []string {
	var lines []string
	for len(b) > 0 {
		i := 0
		for i < len(b) && b[i] != '\n' {
			i++
		}
		lines = append(lines, string(b[:i]))
		b = b[i+1:]
	}
	return lines
}
