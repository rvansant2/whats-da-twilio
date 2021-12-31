package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	viper "github.com/spf13/viper"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say  string `xml:",omitempty"`
	Play string `xml:",omitempty"`
}

func init() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()

	viper.SetDefault("app.twilio.toPhoneNumber", viper.GetString("app.twilio.defaultToNumber"))

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Fatal error config file: default \n", err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/twiml", twiml)
	http.HandleFunc("/call", call)
	http.HandleFunc("/sms", sms)
	http.ListenAndServe(":3000", nil)
}

func twiml(w http.ResponseWriter, r *http.Request) {
	//twiml := TwiML{Say: "Thanks for trying our documentation. Enjoy!"}
	twiml := TwiML{Say: "Thanks for trying our documentation. Enjoy!", Play: "http://demo.twilio.com/docs/classic.mp3"}

	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func call(w http.ResponseWriter, r *http.Request) {
	// Let's set some initial default variables
	accountSid := viper.GetString("app.twilio.accountSid")
	authToken := viper.GetString("app.twilio.authToken")
	urlStr := viper.GetString("app.twilio.urlStr") + accountSid + "/Calls.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", "+"+viper.GetString("app.twilio.toPhoneNumber"))
	v.Set("From", "+"+viper.GetString("app.twilio.fromNumber"))
	v.Set("Url", viper.GetString("app.twilio.appUrl"))
	rb := *strings.NewReader(v.Encode())

	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println("Status code: " + resp.Status)
		w.Write([]byte("Hello What's da Twilio!"))
	}
}

func sms(w http.ResponseWriter, r *http.Request) {
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: viper.GetString("app.twilio.accountSid"),
		Password: viper.GetString("app.twilio.authToken"),
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(viper.GetString("app.twilio.toPhoneNumber"))
	params.SetFrom(viper.GetString("app.twilio.fromNumber"))
	params.SetBody(viper.GetString("app.twilio.defaultTextMessage"))

	_, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("SMS sent successfully!"))
	}
}
