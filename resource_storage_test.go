package starrpg

import (
	"testing"
)

func TestResourceStorageURLPathToStoragePath(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	resourceStorage := &resourceStorageImpl{mapStorage}
	type TestCase struct {
		urlPath string
		expected string
	}
	testCases := []TestCase {
		{"", ""},
		{"/", ""},
		{"model1", ""},
		{"/model1", "1:/model1"},
		{"/model1/id1", "2:/model1/id1"},
		{"/model1/id1/model2", "3:/model1/id1/model2"},
		{"/model1/id1/model2/id2", "4:/model1/id1/model2/id2"},
		{"/model1/id1/model2/id2/model3", "5:/model1/id1/model2/id2/model3"},
		{"/model1/id1/model2/id2/model3/id3", "6:/model1/id1/model2/id2/model3/id3"},
	}
	for _, testCase := range testCases {
		if actual := resourceStorage.urlPathToStoragePath(testCase.urlPath);
		actual != testCase.expected {
			t.Errorf(`urlPathToStoragePath(%#v) is not %#v but %#v`,
				testCase.urlPath,
				testCase.expected,
				actual)
		}
	}
}

func TestResourceStorageURLPathToStorageChildrenPathPrefix(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	resourceStorage := &resourceStorageImpl{mapStorage}
	type TestCase struct {
		urlPath string
		expected string
	}
	testCases := []TestCase {
		{"", ""},
		{"/", "1:/"},
		{"model1", ""},
		{"/model1", "2:/model1/"},
		{"/model1/id1", "3:/model1/id1/"},
		{"/model1/id1/model2", "4:/model1/id1/model2/"},
		{"/model1/id1/model2/id2", "5:/model1/id1/model2/id2/"},
		{"/model1/id1/model2/id2/model3", "6:/model1/id1/model2/id2/model3/"},
		{"/model1/id1/model2/id2/model3/id3", "7:/model1/id1/model2/id2/model3/id3/"},
	}
	for _, testCase := range testCases {
		if actual := resourceStorage.urlPathToStorageChildrenPathPrefix(testCase.urlPath);
		actual != testCase.expected {
			t.Errorf(`urlPathToStoragePath(%#v) is not %#v but %#v`,
				testCase.urlPath,
				testCase.expected,
				actual)
		}
	}
}

func TestResourceStorageGetChildren(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	resourceStorage := NewResourceStorage(mapStorage)
	resourceStorage.Set("/foo", map[string]string{"value": "value"})
	resourceStorage.Set("/foo/1", map[string]string{"value": "value-1"})
	resourceStorage.Set("/bar/1", map[string]string{"value": "value-1"})
	resourceStorage.Set("/foo/2", map[string]string{"value": "value-2"})
	resourceStorage.Set("/foo/2/bar", map[string]string{"value": "value-2-bar"})
	resourceStorage.Set("/foo/abcde", map[string]string{"value": "value-abcde"})
	objs, err := resourceStorage.GetChildren("/foo")
	if err != nil {
		t.Fatal(err)
	}
	if objs == nil {
		t.Error(`resourceStorage.GetChildren("/foo") returns nil`)
	}
	if expected, actual := 3, len(objs); expected != actual {
		t.Errorf(`len(obj) is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-1", objs["1"]["value"]; expected != actual {
		t.Errorf(`objs["1"]["value"] is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-2", objs["2"]["value"]; expected != actual {
		t.Errorf(`objs["2"]["value"] is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-abcde", objs["abcde"]["value"]; expected != actual {
		t.Errorf(`objs["abcde"]["value"] is not %#v but %#v`, expected, actual)
	}
}
