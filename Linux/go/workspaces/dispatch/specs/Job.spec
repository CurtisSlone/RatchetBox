name: Job
role: component
intent: The unit of work the dispatcher delivers - a webhook/job with a destination URL, a payload, an attempt count, and a delivery state. This is a plain data type with no logic beyond a state enum.
api:
  - type JobState string
  - const StatePending JobState = "pending"
  - const StateDelivered JobState = "delivered"
  - const StateDead JobState = "dead"
  - type Job struct { ID string; URL string; Payload []byte; Attempts int; State JobState }
  - func NewJob(id, url string, payload []byte) Job
behavior:
  - "NewJob returns a Job with the given id/url/payload, Attempts == 0, and State == StatePending."
  - "JobState is a string enum with exactly three values: pending, delivered, dead. A new job is pending; a successful delivery makes it delivered; exhausting all retries makes it dead."
  - "Job is a value type carrying data only - no methods beyond the constructor."
constraints: package main; standard library only
