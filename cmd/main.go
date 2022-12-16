package main

import (
	"cyoa/story"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	fileName := flag.String("file", "gopher.json", "the JSON file to generate CYOA story") //let user choose which json file to generate storybook from
	flag.Parse()                                                                           //allow program to access flag val
	fmt.Printf("Using the story in %s.\n", *fileName)

	f, err := os.Open(*fileName) //open the json file and return a pointer to an os.File which implements io.Reader
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f) //return a decoder that reads from the os.File
	var story story.Story
	if err := d.Decode(&story); err != nil { //Decode reads the JSON-encoded value from f and stores it in story map.
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
