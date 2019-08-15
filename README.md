# Mainflux & EdgeX Metering Microservices 

Katie Nguyen   
Intel Corporation 

**Description:** The Mainflux and EdgeX metering microservices are add on containers to each platform respectively. Both services query the messages sent from preconfigured devices to the platform to then output an assortment of total reading / device message counts to a designated port. In addition to displaying message counts, the microservices also output local CPU and memory usage metrics. These message counts can be collected and graphed with the setup of Telegraf, InfluxDB, and Chronograf. Through modifying the Telegraf configuration file, the message counts and system info can be consumed into InfluxDB via the HTTP input plugin. From there, the data can be visualized in Chronograf to track specific message counts as well as local system info. The TIC stack can be setup locally or in the cloud depending on individual implementations. The two services work independently of one another or can be run simultaneously if so desired. 

**Setup Note:** Each microservice is currently setup for one device with two readings (e.g. a sensor with temperature and humidity readings). However, the code is easily modifiable if the device has less than or more than two readings. Additionally, if more than one device needs to be tracked, multiple microservices can be run at once for each platform on different ports.

**Security Note:** As noted in the setup instructions, the protocol used throughout this tutorial is "http", making it insecure if implemented as is. All "http://" endpoints are configurable from the Docker configuration files, and the user can take the steps to make these "https://" if so desired. 

## Mainflux 

