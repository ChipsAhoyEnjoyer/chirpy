package titlecase

import (
	"fmt"
	"testing"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "foo",
			expected: "Foo",
		},
		{
			input:    "FoO",
			expected: "Foo",
		},
		{
			input:    "fOO bar",
			expected: "Foo Bar",
		},
		{
			input:    "0foo 0Bar",
			expected: "0foo 0Bar",
		},
		{
			input:    "fo00oo b4r",
			expected: "Fo00oo B4r",
		},
	}
	for _, test := range tests {
		out := Titlecase(test.input)
		if out != test.expected {
			fmt.Println("Test failed:")
			fmt.Printf("Expected:    %v\n", test.expected)
			fmt.Printf("Actual:      %v\n", out)
			t.Fail()
		}
	}
}
