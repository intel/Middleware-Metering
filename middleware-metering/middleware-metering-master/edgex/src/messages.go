package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"time"
	"strconv"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)
// int variables are initiated for message count types
// string vars are initiated to use in queries
var reading1Count, reading2Count, deviceEvents, totalDeviceCounts, cpu_usage, mem_used int
var reading_request1, reading_request2, event_request string

func query() {
	var msg_size int
	time.Sleep(5 * time.Second)
	//number of characters in a single message
	msg_size = 151
	// sends GET request for the reading1 message metadata object 
	resp, err := http.Get(fmt.Sprintf("%s", reading_request1))
	if err != nil {
		reading1Count = 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reading1Count = 0
	}
	// finds the number of individual messages that have been sent for reading1
	reading1Count = (len(body) / msg_size)

    // sends GET request for the reading2 message metadata object 
	resp, err = http.Get(fmt.Sprintf("%s", reading_request2))
	if err != nil {
		reading2Count = 0
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		reading2Count = 0
	}
	// finds the number of individual messages that have been sent for reading2
	reading2Count = (len(body) / msg_size) 

	// sends GET request for the message events from specified device
	resp, err = http.Get(fmt.Sprintf("%s", event_request))
	if err != nil {
		deviceEvents = 0
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		deviceEvents = 0
	}
	temp := string(body)
	deviceEvents, err = strconv.Atoi(temp)
	if err != nil {
		deviceEvents = 0
	}

	totalDeviceCounts = reading1Count + reading2Count

	percent, _ := cpu.Percent(time.Second,false)
	cpu_usage = int(percent[0])

	vm, _ := mem.VirtualMemory()
	mem_used = int(vm.UsedPercent)
}

func main() {
	var PORT, reading1, reading2, device string
	reading1 = os.Getenv("READING1")
	reading2 = os.Getenv("READING2")
	device = os.Getenv("DEVICE")
	PORT = os.Getenv("PORT")
	reading_request1 = os.Getenv("READING_REQUEST1")
	reading_request2 = os.Getenv("READING_REQUEST2")
	event_request = os.Getenv("EVENT_REQUEST")

	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		query()
		fmt.Fprintf(w, "{ \"%v Messages\": %v,\n",reading1, reading1Count)
		fmt.Fprintf(w, "\"%v Messages\": %v,\n", reading2, reading2Count)
		fmt.Fprintf(w, "\"%v Messages\": %v,\n", device, totalDeviceCounts)
		fmt.Fprintf(w, "\"%v Message Events\": %v, \n\n", device, deviceEvents)
		fmt.Fprintf(w, "\"Percentage of CPU Used\": %v,\n", cpu_usage)
		fmt.Fprintf(w, "\"Percentage of Memory Used\": %v }", mem_used)
	})
	http.ListenAndServe(":" + PORT, nil)
}