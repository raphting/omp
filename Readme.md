One month project
=================

What is this all about?
-----------------------

Please check my [weblog](https://raphaelsprenger.de/blog/) to read about it.

Dependencies
------------

1. OSM data with the source of your choice
1. `go get "github.com/thomersch/gosmparse"`
1. `go get "https://github.com/boltdb/bolt"`

Build and run
-------------

### First Step
Parse an OSM file to get all the Nodes and all the data which contains naturally formed water.
This will create two files "nodes" and "waters" with a lot of data in it:

`go run parseOSM.go`

### Second Step
Create a Key-Value DB consisting of all Node-IDs as keys and the corresponding latitude and longitude in a custom format as value.

`go run match matchNodesWays.go`

### Third Step
Use the database to query the "waters" file for latitude and longitude. Write a new file. Each line contains the coordinates of one way in the format: lat,lon,lat,lon,...

### Visualize
After the third step you can use the render.html file to visualize the data in the webbrowser. See my blog (link above) for images.
