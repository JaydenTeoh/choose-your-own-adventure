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

	// tpl := template.Must(template.New("").Parse(customTmpl))

	h := cyoa.HttpHandler(story) //return makeshift http.Handler interface that implements ServeHttp function that allow us to write story intro to initial response body and also handle future http requests
	mux := http.NewServeMux()
	mux.Handle("/", h)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)) //listen for port 3000 and use h to handler http requests
}

// func pathFn(r *http.Request) string {
// 	path := strings.TrimSpace(r.URL.Path)
// 	if path == "/story" || path == "/story/" || path == "/" {
// 		path = "/story/intro"
// 	}
// 	return path[len("/story/"):]
// }

// var customTmpl = `
// <!DOCTYPE html>
// <html lang="en">

// <head>
//     <meta charset="UTF-8">
//     <title>Choose Your Own Adventure</title>
// 	<style>
// 		body {
// 			font-family: helvetica, arial;
// 		}

// 		h1 {
// 			text-align: center;
// 			position: relative;
// 		}

// 		.page {
// 			width: 80%;
// 			max-width: 500px;
// 			margin: auto;
// 			margin-top: 40px;
// 			margin-bottom: 40px;
// 			padding: 80px;
// 			background: #FFFCF6;
// 			border: 1px solid #eee;
// 			box-shadow: 0 10px 6px -6px #777
// 		}

// 		ul {
// 			border-top: 1px dotted #ccc;
// 			padding-top: 10px;
// 			-webkit-padding-start: 0;
// 		}

// 		li {
// 			padding-top: 10px;
// 		}

// 		a,
// 		a:visited {
// 			text-decoration: none;
// 			color: #6295b5;
// 		}
// 		a:active,
// 		a:hover {
// 			color: #7792a2;
// 		}

// 		p {
// 			text-indent: 1em;
// 		}
//    </style>
// </head>

// <body>
// 	<section class="page">
// 		<h1>{{.Title}}</h1>
// 		{{range .Paragraphs}}
// 		<p>{{.}}</p>
// 		{{end}}
// 		<ul>
// 			{{range .Options}}
// 			<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
// 			{{end}}
// 		</ul>
// 	</section>
// </body>

// </html>
// `

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
