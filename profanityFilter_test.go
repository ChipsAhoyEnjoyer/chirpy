package main

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
		out := profanityFilter(test.input)
		if out != test.expected {
			fmt.Println("Test failed:")
			fmt.Printf("Expected:    %v\n", test.expected)
			fmt.Printf("Actual:      %v\n", out)
			t.Fail()
		}
	}
}

func TestToTitle(t *testing.T) {
	tests := []testCase{
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
		out := ToTitle(test.input)
		if out != test.expected {
			fmt.Println("Test failed:")
			fmt.Printf("Expected:    %v\n", test.expected)
			fmt.Printf("Actual:      %v\n", out)
			t.Fail()
		}
	}
}
