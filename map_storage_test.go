package starrpg

import (
	"testing"
)

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
