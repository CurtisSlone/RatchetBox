package main

import (
	"testing"
)

func TestRunLengthCodec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		encoded  string
		decoded  string
		wantErr  bool
	}{
		{
			name:    "empty string",
			input:   "",
			encoded: "",
			decoded: "",
		},
		{
			name:    "single character",
			input:   "a",
			encoded: "a",
			decoded: "a",
		},
		{
			name:    "two different characters",
			input:   "ab",
			encoded: "ab",
			decoded: "ab",
		},
		{
			name:    "repeated character",
			input:   "aaa",
			encoded: "3a",
			decoded: "aaa",
		},
		{
			name:    "mixed run and non-run",
			input:   "aabcccccaaa",
			encoded: "2ab5c3a",
			decoded: "aabcccccaaa",
		},
		{
			name:    "string with digits",
			input:   "123",
			encoded: "123",
			decoded: "123",
		},
		{
			name:    "run with digit",
			input:   "111",
			encoded: "31",
			decoded: "111",
		},
		{
			name:    "long run",
			input:   "aaaaaaaaaaaaaaaaaaaaaaaaaa",
			encoded: "26a",
			decoded: "aaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		{
			name:    "invalid encoding - digit not followed by character",
			input:   "1",
			encoded: "1",
			wantErr: true,
		},
		{
			name:    "invalid encoding - invalid digit",
			input:   "a2b",
			encoded: "a2b",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.input)
			if encoded != tt.encoded {
				t.Errorf("Encode(%q) = %q, want %q", tt.input, encoded, tt.encoded)
			}

			decoded, err := Decode(tt.encoded)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Decode(%q) should have returned an error", tt.encoded)
				}
				return
			}

			if err != nil {
				t.Errorf("Decode(%q) returned error: %v", tt.encoded, err)
				return
			}

			if decoded != tt.decoded {
				t.Errorf("Decode(%q) = %q, want %q", tt.encoded, decoded, tt.decoded)
			}
		})
	}
}

func FuzzRunLengthCodecRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("a")
	f.Add("aa")
	f.Add("aaa")
	f.Add("a1b2c3")
	f.Add("1234567890")
	f.Add("aabbcc")
	f.Add("aaaaaa")
	f.Add("a1b2c3d4e5f6g7h8i9j0")

	f.Fuzz(func(t *testing.T, input string) {
		encoded := Encode(input)
		decoded, err := Decode(encoded)
		if err != nil {
			t.Fatalf("Decode failed on encoded input: %v", err)
		}
		if decoded != input {
			t.Errorf("Round-trip failed: input %q, got %q", input, decoded)
		}
	})
}

func FuzzRunLengthCodecDecodeEncode(f *testing.F) {
	f.Add("")
	f.Add("a")
	f.Add("aa")
	f.Add("3a")
	f.Add("2ab5c3a")
	f.Add("1234567890")

	f.Fuzz(func(t *testing.T, s string) {
		decoded, err := Decode(s)
		if err != nil {
			// If Decode fails, we don't expect Encode to succeed
			return
		}
		encoded := Encode(decoded)
		if encoded != s {
			t.Errorf("Decode(Encode(%q)) != %q", s, s)
		}
	})
}
