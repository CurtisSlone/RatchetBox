name: RunLengthCodecTest
role: test
intent: Test the byte-pair run-length codec against the pinned scheme, leaning on the self-validating round-trip property plus a few canonical, hand-verified literal examples.
api:
  - func TestRunLengthCodec(t *testing.T)
behavior:
  - "Primary: assert the ROUND-TRIP property Decode(Encode(s)) == s (and nil error) for a table of inputs - empty string, a single byte, a long run, a mixed string, and a string of digits like \"112233\". This property is self-validating and cannot encode a wrong expectation."
  - "Pinned literal checks (hand-verified against the scheme): Encode(\"aaa\") == string([]byte{3,'a'}); Encode(\"\") == \"\"; Encode(\"12\") == string([]byte{1,'1',1,'2'})."
  - "Error checks: Decode of an odd-length input (e.g. string([]byte{3,'a',2})) returns a non-nil error; Decode of an input with a zero count byte (e.g. string([]byte{0,'a'})) returns a non-nil error."
  - "Include a fuzz target FuzzRoundTrip(f) seeded with f.Add on several strings, asserting Decode(Encode(s)) == s for arbitrary string inputs."
  - "Do NOT invent textual count-prefix expectations (like \"3a\") - the scheme is binary byte pairs, not decimal-digit prefixes."
constraints: package: main
