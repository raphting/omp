package main

import "fmt"
import "github.com/boltdb/bolt"

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("WaterBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte("answer"), []byte("42"))

		v := b.Get([]byte("answer"))
		fmt.Printf("%s\n", v)
		return err
	})


}
