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
	resourceObj, err := r.ResourceStorage.Get(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if resourceObj == nil {
		return http.StatusNotFound, nil, nil, err
	}
	acceptHeader := req.Header.Get("Accept")
	jsonQVal := checkAcceptHeader("application/json", acceptHeader)
	xhtmlQVal := checkAcceptHeader("application/xhtml+xml", acceptHeader)
	htmlQVal := checkAcceptHeader("text/html", acceptHeader)
	if jsonQVal == 0 && xhtmlQVal == 0 && htmlQVal == 0 {
		return http.StatusNotAcceptable, nil, nil, nil
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
			obj = resourceObj
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
		responseHeader := map[string]string{"Content-Type": contentType}
		return http.StatusOK, responseHeader, content, nil
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
	responseHeader := map[string]string{"Content-Type": contentType}
	return http.StatusOK, responseHeader, htmlContent, nil
}

func (r *resourceRequestProcessor) getRequestBody(req *http.Request) (requestBody []byte, tooLarge bool, err os.Error) {
	buf := make([]byte, 4096)
	size, err := req.Body.Read(buf)
	if err != nil {
		return nil, false, err
	}
	if _, err := req.Body.Read(make([]byte, 1)); err != os.EOF {
		return nil, true, nil
	}
	return buf[:size], false, nil
}

func (r *resourceRequestProcessor) getScheme() string {
	return "http" // TODO: https?
}

func (r *resourceRequestProcessor) DoPost(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isPostablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
		return http.StatusUnsupportedMediaType, nil, nil, nil
	}
	acceptHeader := req.Header.Get("Accept")
	if jsonQVal := checkAcceptHeader("application/json", acceptHeader); jsonQVal == 0 {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	requestBody, tooLarge, err := r.getRequestBody(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if tooLarge {
		return http.StatusRequestEntityTooLarge, nil, nil, err		
	}
	obj := map[string]string{}
	if err := json.Unmarshal(requestBody, &obj); err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	newID, err := r.ResourceStorage.Create(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if newID == 0 {
		return http.StatusNotFound, nil, nil, nil
	}
	newPath := path + "/" + strconv.Uitoa64(newID)
	if err := r.ResourceStorage.Set(newPath, obj); err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	// TODO: content
	responseHeader := map[string]string{
		"Content-Type": "application/xhtml+xml; charset=utf-8",
		"Location": r.getScheme() + "://" + req.Host + newPath}
	responseContent, err := json.Marshal(obj)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	return http.StatusCreated, responseHeader, responseContent, nil
}

func (r *resourceRequestProcessor) DoPut(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isPuttablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
		return http.StatusUnsupportedMediaType, nil, nil, nil
	}
	acceptHeader := req.Header.Get("Accept")
	if jsonQVal := checkAcceptHeader("application/json", acceptHeader); jsonQVal == 0 {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	requestBody, tooLarge, err := r.getRequestBody(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if tooLarge {
		return http.StatusRequestEntityTooLarge, nil, nil, err		
	}
	obj := map[string]string{}
	if err := json.Unmarshal(requestBody, &obj); err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	oldObj, err := r.ResourceStorage.Get(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if err := r.ResourceStorage.Set(path, obj); err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	responseHeader := map[string]string{
		"Content-Type": "application/xhtml+xml; charset=utf-8"}
	responseContent, err := json.Marshal(obj)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if oldObj != nil {
		return http.StatusOK, responseHeader, responseContent, nil
	}
	responseHeader["Location"] = r.getScheme() + "://" + req.Host + req.URL.Path
	return http.StatusCreated, responseHeader, responseContent, nil
}

func (r *resourceRequestProcessor) DoDelete(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !isDeletablePath(path) {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}
	// TODO: 子リソースの再帰的削除
	return http.StatusNoContent, nil, nil, nil
}
