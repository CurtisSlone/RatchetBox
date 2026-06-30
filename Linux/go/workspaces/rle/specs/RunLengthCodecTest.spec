name: RunLengthCodecTest
role: test
intent: Test the run-length encoding codec with example cases and a fuzz target for round-trip property
api:
  - func TestRunLengthCodec(t *testing.T)
behavior:
  - Test with example cases including empty string, single characters, repeated characters, and strings with digits
  - Test round-trip property with fuzz target for arbitrary inputs
  - Assert that Decode(Encode(s)) == s for all test cases
constraints: package: main
