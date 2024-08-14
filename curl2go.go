package curl2go

import (
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Request represents a parsed curl command
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Cookies map[string]string
	Body    string
}

// Execute sends the request and returns the response body as a byte slice
func (r *Request) Execute() ([]byte, error) {
	req, err := r.ToHTTPRequest()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*ParseCurlCommandFile takes a filename string to a text file containing the curl command and returns a Request struct*/
func ParseCurlCommandFile(filename string) (*Request, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseCurlCommand(string(f))
}

// ParseCurlCommand takes a curl command string and returns a Request struct
func ParseCurlCommand(curlCommand string) (*Request, error) {

	if curlCommand == "" {
		return nil, errors.New("empty curl command")
	}
	req := &Request{
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
	}

	// Remove 'curl' from the beginning of the command if present
	curlCommand = strings.TrimPrefix(curlCommand, "curl")
	curlCommand = strings.TrimSpace(curlCommand)

	// Replace newlines and multiple spaces with a single space
	curlCommand = regexp.MustCompile(`\s+`).ReplaceAllString(curlCommand, " ")

	// Parse URL
	urlRegex := regexp.MustCompile(`'(.*?)'`)
	urlMatches := urlRegex.FindAllStringSubmatch(curlCommand, -1)
	if len(urlMatches) == 0 {
		return nil, errors.New("URL not found in curl command")
	}
	req.URL = urlMatches[0][1]

	// Parse method (default to GET if not specified)
	req.Method = "GET"
	if strings.Contains(curlCommand, "-X") || strings.Contains(curlCommand, "--request") {
		methodRegex := regexp.MustCompile(`-X\s+(\w+)`)
		methodMatch := methodRegex.FindStringSubmatch(curlCommand)
		if len(methodMatch) >= 2 {
			req.Method = methodMatch[1]
		}
	}

	// Parse headers
	headerRegex := regexp.MustCompile(`-H\s+'([^:]+):\s*([^']+)'`)
	headerMatches := headerRegex.FindAllStringSubmatch(curlCommand, -1)
	for _, match := range headerMatches {
		key := match[1]
		value := match[2]
		if strings.ToLower(key) == "cookie" {
			parseCookies(value, req.Cookies)
		} else {
			req.Headers[key] = value
		}
	}

	// Parse body
	bodyRegex := regexp.MustCompile(`--data\s+'([^']+)'`)
	bodyMatch := bodyRegex.FindStringSubmatch(curlCommand)
	if len(bodyMatch) >= 2 {
		req.Body = bodyMatch[1]
	}

	return req, nil
}

// parseCookies helper function to parse the Cookie header
func parseCookies(cookieHeader string, cookieMap map[string]string) {
	cookies := strings.Split(cookieHeader, ";")
	for _, cookie := range cookies {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 {
			cookieMap[parts[0]] = parts[1]
		}
	}
}

// ToHTTPRequest converts the parsed Request to a *http.Request
func (r *Request) ToHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	for key, value := range r.Cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	return req, nil
}
