# go-frappe-client

## Installation

`go get github.com/joeirimpan/go-frappe-client`


## Usage

```golang
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	frappe "github.com/joeirimpan/go-frappe-client"
)

// SampleResp holder struct
type SampleResp struct {
	Message string   `json:"message"`
}

func main() {
	frappeClient, _ := frappe.New("http://localhost:5001/", "username", "password", true)

	// Creating a post request
	s := SampleResp{}
	params := url.Values{}
	params.Set("param1", "value1")
	resp, _ := frappeClient.Do(http.MethodPost, "module.module.my_api_method", params, nil)
	if err := json.Unmarshal(resp.Body, &s); err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)

	// Creating a post request and serialize back to struct
	r := SampleResp{}
	params = url.Values{}
	params.Set("param1", "value1")
	frappeClient.DoJSON(http.MethodPost, "module.module.my_api_method", params, nil, &r)
	fmt.Println(r)
}
```
