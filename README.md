# curl2go

Parses a (single) `curl` command from Dev Tools "Copy as curl" (*nix/shell version) and turns into a Golang request.

## curl command passed as a string example

```go
package main

import (
	"fmt"
	"log"

	"github.com/Loo0D/curl2go"
)

func main() {
	parsed, err := curl2go.ParseCurlCommand(`curl 'https://ifconfig.io/' \
  -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8' \
  -H 'accept-language: en;q=0.8' \
  -H 'cache-control: max-age=0' \
  -H 'dnt: 1' \
  -H 'priority: u=0, i' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-fetch-dest: document' \
  -H 'sec-fetch-mode: navigate' \
  -H 'sec-fetch-site: none' \
  -H 'sec-fetch-user: ?1' \
  -H 'sec-gpc: 1' \
  -H 'upgrade-insecure-requests: 1' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36'
	`)

	if err != nil {
		log.Fatalf("Error parsing curl command: %v", err)
	}

	result, err := parsed.Execute()
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
	}

	fmt.Println(string(result))

}
```


## curl command taken from a file

```go
package main

import (
	"fmt"
	"log"

	"github.com/Loo0D/curl2go"
)

func main() {
	parsed, err := curl2go.ParseCurlCommandFile("./example.curl")

	if err != nil {
		log.Fatalf("Error parsing curl command: %v", err)
	}

	result, err := parsed.Execute()
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
	}

	fmt.Println(string(result))

}
```