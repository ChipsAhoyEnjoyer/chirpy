package utils

import (
	"fmt"
	"testing"
)

type testCase struct {
	input    string
	expected string
}

func TestProfanityFilter(t *testing.T) {
	tests := []testCase{
		{
			input:    " kerfuffle sharbert fornax ",
			expected: " **** **** **** ",
		},
		{
			input:    "What a Kerfuffle !",
			expected: "What a **** !",
		},
		{
			input:    " SHARBERT ",
			expected: " **** ",
		},
		{
			input:    " fornax ",
			expected: " **** ",
		},
		{
			input:    "KERFUFFLE, SharBert, fornax",
			expected: "KERFUFFLE, SharBert, fornax",
		},
	}
	for _, test := range tests {
		out := ProfanityFilter(test.input)
		if out != test.expected {
			fmt.Println("Test failed:")
			fmt.Printf("Expected:    %v\n", test.expected)
			fmt.Printf("Actual:      %v\n", out)
			t.Fail()
		}
	}
}
