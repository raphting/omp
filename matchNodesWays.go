package main

import "bufio"
import "fmt"
import "math"
import "strconv"
import "strings"
import "os"
import "github.com/boltdb/bolt"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func intToByte(myInt uint64) []byte {
	res := make([]byte, 8)
	var cnt uint
	var x byte

	for cnt = 0; cnt < 8; cnt++ {
		x = byte((myInt>>(cnt*8)) & 0xff)
		res[cnt] = x
	}

	return res
}

func main() {
	fmt.Println("Reading Node file")

	nodeFile, err := os.Open("./nodes.sorted")
	check(err)
	defer nodeFile.Close()

	db, err := bolt.Open("second.db", 0600, nil)
	check(err)
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("latlon"))
		return err
	})
	check(err)


	scanner := bufio.NewScanner(nodeFile)
	tx, err := db.Begin(true)
	check(err)
	bucket := tx.Bucket([]byte("latlon"))
	var progress uint64 = 0
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		id, _ := strconv.ParseUint(data[0], 10, 64)
		lat, _ := strconv.ParseFloat(data[1], 32)
		lon, _ := strconv.ParseFloat(data[2], 32)

		key := intToByte(id)
		val1 := intToByte(math.Float64bits(lat))
		val2 := intToByte(math.Float64bits(lon))
		value := append(val1, val2...)

		err := bucket.Put(key, value)
		check(err)

		progress++
		if progress % 100000 == 0 {
			tx.Commit()
			tx, err = db.Begin(true)
			check(err)
			bucket = tx.Bucket([]byte("latlon"))
			fmt.Print("#")
		}
	}
	tx.Commit()
	fmt.Println()
	fmt.Println("End.")
}
