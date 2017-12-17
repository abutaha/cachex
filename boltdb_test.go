package cachex

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenDBFile(t *testing.T) {
	f, err := TempFileString("boltdb-")
	if err != nil {
		t.Fatal(err)
	}

	db, err := GetBoltDB(f)
	if err != nil {
		t.Fatal(err)
	}

	db.Close()
	os.Remove(f)
}

func TestNewBoltCache(t *testing.T) {
	f, _ := TempFileString("boltdb-")
	db, _ := GetBoltDB(f)
	cache, err := NewBoltCache(db, "mybucket")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(cache.DB.Stats())

	db.Close()
	os.Remove(f)
}

func TestSet(t *testing.T) {
	cache, f := GetCacheInstance()

	err := cache.Set("key", "value")
	if err != nil {
		t.Fatal(err.Error())
	}

	cache.DB.Close()
	os.Remove(f)
}

func TestGet(t *testing.T) {
	cache, f := GetCacheInstance()
	cache.Set("key", "value")

	val, err := cache.Get("key")
	if val != "value" {
		t.Fatalf("value: %s, error: %s\n", val, err.Error())
	}

	cache.DB.Close()
	os.Remove(f)
}

func TestDelete(t *testing.T) {
	cache, f := GetCacheInstance()

	t.Logf("Creating key...\n")
	cache.Set("key", "value")

	t.Logf("Getting key...\n")
	val, err := cache.Get("key")
	if val != "value" {
		t.Fatalf("value: %s, error: %s\n", val, err.Error())
	}

	t.Logf("Deleting existing key...\n")
	err = cache.Delete("key")
	if err != nil {
		t.Fatalf("Failure deleting key. Error: %s\n", err.Error())
	}

	t.Logf("Getting key...\n")
	_, err = cache.Get("key")
	if err != nil {
		t.Logf("key deleted. %s", err)
	}

	t.Logf("Deleting non-existing key...\n")
	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Failure deleting key. Error: %s\n", err.Error())
	}

	cache.DB.Close()
	os.Remove(f)
}

func TestGetKeys(t *testing.T) {
	cache, f := GetCacheInstance()

	t.Logf("Creating keys...\n")
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")

	keys, err := cache.GetKeys()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(keys)

	cache.DB.Close()
	os.Remove(f)
}

func TestSearch(t *testing.T) {
	cache, f := GetCacheInstance()

	t.Logf("Creating keys...\n")
	cache.Set("key123asdf", "value1")
	cache.Set("keyfghkey456key", "value2")
	cache.Set("werver", "value3")
	cache.Set("louxcf", "value4")
	cache.Set("weryuip", "value5")

	keys, err := cache.Search("key")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(keys)

	cache.DB.Close()
	os.Remove(f)
}

func GetCacheInstance() (*BoltCache, string) {
	f, _ := TempFileString("boltdb-testsuite-")
	db, _ := GetBoltDB(f)
	cache, _ := NewBoltCache(db, "mybucket")
	return cache, f
}

// TempFileString - Returns a string for a temporary file
func TempFileString(s string) (string, error) {
	if s == "" {
		s = "gotempfile-"
	}

	// Get a tempfile
	file, err := ioutil.TempFile("", s)

	if err != nil {
		return "", err
	}

	// Close the tempfile
	if err := file.Close(); err != nil {
		return "", err
	}

	// Delete the file
	if err := os.Remove(file.Name()); err != nil {
		return "", err
	}

	return file.Name(), nil

}
