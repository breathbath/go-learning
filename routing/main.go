package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type Hash struct {
	Scheme     string     `json:"scheme"`
	Host       string     `json:"host"`
	Path       string     `json:"path"`
	RawPath    string     `json:"-"`
	Parameters posParams `json:"parameters"`
}

func NewHash() Hash {
	return Hash{Parameters: []posParam{}}
}

func (h Hash) IsEmpty() bool {
	return h.Scheme == "" && h.Host == "" && h.Path == "" && len(h.Parameters) == 0
}

type Handler func(h Hash) (string, error)

func JsonHandler(h Hash) (string, error) {
	if h.IsEmpty() {
		return "{}", nil
	}

	b, e := json.Marshal(h)
	return string(b), e
}

type Response struct {
	data string
	err  error
}

type posParam struct {
	name  string
	value string
	pos   int
}

type posParams []posParam

func (pps posParams) MarshalJSON() ([]byte, error) {
	resultString := "{%s}"
	if len(pps) == 0 {
		return []byte("{}"), nil
	}

	posParamStrings := []string{}
	for _, pp := range pps {
		if pp.isEmpty() || pp.name == "" {
			continue
		}
		posParamStrings = append(posParamStrings, fmt.Sprintf(`"%s":"%s"`, pp.name, pp.value))
	}

	resultString = fmt.Sprintf(resultString, strings.Join(posParamStrings, ","))
	return []byte(resultString), nil
}

func (pp posParam) isEmpty() bool {
	return pp.name == "" && pp.value == ""
}

func (r Response) Render() {
	if r.err != nil {
		panic(r.err)
	}

	fmt.Println(r.data)
}

type Resolver struct{}

func (r Resolver) Resolve(route, url string, h Handler) Response {
	hash, err := r.parseBaseUrlParts(url)
	if err != nil {
		return Response{err: err}
	}

	var pathIsMatched bool
	pathIsMatched, hash.Parameters = r.parsePath(hash.RawPath, route)
	if !pathIsMatched {
		hash = NewHash()
	}

	output, err := h(hash)
	return Response{data: output, err: err}
}

func (r Resolver) parseBaseUrlParts(inputUrl string) (h Hash, err error) {
	var parsedUrl *url.URL
	parsedUrl, err = url.Parse(inputUrl)
	if err != nil {
		return
	}

	path := parsedUrl.Path
	if parsedUrl.RawQuery != "" {
		path += "?" + parsedUrl.RawQuery
	}

	h = Hash{
		Scheme:     parsedUrl.Scheme,
		Host:       parsedUrl.Host,
		Path:       path,
		RawPath:    parsedUrl.Path,
		Parameters: []posParam{},
	}
	return
}

func (r Resolver) parsePath(inputUrl, route string) (isMatched bool, params []posParam) {
	isMatched = false
	params = []posParam{}

	routeRegex, posParams := r.buildRouteRegex(route)
	regex := regexp.MustCompile(routeRegex)
	res := regex.FindStringSubmatch(inputUrl)

	if len(res) == 0 {
		return
	}

	isMatched = true

	for _, posParam := range posParams {
		if posParam.pos >= len(res) {
			log.Fatalf(
				"Route regex generation problem: the positional parameter %s is not found in url %s",
				posParam.name,
				inputUrl,
			)
		}
		if res[posParam.pos] != "" {
			posParam.value = res[posParam.pos]
			params = append(params, posParam)
		}
	}

	return
}

func (r Resolver) buildRouteRegex(route string) (regex string, posParams []posParam) {
	posParams = []posParam{}
	if route == "/" || route == "" {
		regex = `^/?$`
		return
	}

	regex = `^/%s/?$`
	regexBody := ""
	i := 0
	route = strings.Replace(route, "[/:", "[?:", -1)
	routeParts := strings.Split(route, "/")
	for _, part := range routeParts[1:] {
		requiredParam := r.getRequiredParam(part)

		if requiredParam != "" {
			regexBody += `/([^/]+)`
			i++
			posParams = append(posParams, posParam{name: requiredParam, pos: i})
		} else {
			regexBody += fmt.Sprintf("/%s", r.getRouteParam(part))
		}

		optionalParam := r.getOptionalParam(part)
		if optionalParam != "" {
			regexBody += `/?([^/]+)?`
			i++
			posParams = append(posParams, posParam{name: optionalParam, pos: i})
		}
	}

	regexBody = strings.TrimLeft(regexBody, "/")

	regex = fmt.Sprintf(regex, regexBody)

	return
}

func (r Resolver) getRouteParam(routePart string) string {
	regex := regexp.MustCompile(`^([^/\[]+)`)
	res := regex.FindStringSubmatch(routePart)

	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func (r Resolver) getRequiredParam(routePart string) string {
	regex := regexp.MustCompile(`^:([^/\[]+)`)
	res := regex.FindStringSubmatch(routePart)

	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func (r Resolver) getOptionalParam(routePart string) string {
	regex := regexp.MustCompile(`\[\?:([^\[]+)]$`)
	res := regex.FindStringSubmatch(routePart)

	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func main() {
	res := Resolver{}
	resp := res.Resolve("/:lang/products/:id/compare/:compareId", "https://ya.ru/en/products/418/compare/420", JsonHandler)
	resp.Render()
}
