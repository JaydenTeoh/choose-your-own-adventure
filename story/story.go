package story

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
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
</head>

<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>

</html>
`

func HttpHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"]) //http.ResponseWriter is an io.Writer => can use it in Template.Execute to write to response body
	if err != nil {
		panic(err)
	}
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
