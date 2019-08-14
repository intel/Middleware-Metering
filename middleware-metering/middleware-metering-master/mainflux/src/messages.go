package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	"time"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	client "github.com/influxdata/influxdb1-client/v2"
)

// json and int vars are to store the message counts and values 
// string vars are to query specific message names in InlfuxDB 
var dev1read1, dev1read2 json.Number
var cpu_usage, mem_used int
var address, dev1read1Name, dev1read2Name string

func query() {
	time.Sleep(5 * time.Second)

	//create a new InfluxDB client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("%s", address),
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	// query specific message counts for the first reading from InfluxDB
	param := map[string]interface{}{"name": dev1read1Name}
	q := client.NewQueryWithParameters("SELECT count(value) FROM messages WHERE \"name\" = $name", "mainflux", "", param)
	countIndex := 1
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			dev1read1 = `0`
		} else {
			tempVal := response.Results[0].Series[0].Values[0]
			dev1read1 = tempVal[countIndex].(json.Number)
		}
	}

	// query specific message counts for the second reading from InfluxDB
	param = map[string]interface{}{"name": dev1read2Name}
	q = client.NewQueryWithParameters("SELECT count(value) FROM messages WHERE \"name\"= $name", "mainflux", "", param)
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			dev1read2 = `0`
		} else {
			tempVal := response.Results[0].Series[0].Values[0]
		    dev1read2 = tempVal[countIndex].(json.Number)
		}
	}

	percent, _ := cpu.Percent(time.Second,false)
	cpu_usage = int(percent[0])

	vm, _ := mem.VirtualMemory()
	mem_used = int(vm.UsedPercent)
}

func main() {
	var PORT, dev1 string
	var dev1Msg uint64
	PORT = os.Getenv("PORT")
	dev1read1Name = os.Getenv("DEV1READ1")
	dev1read2Name = os.Getenv("DEV1READ2")
	dev1 = os.Getenv("DEVICE1")
	address = os.Getenv("ADDR")

	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		query()

		// calculate total device messages from the seperate readings 
		dev1read1Int, err1 := strconv.ParseUint(dev1read1.String(), 10, 64)
		dev1read2Int, err2 := strconv.ParseUint(dev1read2.String(), 10, 64)
		if err1 == nil && err2 == nil {
			dev1Msg = dev1read1Int + dev1read2Int
		}
		
		fmt.Fprintf(w, "{ \"%v Messages\": %v,\n", dev1read1Name, dev1read1)
		fmt.Fprintf(w, "\"%v Messages\": %v,\n", dev1read2Name, dev1read2)
		fmt.Fprintf(w, "\"%v Messages\": %v,\n\n", dev1, dev1Msg)
		fmt.Fprintf(w, "\"Percentage of CPU Used\": %v,\n", cpu_usage)
		fmt.Fprintf(w, "\"Percentage of Memory Used\": %v }", mem_used)
	})
	http.ListenAndServe(":" + PORT, nil)
}