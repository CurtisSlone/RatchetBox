name: Base62Encoder
role: component
intent: Encodes integers to base62 strings for short URL codes
api:
  - func NewBase62Encoder() *Base62Encoder
  - func (e *Base62Encoder) Encode(n int64) string
behavior:
  - Encode should produce base62 strings using characters [0-9a-zA-Z]
  - Should handle zero and negative numbers
  - Should produce consistent output for same input
constraints: Use math/big for large numbers if needed; package: main
