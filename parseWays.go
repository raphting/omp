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

func valToFloats(val []byte) []float64 {
	var cnt uint8 = 0
	var latRaw uint64 = 0
	var lonRaw uint64 = 0
	for cnt < 8 {
		latRaw = latRaw | uint64(val[cnt])<<(cnt*8)
		cnt++
	}

	for cnt < 16 {
		lonRaw = lonRaw | uint64(val[cnt])<<((cnt-8)*8)
		cnt++
	}

	res := make([]float64, 2)
	res[0] = math.Float64frombits(latRaw)
	res[1] = math.Float64frombits(lonRaw)
	return res
}

func main() {
	fmt.Println("Reading Node file")

	wayFile, err := os.Open("./waters")
	check(err)
	defer wayFile.Close()

	db, err := bolt.Open("second.db", 0600, nil)
	check(err)
	defer db.Close()


	scanner := bufio.NewScanner(wayFile)
	check(err)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")

		for id := range data {
			idnum, _ := strconv.ParseUint(data[id], 10, 64)
			key := intToByte(idnum)
			fmt.Println("ID: ", idnum)
			db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("latlon"))
				val := bucket.Get(key)
				if len(val) == 16 {
					fmt.Println(valToFloats(val))
				}
				return nil
			})
		}
	}
	fmt.Println("End.")
}
