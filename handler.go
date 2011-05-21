package starrpg

import (
	"fmt"
	"http"
	"io"
	"json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Storage interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
	Delete(key string) bool
	Inc(key string) (uint64, bool)
}

type DummyStorage map[string][]byte

func (s *DummyStorage) Get(key string) ([]byte, bool) {
	item, ok := (*s)[key]
	return item, ok
}

func (s *DummyStorage) Set(key string, value []byte) {
	(*s)[key] = value
}

func (s *DummyStorage) Delete(key string) bool {
	if _, ok := (*s)[key]; !ok {
		return false
	}
	(*s)[key] = []byte{}, false
	return true
}

func (s *DummyStorage) Inc(key string) (uint64, bool) {
	value, ok := (*s)[key]
	if !ok {
		(*s)[key] = []byte("1")
		return 1, true
	}
	numValue, err := strconv.Atoui64(string(value))
	if err != nil {
		return 0, false
	}
	(*s)[key] = []byte(strconv.Uitoa64(numValue + 1))
	return numValue + 1, true
}

func checkAcceptHeader(mediaType, accept string) float64 {
	splitedMediaType := strings.Split(mediaType, "/", -1)
	if len(splitedMediaType) != 2 {
		return 0
	}
	if splitedMediaType[0] == "*" || splitedMediaType[1] == "*" {
		return 0
	}
	result := 0.0
	for _, mediaRange := range strings.Split(accept, ",", -1) {
		qValue := 1.0
		if i := strings.LastIndex(mediaRange, ";q="); i != -1 {
			newQValue, err := strconv.Atof64(mediaRange[i+3:])
			if err != nil {
				return 0
			}
			mediaRange = mediaRange[0:i]
			qValue = newQValue
		}
		splitedMediaRange := strings.Split(mediaRange, "/", -1)
		if len(splitedMediaRange) != 2 {
			return 0
		}
		if mediaType == mediaRange {
			if result < qValue {
				result = qValue
			}
		} else if splitedMediaRange[1] == "*" {
			if splitedMediaRange[0] == "*" ||
				splitedMediaRange[0] == splitedMediaType[0] {
				if result < qValue {
					result = qValue
 				}
			}
		} else if splitedMediaRange[0] == "*" {
			return 0
		}
	}
	return result
}

