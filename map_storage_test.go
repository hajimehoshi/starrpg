package starrpg

import (
	"testing"
)

func TestMapStorageGetPrefix(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	mapStorage.Set("/foo", map[string]string{"value": "value"})
	mapStorage.Set("/foo/1", map[string]string{"value": "value-1"})
	mapStorage.Set("/bar/1", map[string]string{"value": "value-1"})
	mapStorage.Set("/foo/2", map[string]string{"value": "value-2"})
	mapStorage.Set("/foo/abcde", map[string]string{"value": "value-abcde"})
	objs, err := mapStorage.GetWithPrefix("/foo/")
	if err != nil {
		t.Error(err)
	}
	if objs == nil {
		t.Error(`mapStorage.GetWithPrefix("/foo") returns nil`)
	}
	if expected, actual := 3, len(objs); expected != actual {
		t.Errorf(`len(obj) is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-1", objs["/foo/1"]["value"]; expected != actual {
		t.Errorf(`objs["/foo/1"]["value"] is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-2", objs["/foo/2"]["value"]; expected != actual {
		t.Errorf(`objs["/foo/2"]["value"] is not %#v but %#v`, expected, actual)
	}
	if expected, actual := "value-abcde", objs["/foo/abcde"]["value"]; expected != actual {
		t.Errorf(`objs["/foo/abcde"]["value"] is not %#v but %#v`, expected, actual)
	}
}

func TestMapStorageSet(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	{
		mapStorage.Set("foo", map[string]string{"bar": "baz"})
		expected := "baz"
		m, err := mapStorage.Get("foo")
		if err != nil {
			t.Error(err)
		}
		actual := m["bar"]
		if expected != actual {
			t.Errorf(`mapStorage.Get("foo")["bar"] is not %#v but %#v`, expected, actual)
		}
	}
}

func TestMapStorageInc(t *testing.T) {
	storage := &DummyStorage{}
	mapStorage := NewMapStorage(storage)
	{
		expected := uint64(1)
		actual, err := mapStorage.Inc("foo", "count")
		if err != nil {
			t.Error(err)
		}
		if expected != actual {
			t.Errorf(`mapStorage.Inc("foo", "count") is not %#v but %#v`, expected, actual)
		}
	}
	{
		expected := uint64(2)
		actual, err := mapStorage.Inc("foo", "count")
		if err != nil {
			t.Error(err)
		}
		if expected != actual {
			t.Errorf(`mapStorage.Inc("foo", "count") is not %#v but %#v`, expected, actual)
		}
	}
}
