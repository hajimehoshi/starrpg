package starrpg

import (
	"testing"
)

func TestCheckAcceptHeader(t *testing.T) {
	type TestCase struct {
		mediaType string
		accept string
		expected float64
	}
	testCases := []TestCase{
		{"application/xml",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			1.0},
		{"application/xhtml+xml",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			1.0},
		{"text/html",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			0.9},
		{"text/plain",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			0.8},
		{"image/png",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			1.0},
		{"text/xml",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			0.7},
		{"application/json",
			"application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,text/*;q=0.7,*/*;q=0.5",
			0.5},
		{"text/html",
			"",
			0},
		{"text/html",
			",,,;;;",
			0},
		{"text/html",
			"text/html;q=foo",
			0},
		{"text/html",
			"text/html;q=NaN",
			0},
		{"application/xhtml+xml",
			"application/xhtml+xml,*/html",
			0},
		{"text/html",
			"application/xhtml+xml,*/html",
			0},
	}
	for _, testCase := range testCases {
		if actual := checkAcceptHeader(testCase.mediaType, testCase.accept);
		actual != testCase.expected {
			t.Errorf("checkAcceptHeader(%#v, %#v) is not %#v but %#v",
				testCase.mediaType,
				testCase.accept,
				testCase.expected,
				actual)
		}
	}
}

/*func TestDoPost(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	resourceStorage := NewResourceStorage(mapStorage)
	newPath, err := doPost(resourceStorage, "/foos")
	if err != nil {
		t.Errorf(`doPost(storage, "/foos") failed: %s`, err.String())
	}
	if newPath != "/foos/1" {
		t.Errorf(`newPath is not "/foos/1" but %#v`, newPath)
	}
}*/
