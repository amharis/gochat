// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
)

var chatroomListener = flag.String("chat-backend", "0.0.0.0:8080", "chatroom backend's address")
var redisConnectionString = flag.String("redis-connection-string", "localhost:6379", "redis server's address")
var CHATROOM = "chatroom"
var ctx = context.Background()

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	fmt.Println("Starting Application")
	flag.Parse()
	fmt.Println("input args: ", *chatroomListener, *redisConnectionString)
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*chatroomListener, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func setupLogging() {
	// setup logging, useful in local non-docker setups when stdout is difficult to manage
	/*
		lfname := fmt.Sprintf("chat-%s.log", hub.ID)
		f, err := os.OpenFile(lfname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		//defer to close when you're done with it, not because you think it's idiomatic!
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("Error in closing log file: ", err)
			}
		}(f)
		//set output of logs to f
		log.SetOutput(f)
		fmt.Println("Writing to log file: ", lfname)
	*/
}
