package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ApiKey = "AIzaSyDaB7dhDz0ulPnUhOaCtjDbHrj5nccPStE"
)

func main() {
	url := fmt.Sprintf("https://safebrowsing.googleapis.com/v4/threatMatches:find?key=%s", ApiKey)
	//fmt.Println(url)

	var data = []byte(`{"client": {"clientId": "drixevel", "clientVersion": "1.0.0"}, "threatInfo": {"threatTypes": ["THREAT_TYPE_UNSPECIFIED", "MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"], "platformTypes": ["WINDOWS"], "threatEntryTypes": ["URL"], "threatEntries": [{"url": "https://getyourpaymentirs.com/form/data"}]}}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error while scanning malicious urls from Google: Invalid Status Code (%d returned, should be 200)", resp.StatusCode)
	}

	sb := string(body)
	fmt.Println(sb)

	if len(sb) < 4 {
		fmt.Println("All urls are clean.")
	} else {
		fmt.Println("Some or all urls are malicious.")
	}

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// fmt.Println("response Body:", string(body))
}
