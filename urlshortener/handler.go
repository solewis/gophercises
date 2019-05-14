package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"encoding/json"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	t := redirect{}

	err := yaml.Unmarshal(yml, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	return MapHandler(redirectToMap(t), fallback), nil
}

func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	t := redirect{}

	if err := json.Unmarshal(j, &t); err != nil {
        log.Fatalf("error: %v", err)
		return nil, err
    }

    return MapHandler(redirectToMap(t), fallback), nil
}

type redirect []struct {
	Path string
	Url string
}

func redirectToMap(redirects redirect) map[string]string {
	m := make(map[string]string)
	for _, redirect := range redirects {
		m[redirect.Path] = redirect.Url
	}
	return m
}