### Setup:
- Run Mainflux via Docker (https://github.com/mainflux/mainflux/)
  - Clone mainflux repo, cd into it, run ```make run```
- Use the CLI to setup devices and channels for actual or simulated devices (https://mainflux.readthedocs.io/en/latest/getting-started/)
  - Copy link to cli from releases page on Mainflux GitHub: ```wget LINK```
  - ```tar xvf TAR-FILE```
  - Create user: ```./mainflux-cli users create test@example.com test```
  - Get user token: ```./mainflux-cli users token test@example.com test```
  - Export user token: ```export USERTOKEN=user-token-from-above```
  - Create thing: ```./mainflux-cli things create '{"type": "device", "name" "testThing"}' $USERTOKEN```
  - Get info about thing: ```./mainflux-cli things get all $USERTOKEN```
  - Create channel: ```./mainflux-cli channels create '{"name": "testChannel"}' $USERTOKEN```
  - Get channel info: ```./mainflux-cli channels get all $USERTOKEN```
  - Connect channel to thing: ```./mainflux-cli things connect THINGID CHANNELID $USERTOKEN```
- Send messages from the device across Mainflux
  - Use the CLI to send messages if the device is simulated
    - ```./mainflux-cli messages send CHANNELID '[{"bn":"Dev1", "n":"temp", "v":20}, {"n":"hum","v":34}]' THINGKEY```
- Start the Mainflux -> InfluxDB writer service via Docker(https://github.com/mainflux/mainflux/tree/master/writers/influxdb)
- Start the Mainflux Metering Microservice via Docker
  - Add environmental variables into Docker file before running
    - Example Configuration: replace with correct IP Address
    ```yml
    PORT: 8915
    DEV1READ1: Dev1hum
    DEV1READ2: Dev1temp
    DEVICE1: dev1
    ADDR: http://{IP_Address}:8086
    ```
  - ```docker-compose up -d```
- Navigate to localhost:8915 to view message counts & system info
- Install Telegraf and Chronograf (https://portal.influxdata.com/downloads/)
  - Install InfluxDB if Mainflux / Mainflux Writer is not running on the same system 
  - InfluxDB is already running through the Mainflux writer service if it is running
- Configure the HTTP Telegraf input plugin to look at localhost:8915 with a timeout time of 10 seconds and an input format of 'json' (https://github.com/influxdata/telegraf/tree/master/plugins/inputs/http)
  - Telegraf configuration file is located in /etc/telegraf/telegraf.conf
- Restart Telegraf and attach Chronograf to the InfluxDB on port 8086
  - ```sudo systemctl restart telegraf```
  - Chronograf can be launched from localhost:8888
- Create Chronograf graphs to visualize Mainflux device message count data and system info by creating a new dashboard

### Data Flow: 
![Mainflux](images/mainflux.PNG)

## EdgeX

### Setup:
- Run EdgeX via Docker Compose (https://docs.edgexfoundry.org/Ch-QuickStart.html)
- Setup a device with the platform 
  - If creating a simulated device: (https://docs.edgexfoundry.org/Ch-Walkthrough.html)
- Start the EdgeX Metering Microservice via Docker 
  - Add environmental variables into Docker file before running
    - Example Configuration: replace with correct IP Address
    ```yml
    PORT: 8925
    READING1: caninecount
    READING2: humancount
    DEVICE: countcamera1
    READING_REQUEST1: http://{address}:48080/api/v1/reading/name/caninecount/10000
    READING_REQUEST2: http://{address}:48080/api/v1/reading/name/humancount/10000
    EVENT_REQUEST: http://{address}:48080/api/v1/event/count/countcamera1
    ```
  - ```docker-compose up -d```
- Send messages from the device
  - ```curl -X POST http://localhost:48080/api/v1/event -d '{"device":"countcamera1","readings":[{"name":"caninecount","value":"3"}, {"name":"humancount","value":"2"}]}' ```
- Navigate to localhost:8925 to view message counts & system info
- Install Telegraf, InfluxDB, & Chronograf (https://portal.influxdata.com/downloads/)
- Configure the HTTP Telegraf input plugin to look at localhost:8925 with a timeout of 10 seconds and an input format of 'json' (https://github.com/influxdata/telegraf/tree/master/plugins/inputs/http)
  - Telegraf configuration file is located in /etc/telegraf/telegraf.conf
- Restart Telegraf and attach Chronograf to InfluxDB on port 8086
  - ```sudo systemctl restart telegraf```
  - Chronograf can be launched from localhost:8888
- Create Chronograf graphs to visualize EdgeX device message count data and system info by creating a new dashboard

### Data Flow:
![EdgeX](images/edgex.PNG)

## AWS

### Setup:
- Start EC2 instance
- Adjust the instance's security group settings to allow an inbound rule on port 8888 (tcp protocol)
- SSH into instance and forward ports that are being used by microservice(s)
- Install Telegraf, InfluxDB, and Chronograf (https://portal.influxdata.com/downloads/)
- Configure the Telegraf HTTP input plugin to look at localhost:8915 and/or localhost:8925 with a timeout of 10 seconds and an input format of 'json' (https://github.com/influxdata/telegraf/tree/master/plugins/inputs/http)
  - The Telegraf configuration file can be found in /etc/telegraf/telegraf.conf
- Restart Telegraf
  - ```sudo systemctl restart telegraf``` 
- Navigate to the instance's IP:8888 and configure Chronograf in a browser to point to the InfluxDB instance at localhost:8086
- Create a new dashboard in Chronograf 
- Click on the "Add a Cell to Dashboard" button 
  - Select telegraf.autogen
  - Select http -> url -> pick the designated URL for the graph you wish to configure -> click on the appropriate field that you wish to monitor 
  - Alter the timeframe of the dashboard as needed (e.g. time > now() - 30m)
  - Title the graph accordingly and change the colors in the Visualization tab 
- Repeatedly add cells to the dashboard as needed to track various message counts and system info from Mainflux and/or EdgeX 


## Code Modification Tutorial 
_Use if need to configure a device with less than or more than two readings_

### Mainflux Metering Microservice
- Add additional variables in both the code and Docker configuration files
  - Ex: messages.go
   ```go
   // var declarations at top of code
   var dev1read1, dev1read2, dev1read3 json.Number
   var address, dev1read1Name, dev1read2Name, dev1read3Name string

   // top of main() in the rest of the enviro var declarations
   dev1read3Name = os.Getenv("DEV1READ3")

   // top of HandleFunc with other int conversions
   dev1read3Int, err3 := strconv.ParseUint(dev1read3.String(), 10, 64)
   if err1 == nil && err2 == nil && err3 == nil {
       dev1Msg = dev1read1Int + dev1read2Int + dev1read3Int
   }

   // bottom of the file with other output statements
   fmt.Fprintf(w, "\"%v Messages\": %v,\n", dev1read2Name, dev1read2)
   fmt.Fprintf(w, "\"%v Messages\": %v,\n", dev1read3Name, dev1read3)
   ```
  - Ex: docker-compose.yml
  ```yml
  environment:
    PORT: 8915
    DEV1READ1: Dev1hum
    DEV1READ2: Dev1temp
    DEV1READ3: Dev1time
    ...
  ```
- Add an additional query in the code 
  ```go
  param = map[string]interface{}{"name": dev1read3Name}
	q = client.NewQueryWithParameters("SELECT count(value) FROM messages WHERE \"name\"= $name", "mainflux", "", param)
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			dev1read3 = `0`
		} else {
			tempVal := response.Results[0].Series[0].Values[0]
		  dev1read3 = tempVal[countIndex].(json.Number)
		}
	}
  ```

### EdgeX Metering Microservice 
- Add additional variables in both the code and Docker configuration files
  - Ex: messages.go
  ```go
  // top of the file with the rest of the variable declarations
  var reading1Count, reading2Count, reading3Count, ... int
  var reading_request1, reading_request2, reading_request3, ... string

  // near the end of the query function
  totalDeviceCounts = reading1Count + reading2Count + reading3Count

  // top of main()
  var reading1, reading2, reading3, ... string
  reading3 = os.Getenv("READING3")
  reading_request3 = os.Getenv("READING_REQUEST3")

  //in HandleFunc after query() call
  fmt.Fprintf(w, "\"%v Messages\": %v,\n", reading3, reading3Count)
  ```
  - Ex: docker-compose.yml
  ```yml
  environment:
    PORT:8925
    READING1:caninecount
    READING2:humancount
    READING3:catcount
    ...
    READING_REQUEST3:http://{address}:48080/api/v1/reading/name/catcount/10000
    ...
  ```
- Add an additional query in the code
  ```go
  resp, err = http.Get(fmt.Sprintf("%s", reading_request3))
	if err != nil {
		reading3Count = 0
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		reading3Count = 0
	}
	reading3Count = (len(body) / msg_size) 
  ```

## Testing

Both microservices can be setup and tested through simulated devices. For the purposes of this tutorial, the following simulated devices can be created to mimic actual devices connected to each platform. 

**Mainflux**: a sensor with a temperature and humidity reading

**EdgeX**: a camera that detects and reports the number of humans and dogs in an image

### Mainflux 
- Configure the docker-compose file to the values in the tutorial above
- Move the mainflux_test file to the same directory as the Mainflux CLI
- Open the test file and replace CHANNEL_ID and THING_KEY with the actual values from setup 
- Start Mainflux and the Mainflux Metering Microservice 
- Run the test: ```./mainflux_test```

### EdgeX
- Configure the docker-compose file to the values in the tutorial above
- Move the edgex_test file to wherever is the most convenient
- Start EdgeX and the EdgeX Metering Microservice
- Run the test: ```./edgex_test```

Both tests simulate the devices sending a new message to the platform with updates to their respective fields every 2 minutes. The increase in message counts can be viewed in both the microservice endpoint and in a Chronograf dashboard that is tracking these particular message counts.










