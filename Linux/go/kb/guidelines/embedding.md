# Embedding (composition over inheritance)

Go has no inheritance; it composes via embedding (Effective Go). An embedded field's exported methods
and fields are promoted to the outer type.

- Embed a type by naming it without a field name. The outer type gets the embedded type's methods.
- Embed interfaces to compose larger interfaces; embed structs (often a pointer) to reuse behavior.
- Promoted methods satisfy interfaces, so embedding a `*log.Logger` makes the outer type loggable.
- A name declared on the outer type shadows the same name promoted from an embedded type.

```go
// Struct embedding: Job gets Logger's methods (Print, Printf, ...) promoted.
type Job struct {
	Command string
	*log.Logger
}

func NewJob(cmd string, logger *log.Logger) *Job {
	return &Job{cmd, logger}
}

job := NewJob("build", log.New(os.Stderr, "Job: ", log.Ldate))
job.Println("starting now...") // promoted from *log.Logger

// Interface embedding composes a bigger interface from smaller ones.
type ReadWriter interface {
	io.Reader
	io.Writer
}
```
