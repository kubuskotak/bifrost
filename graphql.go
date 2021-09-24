package bifrost

import (
	"embed"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	graph "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/kubuskotak/bifrost/assets"
	"github.com/rs/zerolog/log"
)

func GetRootSchema(graphql embed.FS, dir string) (string, error) {
	prefix := dir
	var buff strings.Builder
	buff.Reset()

	fs, err := graphql.ReadDir(prefix)
	for _, f := range fs {
		if filepath.Ext(f.Name()) != ".graphql" {
			continue
		}

		pathPattern := path.Join(prefix, f.Name())
		data, err := graphql.ReadFile(pathPattern)
		if err != nil {
			return "", err
		}

		_, err = buff.Write(data)
		if err != nil {
			return "", err
		}
	}
	return buff.String(), err
}

// Graphql handler func
func Graphql(graphql embed.FS, dir string, resolver interface{}, opts ...graph.SchemaOpt) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, erByte := GetRootSchema(graphql, dir)
		if erByte != nil {
			log.Error().Err(erByte)
			_ = ResponseJSONPayload(w, r, http.StatusNoContent, nil)
			return
		}
		sch := graph.MustParseSchema(bytes, resolver, opts...)
		handler := &relay.Handler{Schema: sch}
		handler.ServeHTTP(w, r)
	}
}

// Graph handler func
func Graph(endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index, erIndex := assets.Assets.ReadFile(`index.html`)
		if erIndex != nil {
			log.Error().Err(erIndex)
			_ = ResponseJSONPayload(w, r, http.StatusNoContent, nil)
			return
		}
		tmpl := template.Must(template.New("svelte").Parse(string(index)))
		if err := tmpl.Execute(w, map[string]string{
			"endpoint": endpoint,
		}); err != nil { // Execute template with data
			log.Error().Err(err)
			_ = ResponseJSONPayload(w, r, http.StatusNoContent, nil)
		}
	}
}
