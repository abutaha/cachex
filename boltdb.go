package cachex

import (
	"bytes"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// BoltCache - container to implement Cache interface
type BoltCache struct {
	DB     *bolt.DB
	Bucket []byte
}

// GetBoltDB - returns open BoltDB database with read/write permissions or
// goes down in flames if something bad happends
func GetBoltDB(dbfile string) (*bolt.DB, error) {
	// Check if file exist
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return db, nil
}

// NewBoltCache - returns a new BoltCache instance and create bucket if not exist
func NewBoltCache(db *bolt.DB, bucket string) (*BoltCache, error) {
	c := &BoltCache{
		DB:     db,
		Bucket: []byte(bucket),
	}

	err := c.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(c.Bucket)
		if err != nil {
			return err
		}
		return nil
	})

	return c, err
}

// Set - save given key/value pair to cache. Returns nil for success
func (c *BoltCache) Set(key, value string) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.Bucket)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", c.Bucket)
		}
		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get - Return value of a given key if found
func (c *BoltCache) Get(key string) (value string, err error) {
	err = c.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.Bucket)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", c.Bucket)
		}
		// "Byte slices returned from Bolt are only valid during a transaction."
		var buffer bytes.Buffer
		val := bucket.Get([]byte(key))

		// If it doesn't exist then it will return nil
		if val == nil {
			return fmt.Errorf("key %q not found", key)
		}

		buffer.Write(val)
		value = string(buffer.Bytes())
		return nil
	})
	return
}

// Search - Search for a given prefix pattern
func (c *BoltCache) Search(prefix string) (vals map[string]string, err error) {
	err = c.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(c.Bucket).Cursor()
		vals = make(map[string]string)
		for k, v := c.Seek([]byte(prefix)); bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
			vals[string(k)] = string(v)
		}
		return nil
	})
	return
}

// Delete - deletes specified key
func (c *BoltCache) Delete(key string) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.Bucket)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", c.Bucket)
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetKeys - gets all current keys
func (c *BoltCache) GetKeys() (keys map[string]bool, err error) {
	err = c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.Bucket)

		keys = make(map[string]bool)

		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys[string(k)] = true
		}
		return nil
	})
	return
}
