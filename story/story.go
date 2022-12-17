package story

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl)) //template.Must ensures that we only return if html can be compiled else it will panic with the err
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
	<style>
		body {
			font-family: helvetica, arial;
		}
	
		h1 {
			text-align: center;
			position: relative;
		}
	
		.page {
			width: 80%;
			max-width: 500px;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
			padding: 80px;
			background: #FFFCF6;
			border: 1px solid #eee;
			box-shadow: 0 10px 6px -6px #777
		}

		ul {
			border-top: 1px dotted #ccc;
			padding-top: 10px;
			-webkit-padding-start: 0;
		}

		li {
			padding-top: 10px;
		}

		a,
		a:visited {
			text-decoration: none;
			color: #6295b5;
		}
		a:active,
		a:hover {
			color: #7792a2;
		}

		p {
			text-indent: 1em;
		}
   </style>
</head>

<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
		<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
			<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</section>
</body>

</html>
`

type HandlerOpt func(h *handler)

func WithTemplate(t *template.Template) HandlerOpt {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFn(pathFn func(r *http.Request) string) HandlerOpt {
	return func(h *handler) {
		h.pathFn = pathFn
	}
}

func HttpHandler(s Story, opts ...HandlerOpt) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path) //get the URL Path of the http requests with any whitespaces removed
	if path == "" || path == "/" {
		path = "/intro" //if response has no URL path, assume user is starting story from scratch (restart from intro)
	}
	return path[1:] //remove the '/' in path
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	//check if path name exists in the JSON file (e.g. '/new-york' would be a valid path URL as 'new-york' is a key to a Chapter in Story map)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter) //http.ResponseWriter is an io.Writer => can use it in Template.Execute to write story chapter to response body
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r) //return a decoder that reads from the os.File
	var story Story
	if err := d.Decode(&story); err != nil { //Decode reads the JSON-encoded value from f and stores it in story map.
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
