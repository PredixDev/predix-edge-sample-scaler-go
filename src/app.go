package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type messageBody struct {
	Attributes map[string]string `json:"attributes"`
	Datapoints [][]interface{}   `json:"datapoints"`
	Name       string            `json:"name"`
}

type message struct {
	Body      []messageBody `json:"body"`
	MessageID string        `json:"messageId"`
}

//callback that indicates when the MQTT client is connected
var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Printf("Connection state changed")
}

//callback that indicates when the MQTT client connection is lost
var connectionLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, Error error) {
	fmt.Println("Connection unexpectedly lost:", Error)
}

func scaleData(msg message, tagToMatch string) message {

	bodyLength := len(msg.Body)
	for i := 0; i < bodyLength; i++ {
		tagName := msg.Body[i].Name
		if tagName == tagToMatch {
			data := msg.Body[i].Datapoints[0][1]
			doubleData, ok := data.(float64)
			if ok {
				doubleData = doubleData * 1000
				//put datapoint back into data
				msg.Body[i].Datapoints[0][1] = doubleData
				msg.Body[i].Name = tagName + ".scaled_by_1000" //Rename tag
			} else {
				fmt.Printf("Not a double")
			}
		}
	}
	return msg

}

//define a function for the default message handler
var defaultMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	var m message

	//Convert the message to JSON that fits our message structure
	err := json.Unmarshal(msg.Payload(), &m)
	if err != nil {
		fmt.Println(err)
	}
	//Take any tags with this name and multiply their data value by 1000
	for i := range m.Body  {
		data := m.Body[i].Datapoints[0][1]
		doubleData, ok := data.(float64)
		if ok {
			doubleData = doubleData * 1000
			//put datapoint back into data
			m.Body[i].Datapoints[0][1] = doubleData
			m.Body[i].Name = m.Body[i].Name + ".scaled_x_1000" //Rename tag
			fmt.Println("Changed data for "+m.Body[i].Name)
		}
	}

	//Turn data back into byte array
	var newMessage []byte
	newMessage, err = json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	//Republish to app_data topic
	pub, ok := os.LookupEnv("PUB_TOPIC")
	if !ok {
		fmt.Println("PUB_TOPIC environment variable not set")
		os.Exit(3)
	}
	if token := client.Publish(pub, 0, false, newMessage); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

}

//Used to set environment variables if running locally. Otherwise, environment variables are set in the docker-compose file
func setEnvironmentVariables() {

	os.Setenv("CLIENT_ID", "go-simple")
	os.Setenv("BROKER", "127.0.0.1:1883")
	os.Setenv("PUB_TOPIC", "timeseries_data")
	os.Setenv("SUB_TOPIC", "app_data")
	os.Setenv("TAG_NAME", "My.App.DOUBLE1")
}

func main() {

	//Get the command line arguments to see if we are running it locally
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "local" {
		setEnvironmentVariables()
	}

	//Read in broker, clientID and subscribe topic environment variables
	//If they aren't set, we can't proceed, so exit the program
	broker, ok := os.LookupEnv("BROKER")
	if !ok {
		fmt.Println("BROKER environment variable not set")
		os.Exit(3)
	}
	clientID, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		fmt.Println("CLIENT_ID environment variable not set")
		os.Exit(3)
	}
	sub, ok := os.LookupEnv("SUB_TOPIC")
	if !ok {
		fmt.Println("SUB_TOPIC environment variable not set")
		os.Exit(3)
	}

	//Open a channel to keep the connection open
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	//create a ClientOptions struct setting the broker address, clientid
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetAutoReconnect(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(500)

	//subscribe to the topic specificed in the "SUB_TOPIC" environment variable and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription

	if token := c.Subscribe(sub, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	<-channel
}
