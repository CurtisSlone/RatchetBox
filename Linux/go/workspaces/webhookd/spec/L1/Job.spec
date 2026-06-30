name: Job
role: data
intent: One webhook delivery job - a payload to POST to a destination URL, plus its current delivery state. This is the core data record the whole dispatcher moves through the queue and (later) the write-ahead log.
api:
  - type Job struct { ID string; URL string; Payload []byte; State string; Attempts int }
  - func NewJob(id, url string, payload []byte) *Job
behavior:
  - "State is one of the string constants StatePending = \"pending\", StateDelivered = \"delivered\", StateFailed = \"failed\". Declare them as exported consts."
  - "NewJob returns &Job{ID: id, URL: url, Payload: payload, State: StatePending, Attempts: 0}."
  - "Job is a plain data holder - no methods beyond the constructor. It must be JSON-serialisable with encoding/json (all fields exported), because later layers persist it to a write-ahead log as one JSON object per line."
constraints: package main; standard library only
