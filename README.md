# cachex

**cachex** aims to simplify usage of different caching/storage layers by using an interface that should be implemented by each cache engine.
At the moment, it supports only [boltdb](https://github.com/boltdb/bolt).

## Usage

```go
package main

import (
        "github.com/abutaha/cachex"
)

func main() {
        db, _ := cachex.GetBoltDB("bolt.db")
        cache, _ := cachex.NewBoltCache(db, "mybucket")
        err := cache.Set("mykey", "myvalue")
        val, err := cache.Get("mykey")
        err = cache.Delete("mykey")
        foundKeys, err := cache.Search("key")
        keys, err := cache.GetKeys()
}
```

