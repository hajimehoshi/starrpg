package starrpg

import (
	"testing"
)

func TestIsPostablePath(t *testing.T) {
	type TestCase struct {
		path string
		expected bool
	}
	testCases := []TestCase{
		{"/games", true},
		{"/games/1", false},
		{"/games/1/maps", false},
		{"/games/1/maps/2", false},
		{"/games/", false},
		{"/games/1/", false},
		{"/games/1/maps/", false},
		{"/games/1/maps/2/", false},
		{"/foos/", false},
		{"/foos/1", false},
		{"/games/1/foos", false},
		{"/games/1/foos/2", false},
		{"", false},
		{"games", false},
	}
	r := &resourceRequestProcessor{nil}
	for _, testCase := range testCases {
		if actual := r.isPostablePath(testCase.path); actual != testCase.expected {
			t.Errorf("r.isPostablePath(%#v) is not %#v but %#v",
			testCase.path,
			testCase.expected,
			actual)
		}
	}
}

func TestIsPuttablePath(t *testing.T) {
	type TestCase struct {
		path string
		expected bool
	}
	testCases := []TestCase{
		{"/games", false},
		{"/games/1", true},
		{"/games/1/maps", false},
		{"/games/1/maps/2", true},
		{"/games/", false},
		{"/games/1/", false},
		{"/games/1/maps/", false},
		{"/games/1/maps/2/", false},
		{"/foos/", false},
		{"/foos/1", false},
		{"/games/1/foos", false},
		{"/games/1/foos/2", true},
		{"", false},
		{"games", false},
	}
	r := &resourceRequestProcessor{nil}
	for _, testCase := range testCases {
		if actual := r.isPuttablePath(testCase.path); actual != testCase.expected {
			t.Errorf("r.isPuttablePath(%#v) is not %#v but %#v",
			testCase.path,
			testCase.expected,
			actual)
		}
	}
}

func TestIsDeletablePath(t *testing.T) {
	type TestCase struct {
		path string
		expected bool
	}
	testCases := []TestCase{
		{"/games", false},
		{"/games/1", true},
		{"/games/1/maps", false},
		{"/games/1/maps/2", false},
		{"/games/", false},
		{"/games/1/", false},
		{"/games/1/maps/", false},
		{"/games/1/maps/2/", false},
		{"/foos/", false},
		{"/foos/1", false},
		{"/games/1/foos", false},
		{"/games/1/foos/2", false},
		{"", false},
		{"games", false},
	}
	r := &resourceRequestProcessor{nil}
	for _, testCase := range testCases {
		if actual := r.isDeletablePath(testCase.path); actual != testCase.expected {
			t.Errorf("r.isDeletablePath(%#v) is not %#v but %#v",
				testCase.path,
				testCase.expected,
				actual)
		}
	}
}
