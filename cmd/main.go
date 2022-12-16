package main

import (
	cyoa "cyoa/story"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3000, "port to start the CYOA application")                   //generate a port number to listen for later
	fileName := flag.String("file", "gopher.json", "the JSON file to generate CYOA story") //let user choose which json file to generate storybook from
	flag.Parse()                                                                           //allow program to access flag val
	fmt.Printf("Using the story in %s.\n", *fileName)

	f, err := os.Open(*fileName) //open the json file and return a pointer to an os.File which implements io.Reader
	checkErr(err)

	story, err := cyoa.JsonStory(f) //decode the json file using io.Reader and return a Story map with accessible struct values
	checkErr(err)

	h := cyoa.HttpHandler(story) //return makeshift http.Handler interface that implements ServeHttp function that allow us to write story intro to initial response body and also handle future http requests
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h)) //listen for port 3000 and use h to handler http requests
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
