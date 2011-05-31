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

type ResourceStorage interface {
	Get(urlPath string) (map[string]string, os.Error)
	GetChildren(urlPath string) (map[string]map[string]string, os.Error)
	Set(urlPath string, obj map[string]string) os.Error
	Delete(urlPath string) (bool, os.Error)
	Create(urlPath string) (uint64, os.Error)
}

type resourceRequestProcessor struct {
	ResourceStorage
}

func NewResourceRequestProcessor(resourceStorage ResourceStorage) RequestProcessor {
	return &resourceRequestProcessor{resourceStorage}
}

func (r *resourceRequestProcessor) isPostablePath(path string) bool {
	var pathRegExp = regexp.MustCompile(`^/games$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func (r *resourceRequestProcessor) isPuttablePath(path string) bool {
	pathRegExp := regexp.MustCompile(`^/games/[a-zA-Z0-9_\-]+(/[a-zA-Z0-9_\-]+/[a-zA-Z0-9_\-]+)*$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func (r *resourceRequestProcessor) isDeletablePath(path string) bool {
	pathRegExp := regexp.MustCompile(`^/games/[a-zA-Z0-9_\-]+$`)
	if pathRegExp.MatchString(path) {
		return true
	}
	return false
}

func (r *resourceRequestProcessor) getAllowHeader(path string) string {
	allow := "OPTIONS, GET, HEAD"
	if r.isPostablePath(path) {
		allow += ", POST"
	}
	if r.isPuttablePath(path) {
		allow += ", PUT"
	}
	if r.isDeletablePath(path) {
		allow += ", DELETE"
	}
	return allow
}

func (r *resourceRequestProcessor) DoOptions(req *http.Request) (int, map[string]string, os.Error) {
	responseHeader := map[string]string{"Content-Length": "0"}
	path := req.URL.Path
	if path == "*" {
		responseHeader["Allow"] = "OPTIONS, GET, HEAD, POST, PUT, DELETE"
	} else {
		responseHeader["Allow"] = r.getAllowHeader(path)
	}
	return http.StatusOK, responseHeader, nil
}

func (r *resourceRequestProcessor) returnsFile(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	fileContent, err := GetFileFromCache(filepath.Join("public", path))
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if fileContent == nil {
		return http.StatusNotFound, nil, nil, nil
	}
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
	status := http.StatusOK
	if checkAcceptHeader(contentType, req.Header.Get("Accept")) == 0 {
		status  = http.StatusNotAcceptable
	}
	return status, map[string]string{"Content-Type":contentType}, fileContent, nil
}

func (r *resourceRequestProcessor) returnsResource(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	resourceObj, err := r.ResourceStorage.Get(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if len(path) == 0 {
		return http.StatusBadRequest, nil, nil, nil
	}
	// remove the last slash
	if 1 < len(path) && path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	if path == "/" {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	slashCount := strings.Count(path, "/")
	if slashCount == 0 {
		return http.StatusBadRequest, nil, nil, nil
	}
	if slashCount == 1 {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	var obj interface{}
	switch slashCount % 2 {
	case 0:
		if resourceObj == nil {
			return http.StatusNotFound, nil, nil, err
		}
		obj = resourceObj
	case 1:
		// resourceObj may be nil
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

func (r *resourceRequestProcessor) returnsHTMLFile(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	resourceObj, err := r.ResourceStorage.Get(path)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if resourceObj == nil {
		return http.StatusNotFound, nil, nil, err
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

func (r *resourceRequestProcessor) DoGet(req *http.Request) (int, map[string]string, []byte, os.Error) {
	if status, header, content, err := r.returnsFile(req); status != http.StatusNotFound || err != nil {
		return status, header, content, err
	}
	acceptHeader := req.Header.Get("Accept")
	jsonQVal := checkAcceptHeader("application/json", acceptHeader)
	xhtmlQVal := checkAcceptHeader("application/xhtml+xml", acceptHeader)
	htmlQVal := checkAcceptHeader("text/html", acceptHeader)
	if jsonQVal == 0 && xhtmlQVal == 0 && htmlQVal == 0 {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	if xhtmlQVal <= xhtmlQVal && htmlQVal <= jsonQVal {
		return r.returnsResource(req)
	}
	return r.returnsHTMLFile(req)
}

func (r *resourceRequestProcessor) DoHead(req *http.Request) (int, map[string]string, os.Error) {
	status, responseHeader, _, err := r.DoGet(req)
	return status, responseHeader, err
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
	if !r.isPostablePath(path) {
		responseHeader := map[string]string{"Allow":r.getAllowHeader(req.URL.Path)}
		return http.StatusMethodNotAllowed, responseHeader, nil, nil
	}
	if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
		return http.StatusUnsupportedMediaType, nil, nil, nil
	}
	acceptHeader := req.Header.Get("Accept")
	if checkAcceptHeader("application/json", acceptHeader) == 0 {
		return http.StatusNotAcceptable, nil, nil, nil
	}
	requestBody, tooLarge, err := r.getRequestBody(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	if tooLarge {
		return http.StatusRequestEntityTooLarge, nil, nil, nil
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
	responseContent, err := json.Marshal(obj)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	responseHeader := map[string]string{
		"Content-Type": "application/xhtml+xml; charset=utf-8",
		"Location": r.getScheme() + "://" + req.Host + newPath}
	return http.StatusCreated, responseHeader, responseContent, nil
}

func (r *resourceRequestProcessor) DoPut(req *http.Request) (int, map[string]string, []byte, os.Error) {
	path := req.URL.Path
	if !r.isPuttablePath(path) {
		responseHeader := map[string]string{"Allow":r.getAllowHeader(req.URL.Path)}
		return http.StatusMethodNotAllowed, responseHeader, nil, nil
	}
	if !regexp.MustCompile("^application/json;?").MatchString(req.Header.Get("Content-Type")) {
		return http.StatusUnsupportedMediaType, nil, nil, nil
	}
	acceptHeader := req.Header.Get("Accept")
	if checkAcceptHeader("application/json", acceptHeader) == 0 {
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
	if !r.isDeletablePath(path) {
		responseHeader := map[string]string{"Allow":r.getAllowHeader(req.URL.Path)}
		return http.StatusMethodNotAllowed, responseHeader, nil, nil
	}
	// TODO: 子リソースの再帰的削除
	return http.StatusNoContent, nil, nil, nil
}
