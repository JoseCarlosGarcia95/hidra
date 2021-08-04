// HTTP scenario do test using HTTP protocol
package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"

	"github.com/hidracloud/hidra/models"
	"github.com/hidracloud/hidra/scenarios"
)

// Represent an http scenario
type HttpScenario struct {
	models.Scenario

	Url      string
	Method   string
	Response *http.Response
	Body     string
	Redirect string
	Headers  map[string]string
	Client   *http.Client
}

// Set user agent
func (h *HttpScenario) setUserAgent(c map[string]string) ([]models.CustomMetric, error) {
	var ok bool
	if h.Headers["User-Agent"], ok = c["user-agent"]; !ok {
		return nil, fmt.Errorf("user-agent parameter missing")
	}
	return nil, nil
}

// Add new HTTP header
func (h *HttpScenario) addHttpHeader(c map[string]string) ([]models.CustomMetric, error) {
	if _, ok := c["key"]; !ok {
		return nil, fmt.Errorf("key parameter missing")
	}
	if _, ok := c["value"]; !ok {
		return nil, fmt.Errorf("value parameter missing")
	}

	h.Headers[c["key"]] = c["value"]
	return nil, nil
}

// Send a request depends of the method
func (h *HttpScenario) requestByMethod(c map[string]string) ([]models.CustomMetric, error) {
	var err error

	body := ""

	if _, ok := c["body"]; ok {
		body = c["body"]
	}

	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	h.Client.Jar = jar

	req, err := http.NewRequest(h.Method, h.Url, strings.NewReader(body))

	if err != nil {
		return nil, err
	}

	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.Client.Do(req)

	if err != nil {
		return nil, err
	}

	h.Response = resp

	b, err := ioutil.ReadAll(h.Response.Body)

	if err != nil {
		return nil, err
	}

	h.Body = strings.ToLower(string(b))
	h.Response.Body.Close()

	return nil, err
}

// Make http request to given URL
func (h *HttpScenario) request(c map[string]string) ([]models.CustomMetric, error) {
	var err error
	var ok bool

	if h.Url, ok = c["url"]; !ok {
		return nil, fmt.Errorf("url parameter missing")
	}

	h.Method = "GET"

	if _, ok = c["method"]; ok {
		h.Method = strings.ToUpper(c["method"])
	}

	_, err = h.requestByMethod(c)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Check if status code match
func (h *HttpScenario) statusCodeShouldBe(c map[string]string) ([]models.CustomMetric, error) {
	if h.Response == nil {
		return nil, fmt.Errorf("request should be initialized first")
	}

	if _, ok := c["statusCode"]; !ok {
		return nil, fmt.Errorf("statusCode parameter missing")
	}

	if strconv.Itoa(h.Response.StatusCode) != c["statusCode"] {
		return nil, fmt.Errorf("statusCode expected %s, but %d", c["statusCode"], h.Response.StatusCode)
	}

	return nil, nil
}

func (h *HttpScenario) bodyShouldContain(c map[string]string) ([]models.CustomMetric, error) {
	if _, ok := c["search"]; !ok {
		return nil, fmt.Errorf("search parameter missing")
	}

	if !strings.Contains(h.Body, strings.ToLower(c["search"])) {
		return nil, fmt.Errorf("expected %s in body, but not found", c["search"])
	}

	return nil, nil
}

func (h *HttpScenario) shouldRedirectTo(c map[string]string) ([]models.CustomMetric, error) {
	if _, ok := c["url"]; !ok {
		return nil, fmt.Errorf("url parameter missing")
	}

	if h.Response.Header.Get("Location") != c["url"] {
		return nil, fmt.Errorf("expected redirect to %s, but got %s", c["url"], h.Response.Header.Get("Location"))
	}

	return nil, nil
}

// Clear parameters
func (h *HttpScenario) clear(c map[string]string) ([]models.CustomMetric, error) {
	h.Url = ""
	h.Response = nil
	h.Method = ""
	h.Headers = make(map[string]string)

	return nil, nil
}

func (h *HttpScenario) Init() {
	h.StartPrimitives()

	h.Headers = make(map[string]string)
	h.Headers["User-Agent"] = "hidra/monitoring"

	h.Client = &http.Client{}

	h.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	h.RegisterStep("request", h.request)
	h.RegisterStep("statusCodeShouldBe", h.statusCodeShouldBe)
	h.RegisterStep("setUserAgent", h.setUserAgent)
	h.RegisterStep("addHttpHeader", h.addHttpHeader)
	h.RegisterStep("bodyShouldContain", h.bodyShouldContain)
	h.RegisterStep("shouldRedirectTo", h.shouldRedirectTo)

	h.RegisterStep("clear", h.clear)
}

func init() {
	scenarios.Add("http", func() models.IScenario {
		return &HttpScenario{}
	})
}
