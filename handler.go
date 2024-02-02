package urlshortener

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
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

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error){
	pathUrls, err := parseJson(jsonBytes)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func BOLThanlder(db *bolt.DB, fallback http.Handler)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("pathsToUrls"))
			if bucket != nil {
				cursor := bucket.Cursor()
				for path, url := cursor.First(); path != nil; path, url = cursor.Next() {
					if string(path) == r.URL.Path {
						http.Redirect(w, r, string(url), http.StatusFound)
						return nil
					}
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
		fallback.ServeHTTP(w, r)
	}
}

func parseYaml(data []byte)([]pathURL, error){
	var pathUrls []pathURL
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil{
		fmt.Println(err)
	}
	return pathUrls, nil
}

func parseJson(data []byte)([]pathURL, error){
	var pathUrls []pathURL
	err := json.Unmarshal(data, &pathUrls)
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


//TODO:
// Update the main/main.go source file to accept a YAML file as a flag and then load the YAML from a file rather than from a string. DONE
// Build a JSONHandler that serves the same purpose, but reads from JSON data. DONE
// Build a Handler that doesn’t read from a map but instead reads from a database. Whether you use BoltDB, SQL, or something else is entirely up to you.
// VER COMO AGREGAR A LA DB POR COMANDO O ALGO ASI PODEMOS TESTEAR
// TESTS