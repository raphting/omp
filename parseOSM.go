package main

import "fmt"
import "reflect"
import "strconv"
import "sync/atomic"
import "os"
import "github.com/thomersch/gosmparse"
import "github.com/streamrail/concurrent-map"


type dataHandler struct {
    waterPoints cmap.ConcurrentMap
    nodePoints cmap.ConcurrentMap
    waterwayCounter uint64
}

func (d *dataHandler) ReadRelation(r gosmparse.Relation) {}
func (d *dataHandler) ReadNode(n gosmparse.Node) {
    latlon := [2]float32{n.Lat, n.Lon}
    d.nodePoints.Set(strconv.FormatInt(n.ID, 10), latlon)
}

func (d *dataHandler) ReadWay(w gosmparse.Way) {
    val, ok := w.Tags["waterway"]
    if ok && (val == "river" || val == "riverbank" || val == "stream") {
        incCtr := atomic.AddUint64(&d.waterwayCounter, 1)
        d.waterPoints.Set(strconv.FormatUint(incCtr, 10), w.NodeIDs)
    }
}

func main() {
    r, err := os.Open("hamburg-latest.osm.pbf")
    if err != nil {
        panic(err)
    }
    dec := gosmparse.NewDecoder(r)
    //Parse data to find water
    d := dataHandler{waterwayCounter: 0, waterPoints: cmap.New(), nodePoints: cmap.New()}
    err = dec.Parse(&d)
    if err != nil {
        panic(err)
    }

    fmt.Println(d.waterPoints.Count())

    for i := 1; i < d.waterPoints.Count(); i++ {
        nodes, _ := d.waterPoints.Get(strconv.Itoa(i))
        r := reflect.ValueOf(nodes)
        for node := 0; node < r.Len(); node++ {
            latlon, _ := d.nodePoints.Get(strconv.FormatInt(r.Index(node).Int(), 10))
            fmt.Println(latlon)
        }

        fmt.Println("----------------------------")
    }
}
