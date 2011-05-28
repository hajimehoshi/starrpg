package starrpg

import (
	//"strconv"
	"testing"
)

func TestCheckAcceptHeader(t *testing.T) {
	type TestCase struct {
		mediaType string
		accept string
		expectedResult float64
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
		actual != testCase.expectedResult {
			t.Errorf("checkAcceptHeader(%#v, %#v) is not %#v but %#v",
				testCase.mediaType,
				testCase.accept,
				testCase.expectedResult,
				actual)
		}
	}
}

func TestIsPostablePath(t *testing.T) {
	type TestCase struct {
		path string
		expectedResult bool
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
	for _, testCase := range testCases {
		if actual := isPostablePath(testCase.path); actual != testCase.expectedResult {
			t.Errorf("isGettablePath(%#v) is not %#v but %#v",
			testCase.path,
			testCase.expectedResult,
			actual)
		}
	}
}

func TestIsPuttablePath(t *testing.T) {
	type TestCase struct {
		path string
		expectedResult bool
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
		{"/games/1/foos/2", false},
		{"", false},
		{"games", false},
	}
	for _, testCase := range testCases {
		if actual := isPuttablePath(testCase.path); actual != testCase.expectedResult {
			t.Errorf("isGettablePath(%#v) is not %#v but %#v",
			testCase.path,
			testCase.expectedResult,
			actual)
		}
	}
}

func TestIsDeletablePath(t *testing.T) {
	type TestCase struct {
		path string
		expectedResult bool
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
		{"/games/1/foos/2", false},
		{"", false},
		{"games", false},
	}
	for _, testCase := range testCases {
		if actual := isDeletablePath(testCase.path); actual != testCase.expectedResult {
			t.Errorf("isGettablePath(%#v) is not %#v but %#v",
			testCase.path,
			testCase.expectedResult,
			actual)
		}
	}
}

func TestDoPost(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	newPath, err := doPost(mapStorage, "/foos")
	if err != nil {
		t.Errorf(`doPost(storage, "/foos") failed: %s`, err.String())
	}
	if newPath != "/foos/1" {
		t.Errorf(`newPath is not "/foos/1" but %#v`, newPath)
	}
	mapStorage.Get("/foos")
	
	/*innerCount, err := strconv.Atoui64(string(innerCountBytes))
	if err != nil {
		t.Errorf(`strconv.Atoui64(string(innerCountBytes)) failed: %s`, err.String())
	}
	if innerCount != 1 {
		t.Errorf(`innerCount is not 1 but %#v`, innerCount)
	}*/
	/*newCollectionBytes := (*storage)["/foos"]
	var newCollection map[string]map[string]string
	if err := json.Unmarshal(newCollectionBytes, &newCollection); err != nil {
		t.Errorf(`json.Unmarshal(newValueBytes, newValue) failed: %s`, err.String())
	}
	newValue := newCollection["1"]
	name := newValue["name"]
	if name != "" {
		t.Errorf(`name is not "" but %#v`,  name)
	}
	newItemBytes := (*storage)["/foos/1"]
	if string(newItemBytes) != "{}" {
		t.Errorf(`newItemBytes is not %#v but %#v`, "{}", newItemBytes)
	}*/
}

/*func TestJson(t *testing.T) {
	data := map[string]map[string]string {
		"1": {"name": "ii"},
		"2": {"name": "fdajslkfjkl"},
	}
	json, err := json.Marshal(data)
	fmt.Println(string(json), err)
}*/
