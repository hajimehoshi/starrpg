package starrpg

import (
	"http"
	"json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"path/filepath"
)

type resourceRequestProcessor struct {
	ResourceStorage
}

func NewResourceRequestProcessor(resourceStorage ResourceStorage) RequestProcessor {
	return &resourceRequestProcessor{resourceStorage}
}

func (r *resourceRequestProcessor) DoOptions(req *http.Request) (int, map[string]string, os.Error) {
	responseHeader := map[string]string{"Content-Length": "0"}
	path := req.URL.Path
	if path == "*" {
		responseHeader["Allow"] = "OPTIONS, GET, HEAD, POST, PUT, DELETE"
	} else {
		responseHeader["Allow"] = getAllowHeader(path)
	}
	return http.StatusOK, responseHeader, nil
}

func (r *resourceRequestProcessor) DoHead(req *http.Request) (int, map[string]string, os.Error) {
	status, responseHeader, _, err := r.DoGet(req)
	return status, responseHeader, err
}

func (r *resourceRequestProcessor) DoGet(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	fileContent, err := GetFileFromCache(filepath.Join("public", path))
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if fileContent != nil {
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
		return http.StatusOK, map[string]string{"Content-Type":contentType}, fileContent, nil
	}
	acceptHeader := req.Header.Get("Accept")
	jsonQVal := checkAcceptHeader("application/json", acceptHeader)
	xhtmlQVal := checkAcceptHeader("application/xhtml+xml", acceptHeader)
	htmlQVal := checkAcceptHeader("text/html", acceptHeader)
	if jsonQVal == 0 && xhtmlQVal == 0 && htmlQVal == 0 {
		// NOT ACCEPTABLE?
		return http.StatusNotFound, nil, nil, nil
	}
	if xhtmlQVal <= jsonQVal && htmlQVal <= jsonQVal {
		if len(path) == 0 {
			return http.StatusNotFound, nil, nil, nil
		}
		if 1 < len(path) && path[len(path) - 1] == '/' {
			path = path[:len(path) - 1]
		}
		slashCount := strings.Count(path, "/")
		if slashCount == 0 || slashCount == 1 {
			return http.StatusNotFound, nil, nil, nil
		}
		var obj interface{}
		switch slashCount % 2 {
		case 0:
			obj2, err := r.ResourceStorage.Get(path)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, err
			}
			obj = obj2
		case 1:
			obj2, err := r.ResourceStorage.GetChildren(path)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, err
			}
			obj = obj2
		}
		if obj == nil {
			return http.StatusNotFound, nil, nil, nil
		}
		content, err := json.Marshal(obj)
		if err != nil {
			return http.StatusInternalServerError, nil, nil, err
		}
		contentType := "application/json; charset=utf-8"
		return http.StatusOK, map[string]string{"Content-Type": contentType}, content, nil
	}
	htmlPath := getHTMLViewPath(path)
	if htmlPath == "" {
		return http.StatusNotFound, nil, nil, nil
	}
	htmlContent, err := GetFileFromCache(htmlPath)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if htmlContent == nil {
		return http.StatusNotFound, nil, nil, nil
	}
	contentType := "application/xhtml+xml; charset=utf-8";
	return http.StatusOK, map[string]string{"Content-Type": contentType}, htmlContent, nil
}

func (r *resourceRequestProcessor) DoPost(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isPostablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	newID, err := r.ResourceStorage.Create(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if newID == 0 {
		return http.StatusNotFound, nil, nil, nil
	}
	newPath := path + "/" + strconv.Uitoa64(newID)
	// TODO: これはやるべきか?
	/*if err := rs.Set(newItemPath, map[string]string{}); err != nil {
		return "", err
	}*/
	// TODO: https?
	location := "http://" + req.Host + newPath
	return http.StatusCreated, map[string]string{"Location": location}, nil, nil
}

func (r *resourceRequestProcessor) DoPut(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isPuttablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
		return http.StatusUnsupportedMediaType, nil, nil, nil
	}
	buf := make([]byte, 4096)
	size, err := req.Body.Read(buf)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if _, err := req.Body.Read(make([]byte, 1)); err != os.EOF {
		return http.StatusRequestEntityTooLarge, nil, nil, err
	}
	requestBody := buf[:size]
	obj := map[string]string{}
	if err := json.Unmarshal(requestBody, &obj); err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	if err := r.ResourceStorage.Set(path, obj); err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	return http.StatusOK, nil, nil, nil
}

func (r *resourceRequestProcessor) DoDelete(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isDeletablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	// TODO: 子リソースの再帰的削除
	return http.StatusOK, nil, nil, nil
}
