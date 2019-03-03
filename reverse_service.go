package reverseproxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/mallbook/commandline"
)

var (
	logger = log.New(os.Stdout, "reverseproxy", log.Lshortfile)
)

type route struct {
	SubPath    string `josn:"subPath"`
	HTTPMethod string `json:"httpMethod"`
}

type reverseConfig struct {
	RootPath   string   `json:"rootPath"`
	TargetPath string   `json:"targetPath"`
	ProxyPass  []string `json:"proxyPass"`
	Routes     []route  `json:"routes"`
}

type config struct {
	ReverseProxy map[string]reverseConfig `json:"reverseProxy"`
}

func loadConfig(fileName string) (conf config, err error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Println("ReadFile: ", err.Error())
		return
	}

	if err = json.Unmarshal(bytes, &conf); err != nil {
		logger.Println("Unmarshal: ", err.Error())
		return
	}

	return
}

type reverseService struct {
	config reverseConfig
}

func newReverseService(config reverseConfig) *restful.WebService {
	s := &reverseService{
		config: config,
	}

	ws := new(restful.WebService)

	ws.Path(s.config.RootPath)
	for _, route := range s.config.Routes {
		ws.Route(ws.Method(route.HTTPMethod).Path(route.SubPath).To(s.proxy))
	}

	return ws
}

func (s reverseService) proxy(req *restful.Request, resp *restful.Response) {
	target := s.config.ProxyPass[0] + s.config.TargetPath + strings.TrimPrefix(req.Request.URL.Path, s.config.RootPath)
	if req.Request.URL.RawQuery != "" {
		target += "?" + req.Request.URL.RawQuery
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		logger.Printf("Parse url fail, err = %s", err.Error())
		return
	}

	director := func(req *http.Request) {
		req.URL = targetURL
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	p := httputil.ReverseProxy{
		Director: director,
	}
	p.ServeHTTP(resp.ResponseWriter, req.Request)
}

func selectProxyPass(proxyPass []string, policy string) string {
	return proxyPass[0]
}

func init() {
	prefix := commandline.GetPrefixPath()
	configFile := prefix + "/etc/conf/reverse_proxy.json"
	config, err := loadConfig(configFile)
	if err != nil {
		logger.Printf("load config file(%s) fail, err = %s", configFile, err.Error())
		return
	}

	for _, reverseConfig := range config.ReverseProxy {
		s := newReverseService(reverseConfig)
		restful.Add(s)
	}
}
