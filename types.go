package main

import (
	"time"
)

//The Message type will encapsulate the Message string itself, but we have also added
//the Name and When fields that respectively hold the user's name and a timestamp of
//when the Message was sent.
type Message struct {
	Name    string
	Message string
	When    time.Time
	Hub     string
}
