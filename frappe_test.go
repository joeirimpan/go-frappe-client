package frappe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

// SampleResp holder struct
type SampleResp struct {
	Message string `json:"message"`
}

func TestFrappeClient(t *testing.T) {
	auth := LoginAuth{
		userName: "username",
		password: "password",
	}
	frappeClient, _ := New("http://localhost:5001/", &auth, true)

	// Creating a post request
	s := SampleResp{}
	params := url.Values{}
	params.Set("param1", "value1")
	resp, _ := frappeClient.Do(http.MethodPost, "erpnext.accounts.doctype.subscription.subscription.ping", params, nil)
	if err := json.Unmarshal(resp.Body, &s); err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}
