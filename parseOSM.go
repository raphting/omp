package main

import "fmt"
import "strings"
import "os"
import "github.com/thomersch/gosmparse"

type dataHandler struct {
    nodeChan chan gosmparse.Node
    waterChan chan []int64
}

func (d *dataHandler) ReadRelation(r gosmparse.Relation) {}
func (d *dataHandler) ReadNode(n gosmparse.Node) {
    d.nodeChan <- n
}

func (d *dataHandler) ReadWay(w gosmparse.Way) {
    val, ok := w.Tags["waterway"]
    if ok && (val == "river" || val == "riverbank" || val == "stream") {
        d.waterChan <- w.NodeIDs
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func StoreNode(nodeChan <-chan gosmparse.Node, nodeFile *os.File) {
    for {
        node := <-nodeChan
        res := []string{fmt.Sprintf("%v", node.ID), fmt.Sprintf("%f", node.Lat), fmt.Sprintf("%f", node.Lon)}
        nodeFile.WriteString(strings.Join(res, ","))
        nodeFile.WriteString("\n")
    }
}

func StoreWater(waterChan <-chan []int64, waterFile *os.File) {
    for {
        ids := <-waterChan
        for cnt := range ids[:len(ids)-1] {
            waterFile.WriteString(fmt.Sprintf("%v", ids[cnt]))
            waterFile.WriteString(",")
        }

        waterFile.WriteString(fmt.Sprintf("%v", ids[len(ids)-1]))
        waterFile.WriteString("\n")
    }
}

func main() {
    //Open original file and prepare Decoder
    r, err := os.Open("hamburg-latest.osm.pbf")
    check(err)
    defer r.Close()

    nodeFile, err := os.Create("./nodes")
    check(err)
    defer nodeFile.Close()
    waterFile, err := os.Create("./waters")
    check(err)
    defer waterFile.Close()

    nodeChan := make(chan gosmparse.Node, 8)
    waterChan := make(chan []int64, 8)
    go StoreNode(nodeChan, nodeFile)
    go StoreWater(waterChan, waterFile)

    dec := gosmparse.NewDecoder(r)
    err = dec.Parse(&dataHandler{nodeChan: nodeChan, waterChan: waterChan})
    check(err)

    //Wait for all channels beeing written to files
    for {
        if len(nodeChan) == 0 && len(waterChan) == 0 {
            break
        }
    }

    nodeFile.Sync()
    waterFile.Sync()
    fmt.Println("End.")
}
