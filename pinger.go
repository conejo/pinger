package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type config struct {
	AccountSid string `json:"accountSid,omitempty"`
	AuthToken  string `json:"authToken,omitempty"`
	FromPhone  string `json:"fromPhone,omitempty"`
	ToPhone    string `json:"toPhone,omitempty"`
	SmsURL     string `json:"smsUrl,omitempty"`
}

/*
{
    "sid":"",
    "date_created":"Wed, 10 Oct 2012 17:57:42 +0000",
    "date_updated":"Wed, 10 Oct 2012 17:57:42 +0000",
    "date_sent":null,
    "account_sid":"",
    "to":"+14108675309",
    "from":"+15005550006",
    "body":"hi",
    "status":"queued",
    "direction":"outbound-api",
    "api_version":"2010-04-01",
    "price":null,
    "uri":""
}
*/
type smsResponse struct {
	Sid        string `json:"sid,omitempty"`
	AccountSid string `json:"accountSid,omitempty"`
	To         string `json:"to,omitempty"`
	From       string `json:"from,omitempty"`
	Status     string `json:"status,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Price      string `json:"price,omitempty"`
}

func main() {

	var configFile string
	flag.StringVar(&configFile, "configFile", "settings.json", "Specify the location of the config.json file.")
	flag.Parse()

	log.Printf("Using config values in: %s", configFile)

	conf := config{}
	readConfig(configFile, &conf)

	hostname := getHostname()
	addresses := getIPAddresses()

	log.Printf("Hostname: %s", hostname)
	log.Printf("Addresses: %s", strings.Join(addresses, ", "))

	smsResp := sendSms(&conf, hostname, addresses)
	log.Printf("Finished. Status: %s", smsResp.Status)
}

func sendSms(conf *config, hostname string, addresses []string) *smsResponse {
	var smsResp smsResponse

	requrl := fmt.Sprintf(conf.SmsURL, conf.AccountSid)

	body := url.Values{}
	body.Set("To", conf.ToPhone)
	body.Set("From", conf.FromPhone)
	body.Set("Body", fmt.Sprintf("Hostname : %s\nAddresses: %s", hostname, strings.Join(addresses, ", ")))

	client := &http.Client{}

	br := *strings.NewReader(body.Encode())

	req, err := http.NewRequest(http.MethodPost, requrl, &br)
	if err != nil {
		log.Fatalf("Error building http request: %s", err)
	}

	req.SetBasicAuth(conf.AccountSid, conf.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making http request: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("received unsuccessful http statusCode: %d. \n%s", resp.StatusCode, string(respBody))
	}

	err = json.Unmarshal(respBody, &smsResp)
	if err != nil {
		log.Fatalf("received error unmarshaling json response. err: %s\n%s", err, string(respBody))
	}

	return &smsResp
}

func readConfig(filename string, conf *config) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	err = json.Unmarshal(c, conf)
	if err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}
}

// getHostname returns the machines hostname. It exits if an error is encountered.
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to determine hostname. quitting. %s", err)
	}

	return hostname
}

// getIPAddresses returns the configured ipaddresses on the machine.
// It exits if an error is encountered.
// adapted from: https://gist.github.com/jniltinho/9787946
func getIPAddresses() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Unable to get ip address. quitting. %s", err)
	}

	var addresses []string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}

	return addresses
}
