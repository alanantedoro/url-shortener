package urlshortener

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(ymlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(ymlBytes, &pathUrls)
	if err != nil{
		fmt.Println(err)
	}

	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls{
		pathsToUrls[pu.Path] = pu.Url
	}
	return MapHandler(pathsToUrls, fallback), nil
}

type pathURL struct{
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}