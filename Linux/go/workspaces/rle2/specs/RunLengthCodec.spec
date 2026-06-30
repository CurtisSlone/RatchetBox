name: RunLengthCodec
role: component
intent: A run-length codec using an UNAMBIGUOUS byte-pair scheme. Each maximal run of identical bytes is encoded as two bytes - a count byte followed by the data byte - so digits in the input are ordinary data bytes and can never be confused with counts. This makes Encode/Decode an exact, lossless round-trip for every input.
api:
  - func Encode(input string) string
  - func Decode(s string) (string, error)
behavior:
  - "ENCODING SCHEME (pinned, byte-pair): for each maximal run of an identical byte b that repeats n times where 1 <= n <= 255, emit exactly two bytes - byte(n) then b. A run longer than 255 is split into consecutive chunks each of length <= 255 (e.g. 300 of 'a' -> byte(255),'a', byte(45),'a')."
  - "Encode of the empty string is the empty string."
  - "Encode examples: Encode(\"aaa\") == string([]byte{3,'a'}); Encode(\"aaab\") == string([]byte{3,'a',1,'b'}); Encode(\"12\") == string([]byte{1,'1',1,'2'}) because each digit is a run of length 1 - the count byte (1) is distinct from the data byte ('1')."
  - "DECODING: read the input as consecutive (count, data) byte pairs; for each pair emit the data byte repeated count times."
  - "Decode returns an error if the input length is odd (a count byte with no following data byte), or if any count byte is 0 (an invalid run length)."
  - "ROUND-TRIP INVARIANT: Decode(Encode(s)) must equal s and a nil error, for every input string s (including strings containing digits, control bytes, and runs longer than 255)."
constraints: package: main