func handleHome(conn http.ResponseWriter, req *http.Request) {
	if checkAcceptHeader("application/xhtml+xml", req.Header.Get("Accept")) <= 0 {
		conn.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	data := `<?xml version="1.0" encoding="utf-8" ?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="ja">
  <body>
    <p>Hello, World!</p>
  </body>
</html>`
	// TODO: gzip
	conn.Header().Set("Content-Type", "application/xhtml+xml; charset=utf-8")
	if _, err := io.WriteString(conn, data); err != nil {
		log.Print("io.WriteString: ", err)
	}
}

func isGettablePath(path string) bool {
	var pathRegExp = regexp.MustCompile("^/games(/[a-zA-Z0-9_\\-]+(/(maps|planes|items)(/[a-zA-Z0-9_\\-]+)?)?)?$")
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func isPostablePath(path string) bool {
	var pathRegExp = regexp.MustCompile("^/games$")
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func isPuttablePath(path string) bool {
	var pathRegExp = regexp.MustCompile("^/games/[a-zA-Z0-9_\\-]+(/(maps|planes|items)/[a-zA-Z0-9_\\-]+)?$")
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func isDeletablePath(path string) bool {
	var pathRegExp = regexp.MustCompile("^/games/[a-zA-Z0-9_\\-]+(/(maps|planes|items)/[a-zA-Z0-9_\\-]+)?$")
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func getAllowHeader(path string) string {
	allow := "OPTIONS"
	if (isGettablePath(path)) {
		allow += ", GET, HEAD"
	}
	if (isPostablePath(path)) {
		allow += ", POST"
	}
	if (isPuttablePath(path)) {
		allow += ", PUT"
	}
	if (isDeletablePath(path)) {
		allow += ", DELETE"
	}
	return allow
}

func sendResponseMethodNotAllowed(conn http.ResponseWriter, req *http.Request) {
	conn.Header().Set("Allow", getAllowHeader(req.URL.Path))
	conn.WriteHeader(http.StatusMethodNotAllowed)
}

func doPost(storage Storage, path string) (string, os.Error) {
	newID, ok := storage.Inc(path + "/*count")
	if !ok {
		return "", os.NewError(fmt.Sprintf(`storage.Inc(%#v + "/inner-count") failed!`, path))
	}
	values, ok := storage.Get(path)
	if !ok {
		values = []byte("{}")
	}
	var items map[string]map[string]string
	if err := json.Unmarshal(values, &items); err != nil {
		return "", err
	}
	items[strconv.Uitoa64(newID)] = map[string]string{"name": ""}
	newBytes, err := json.Marshal(items)
	if err != nil {
		return "", os.NewError(fmt.Sprintf(`json.Marshal(%#v) failed!`, items))
	}
	storage.Set(path, newBytes)
	return path + "/" + strconv.Uitoa64(newID), nil
}

type ResourceHandler struct {
	Storage
}

const (
	httpMethodNone = iota
	httpMethodOptions
	httpMethodGet
	httpMethodHead
	httpMethodPost
	httpMethodPut
	httpMethodDelete
)

func (r *ResourceHandler) Handle(conn http.ResponseWriter, req *http.Request) {
	if checkAcceptHeader("application/json", req.Header.Get("Accept")) <= 0 {
		conn.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	path := req.URL.Path
	httpMethod := httpMethodNone
	switch req.Method {
	case "OPTIONS":
		httpMethod = httpMethodOptions
	case "GET":
		httpMethod = httpMethodGet
	case "HEAD":
		httpMethod = httpMethodHead
	case "POST":
		httpMethod = httpMethodPost
	case "PUT":
		httpMethod = httpMethodPut
	case "DELETE":
		httpMethod = httpMethodDelete
	}
	// TODO: 権限チェック
	switch httpMethod {
	case httpMethodOptions:
		conn.Header().Set("Content-Length", "0")
		if path == "*" {
			conn.Header().Set("Allow", "OPTIONS, GET, HEAD, POST, PUT, DELETE")
		} else {
			conn.Header().Set("Allow", getAllowHeader(path))
		}
		conn.WriteHeader(http.StatusOK)
	case httpMethodGet, httpMethodHead:
		if (!isGettablePath(path)) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		data, ok := r.Storage.Get(path)
		if !ok {
			http.NotFound(conn, req)
			return
		}
		conn.Header().Set("Content-Type", "application/json; charset=utf-8")
		conn.WriteHeader(http.StatusOK)
		if httpMethod == httpMethodHead {
			return
		}
		conn.Write(data)
	case httpMethodPost:
		if (!isPostablePath(path)) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		newPath, err := doPost(r.Storage, path)
		if err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		newURL := req.URL.Scheme + "://" + req.URL.RawAuthority + newPath
		conn.Header().Set("Location", newURL)
		conn.WriteHeader(http.StatusCreated)
	case httpMethodPut:
		if (!isPuttablePath(path)) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		buf := make([]byte, 4096)
		size, err := io.ReadFull(req.Body, buf)
		if err == nil {
			conn.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		if err != os.EOF {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := buf[:size]
		// TODO: JSON 形式チェック
		r.Storage.Set(path, body)
		conn.WriteHeader(http.StatusOK)
	case httpMethodDelete:
		if (!isDeletablePath(path)) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		// TODO: 子リソースの再帰的削除
		if !r.Storage.Delete(path) {
			conn.WriteHeader(http.StatusNotFound)
			return
		}
		conn.WriteHeader(http.StatusOK)
	default:
		sendResponseMethodNotAllowed(conn, req)
	}
}

func Handler(conn http.ResponseWriter, req *http.Request) {
	switch path := req.URL.Path; {
	case path == "/":
		handleHome(conn, req)
	default:
		resourceHandler := &ResourceHandler{&DummyStorage{}}
		resourceHandler.Handle(conn, req)
	}
}
