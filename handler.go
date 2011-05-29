package starrpg

import (
	"http"
	"json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Storage interface {
	Get(key string) []byte
	GetWithPrefix(key string) (map[string][]byte)
	Set(key string, value []byte)
	Delete(key string) bool
	Update(key string, f func([]byte) ([]byte, os.Error)) os.Error
}

type MapStorage interface {
	Get(key string) (map[string]string, os.Error)
	GetWithPrefix(prefix string) (map[string]map[string]string, os.Error)
	Set(key string, obj map[string]string) os.Error
	Delete(key string) bool
	Inc(key, subKey string) (uint64, os.Error)
}

type ResourceStorage interface {
	Get(urlPath string) (map[string]string, os.Error)
	GetChildren(urlPath string) (map[string]map[string]string, os.Error)
	Set(urlPath string, obj map[string]string) os.Error
	Create(urlPath string) (uint64, os.Error)
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
	htmlPath := getHTMLViewPath(req.URL.Path)
	data, err := GetFileFromCache(htmlPath)
	if err != nil {
		log.Print(err)
		conn.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		log.Print(htmlPath + " was not found")
		conn.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: gzip
	conn.Header().Set("Content-Type", "application/xhtml+xml; charset=utf-8")
	if _, err := conn.Write(data); err != nil {
		log.Print("io.WriteString: ", err)
	}
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

func getHTTPMethod(method string) int {
	switch method {
	case "OPTIONS":
		return httpMethodOptions
	case "GET":
		return httpMethodGet
	case "HEAD":
		return httpMethodHead
	case "POST":
		return httpMethodPost
	case "PUT":
		return httpMethodPut
	case "DELETE":
		return httpMethodDelete
	}
	return httpMethodNone
}

func isPostablePath(path string) bool {
	var pathRegExp = regexp.MustCompile(`^/games$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func isPuttablePath(path string) bool {
	pathRegExp := regexp.MustCompile(`^/games/[a-zA-Z0-9_\-]+(/(maps|items)/[a-zA-Z0-9_\-]+)?$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func isDeletablePath(path string) bool {
	pathRegExp := regexp.MustCompile(`^/games/[a-zA-Z0-9_\-]+(/(maps|items)/[a-zA-Z0-9_\-]+)?$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func getAllowHeader(path string) string {
	allow := "OPTIONS, GET, HEAD"
	if isPostablePath(path) {
		allow += ", POST"
	}
	if isPuttablePath(path) {
		allow += ", PUT"
	}
	if isDeletablePath(path) {
		allow += ", DELETE"
	}
	return allow
}

func sendResponseNotFound(conn http.ResponseWriter, req *http.Request) {
	acceptHeader := req.Header.Get("Accept")
	xhtmlQVal := checkAcceptHeader("application/xhtml+xml", acceptHeader)
	htmlQVal := checkAcceptHeader("text/html", acceptHeader)
	if xhtmlQVal == 0 && htmlQVal == 0 {
		conn.WriteHeader(http.StatusNotFound)
		return
	}
	notFoundFilePath := filepath.Join("views", "not_found.html")
	content, err := GetFileFromCache(notFoundFilePath)
	if err != nil {
		log.Print(err)
		conn.WriteHeader(http.StatusInternalServerError)
		return
	}
	if content == nil {
		log.Print(notFoundFilePath + " was not found")
		conn.WriteHeader(http.StatusInternalServerError)
		return
	}
	conn.Header().Set("Content-Type", "application/xhtml+xml; charset=utf-8")
	conn.WriteHeader(http.StatusNotFound)
	conn.Write(content)
}

func sendResponseMethodNotAllowed(conn http.ResponseWriter, req *http.Request) {
	conn.Header().Set("Allow", getAllowHeader(req.URL.Path))
	conn.WriteHeader(http.StatusMethodNotAllowed)
}

func getHTMLViewPath(path string) string {
	if path == "/" {
		return filepath.Join("views", "home.html")
	}
	regexpGame := regexp.MustCompile(`^/games/[a-zA-Z0-9_\-]+$`)
	if regexpGame.MatchString(path) {
		return filepath.Join("views", "game.html")
	}
	return ""
}

func doGet(rs ResourceStorage, path string, acceptHeader string) (string, []byte, os.Error) {
	{
		content, err := GetFileFromCache(filepath.Join("public", path))
		if err != nil {
			return "", nil, err
		}
		if content != nil {
			// TODO: check acceptHeader?
			var contentType string
			switch {
			case strings.HasSuffix(path, ".js"):
				contentType = "text/javascript; charset=utf-8"
			case strings.HasSuffix(path, ".css"):
				contentType = "text/css; charset=utf-8"
			case path == "/favicon.ico":
				contentType = "image/png" // TODO: modify?
			default:
				contentType = "application/octet-stream"
			}
			return contentType, content, nil
		}
	}
	jsonQVal := checkAcceptHeader("application/json", acceptHeader)
	xhtmlQVal := checkAcceptHeader("application/xhtml+xml", acceptHeader)
	htmlQVal := checkAcceptHeader("text/html", acceptHeader)
	if jsonQVal == 0 && xhtmlQVal == 0 && htmlQVal == 0 {
		return "", nil, nil
	}
	if xhtmlQVal <= jsonQVal && htmlQVal <= jsonQVal {
		if len(path) == 0 {
			return "", nil, nil
		}
		if 1 < len(path) && path[len(path) - 1] == '/' {
			path = path[:len(path) - 1]
		}
		slashCount := strings.Count(path, "/")
		if slashCount == 0 || slashCount == 1 {
			return "", nil, nil
		}
		var obj interface{}
		switch slashCount % 2 {
		case 0:
			obj2, err := rs.Get(path)
			if err != nil {
				return "", nil, err
			}
			obj = obj2
		case 1:
			obj2, err := rs.GetChildren(path)
			if err != nil {
				return "", nil, err
			}
			obj = obj2
		}
		if obj == nil {
			return "", nil, nil
		}
		content, err := json.Marshal(obj)
		if err != nil {
			return "", nil, err
		}
		return "application/json; charset=utf-8", content, nil
	}
	htmlPath := getHTMLViewPath(path)
	if htmlPath == "" {
		return "", nil, nil
	}
	content, err := GetFileFromCache(htmlPath)
	if err != nil {
		return "", nil, err
	}
	if content == nil {
		return "", nil, nil
	}
	return "application/xhtml+xml; charset=utf-8", content, nil
}

func doPost(rs ResourceStorage, path string) (string, os.Error) {
	newID, err := rs.Create(path)
	if err != nil {
		return "", err
	}
	if newID == 0 {
		return "", nil
	}
	newItemPath := path + "/" + strconv.Uitoa64(newID)
	// TODO: これはやるべきか?
	if err := rs.Set(newItemPath, map[string]string{}); err != nil {
		return "", err
	}
	return newItemPath, nil
}

type ResourceHandler struct {
	ResourceStorage
}

func (r *ResourceHandler) Handle(conn http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	httpMethod := getHTTPMethod(req.Method)
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
		contentType, content, err := doGet(r.ResourceStorage, path, req.Header.Get("Accept"))
		if err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: returns 406?
		if len(content) == 0 {
			log.Print(path, " not found")
			sendResponseNotFound(conn, req)
			return
		}
		conn.Header().Set("Content-Type", contentType)
		conn.WriteHeader(http.StatusOK)
		if httpMethod == httpMethodHead {
			return
		}
		conn.Write(content)
	case httpMethodPost:
		if !isPostablePath(path) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		newPath, err := doPost(r.ResourceStorage, path)
		if err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		if newPath == "" {
			sendResponseNotFound(conn, req)
			return
		}
		// TODO: fix schema
		newURL := "http://" + req.Host + newPath
		conn.Header().Set("Location", newURL)
		conn.WriteHeader(http.StatusCreated)
	case httpMethodPut:
		if !isPuttablePath(path) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
			conn.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		buf := make([]byte, 4096)
		size, err := req.Body.Read(buf)
		if err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := req.Body.Read(make([]byte, 1)); err != os.EOF {
			conn.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		body := buf[:size]
		obj := map[string]string{}
		if err := json.Unmarshal(body, &obj); err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.ResourceStorage.Set(path, obj); err != nil {
			log.Print(err)
			conn.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn.WriteHeader(http.StatusOK)
	case httpMethodDelete:
		if !isDeletablePath(path) {
			sendResponseMethodNotAllowed(conn, req)
			return
		}
		// TODO: 子リソースの再帰的削除
		/*if !r.MapStorage.Delete(urlPathToStoragePath(path)) {
			sendResponseNotFound(conn, req)
			return
		}*/
		conn.WriteHeader(http.StatusOK)
	default:
		sendResponseMethodNotAllowed(conn, req)
	}
}

var (
	storage_ = &DummyStorage{}
	mapStorage_ = NewMapStorage(storage_)
	resourceStorage_ = NewResourceStorage(mapStorage_)
)

func Handler(conn http.ResponseWriter, req *http.Request) {
	switch path := req.URL.Path; {
	case path == "/":
		handleHome(conn, req)
	default:
		resourceHandler := &ResourceHandler{resourceStorage_}
		resourceHandler.Handle(conn, req)
	}
}
