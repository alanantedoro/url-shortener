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
	pathUrls, err := parseYaml(ymlBytes)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(data []byte)([]pathURL, error){
	var pathUrls []pathURL
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil{
		fmt.Println(err)
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathURL) map[string]string{
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls{
		pathsToUrls[pu.Path] = pu.Url
	}
	return pathsToUrls
}

type pathURL struct{
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}


//TODO: extras and tests 