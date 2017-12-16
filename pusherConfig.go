package main

import "fmt"

type PusherConfig struct {
	AppId   string
	Key     string
	Secret  string
	Cluster string
	Secure  bool
	Channel string
	Event   string
}

// Update with configurations from Pusher.com
func GetPusherConfig() PusherConfig {
	return PusherConfig{
		AppId:   "",
		Key:     "",
		Secret:  "",
		Cluster: "us2",
		Secure:  true,
		Channel: "my-channel",
		Event:   "my-event"}
}

func TestFmt() {
	fmt.Print("Testing")
}