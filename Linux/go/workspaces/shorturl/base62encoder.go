package main

type Base62Encoder struct {
	chars string
}

func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{
		chars: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
}

func (e *Base62Encoder) Encode(n int64) string {
	if n == 0 {
		return "0"
	}

	if n < 0 {
		// For negative numbers, we'll encode the absolute value and prepend a minus sign
		return "-" + e.Encode(-n)
	}

	var result string
	for n > 0 {
		remainder := n % 62
		result = string(e.chars[remainder]) + result
		n /= 62
	}

	return result
}
