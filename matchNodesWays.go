package main

import "bufio"
import "fmt"
import "math"
import "strconv"
import "strings"
import "os"

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

func main() {
	fmt.Println("Reading Node file")

	nodeFile, err := os.Open("./nodes")
	check(err)
	defer nodeFile.Close()

	scanner := bufio.NewScanner(nodeFile)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		id, _ := strconv.ParseUint(data[0], 10, 64)
		lat, _ := strconv.ParseFloat(data[1], 32)
		lon, _ := strconv.ParseFloat(data[2], 32)

		key := intToByte(id)
		val1 := intToByte(math.Float64bits(lat))

		val2 := intToByte(math.Float64bits(lon))
		value := append(val1, val2...)

		fmt.Println(key, value)
	}

}
