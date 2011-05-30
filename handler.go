package starrpg

import (
	"http"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type RequestProcessor interface {
	DoOptions(req *http.Request) (int, map[string]string, os.Error)
	DoGet(req *http.Request) (int, map[string]string, []byte, os.Error)
	DoHead(req *http.Request) (int, map[string]string, os.Error)
	DoPost(req *http.Request) (int, map[string]string, []byte, os.Error)
	DoPut(req *http.Request) (int, map[string]string, []byte, os.Error)
	DoDelete(req *http.Request) (int, map[string]string, []byte, os.Error)
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
	// TODO: check method
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

type ResourceHandler struct {
	ResourceRequestProcessor RequestProcessor
}

func (r *ResourceHandler) sendResponseNotFound(conn http.ResponseWriter, req *http.Request) {
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

func (r *ResourceHandler) Handle(conn http.ResponseWriter, req *http.Request) {
	requestProcessor := r.ResourceRequestProcessor
	// TODO: check authority
	var status int = http.StatusInternalServerError
	var header map[string]string
	var content []byte
	var err os.Error
	// TODO: regarding for _method parameter
	switch req.Method {
	case "OPTIONS":
		status, header, err = requestProcessor.DoOptions(req)
	case "GET":
		status, header, content, err = requestProcessor.DoGet(req)
	case "HEAD":
		status, header, err = requestProcessor.DoHead(req)
	case "POST":
		status, header, content, err = requestProcessor.DoPost(req)
	case "PUT":
		status, header, content, err = requestProcessor.DoPut(req)
	case "DELETE":
		status, header, content, err = requestProcessor.DoDelete(req)
	default:
		status = http.StatusMethodNotAllowed
	}
	if err != nil {
		log.Print(err)
	}
	for key, value := range header {
		conn.Header().Set(key, value)
	}
	if status == http.StatusNotFound {
		r.sendResponseNotFound(conn, req)
		return
	}
	conn.WriteHeader(status)
	if content == nil {
		return
	}
	conn.Write(content)
}

var (
	storage_ = &DummyStorage{}
	mapStorage_ = NewMapStorage(storage_)
	resourceStorage_ = NewResourceStorage(mapStorage_)
	resourceRequestProcessor_ = NewResourceRequestProcessor(resourceStorage_)
)

func Handler(conn http.ResponseWriter, req *http.Request) {
	switch path := req.URL.Path; {
	case path == "/":
		handleHome(conn, req)
	default:
		resourceHandler := &ResourceHandler{resourceRequestProcessor_}
		resourceHandler.Handle(conn, req)
	}
}
