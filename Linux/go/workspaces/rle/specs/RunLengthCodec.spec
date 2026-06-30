name: RunLengthCodec
role: component
intent: A run-length encoding codec that compresses runs of the same byte and decompresses them, handling all edge cases including digits, long runs, and empty strings
api:
  - func Encode(input string) string
  - func Decode(s string) (string, error)
behavior:
  - Encode should compress runs of the same byte
  - Decode should reverse the encoding and return an error for invalid input
  - Decode(Encode(s)) should equal s for any input string
constraints: package: main
