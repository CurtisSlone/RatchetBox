package main

import (
	"errors"
	"testing"
)

func TestRunLengthCodec(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		err  error
	}{
		{
			name: "empty string",
			in:   "",
			want: "",
			err:  nil,
		},
		{
			name: "single byte",
			in:   "a",
			want: "a",
			err:  nil,
		},
		{
			name: "long run",
			in:   "aaaaaaaaaa", // 10 a's
			want: string([]byte{10, 'a'}),
			err:  nil,
		},
		{
			name: "mixed string",
			in:   "aaab",
			want: string([]byte{3, 'a', 1, 'b'}),
			err:  nil,
		},
		{
			name: "digits",
			in:   "12",
			want: string([]byte{1, '1', 1, '2'}),
			err:  nil,
		},
		{
			name: "long run split",
			in:   "a" + "a", // 256 a's
			want: string([]byte{255, 'a', 1, 'a'}),
			err:  nil,
		},
		{
			name: "odd length input",
			in:   string([]byte{3, 'a', 2}), // odd length
			want: "",
			err:  errors.New(""),
		},
		{
			name: "zero count",
			in:   string([]byte{0, 'a'}),
			want: "",
			err:  errors.New(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.in)
			if tt.err == nil {
				decoded, err := Decode(encoded)
				if err != nil {
					t.Errorf("Decode(Encode(%q)) failed: %v", tt.in, err)
				}
				if decoded != tt.in {
					t.Errorf("Decode(Encode(%q)) = %q, want %q", tt.in, decoded, tt.in)
				}
			} else {
				_, err := Decode(tt.in)
				if err == nil {
					t.Errorf("Decode(%q) should have failed but didn't", tt.in)
				}
			}
		})
	}
}

func FuzzRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("a")
	f.Add("aa")
	f.Add("aaa")
	f.Add("aabbcc")
	f.Add("1234567890")
	f.Add("a" + "a") // 256 a's

	f.Fuzz(func(t *testing.T, input string) {
		encoded := Encode(input)
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Decode(Encode(%q)) failed: %v", input, err)
		}
		if decoded != input {
			t.Errorf("Decode(Encode(%q)) = %q, want %q", input, decoded, input)
		}
	})
}
