package gobotween

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type configuration struct {
	Secret       string `json:"SLACK_SIGNING_SECRET"`
	WebHookURL   string `json:"WEBHOOK_URL"`
	SlashHandler string `json:"SLASH_HANDLER"`
}

var (
	config     *configuration
	urlTarget  *url.URL
	httpClient http.Client
)

func readConfig() {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}

	d := json.NewDecoder(cfgFile)
	config = &configuration{}
	if err = d.Decode(config); err != nil {
		log.Fatalf("Decode: %v", err)
	}

}

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
)

func Gobotween(w http.ResponseWriter, r *http.Request) {
	// Send an empty string back to Slack immediately to acknowledge
	// receipt of the request.
	w.Write([]byte(""))

	if config == nil {
		readConfig()
		urlTarget, _ = url.Parse(config.SlashHandler)
		httpClient = http.Client{}
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %v", err)
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}

	result, err := verifyWebHook(r, bodyBytes, config.Secret)
	if err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	if !result {
		log.Fatalf("signatures did not match.")
	}
	forwardRequest(r, bodyBytes)

}

func forwardRequest(r *http.Request, body []byte) {

	bodyString := string(body)
	bodyString = bodyString + "&webhook_url=" + config.WebHookURL

	fwdReq, err := http.NewRequest(r.Method, urlTarget.String(), bytes.NewReader([]byte(bodyString)))

	fwdReq.Header = make(http.Header)
	token, err := getOAuthToken()
	if err != nil {
		log.Fatalf("No OAuth Token: %v", err)
	}
	fwdReq.Header.Add("Authorization", "bearer "+token)
	fwdReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(fwdReq)
	if err != nil {
		log.Fatalf(err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

}

func getOAuthToken() (string, error) {

	targetFunctionURL := config.SlashHandler
	metadataServerTokenURL := "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience="

	tokenRequestURL := metadataServerTokenURL + targetFunctionURL

	authReq, err := http.NewRequest("GET", tokenRequestURL, nil)
	authReq.Header = make(http.Header)
	authReq.Header.Add("Metadata-Flavor", "Google")

	resp, err := httpClient.Do(authReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	return bodyString, nil

}

// veryfyWebhook uses signature verification instead of tokens.
// code below this point based on https://github.com/GoogleCloudPlatform/golang-samples/blob/master/functions/slack/search.go
func verifyWebHook(r *http.Request, body []byte, slackSigningSecret string) (bool, error) {

	timeStamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	t, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %v", timeStamp, err)
	}

	if ageOk, age := checkTimestamp(t); !ageOk {
		//return false, &oldTimeStampError{fmt.Sprintf("checkTimestamp(%v): %v %v", t, ageOk, age)}
		return false, fmt.Errorf("checkTimestamp(%v): %v %v", t, ageOk, age)
	}

	if timeStamp == "" || slackSignature == "" {
		return false, fmt.Errorf("either timeStamp or signature headers were blank")
	}

	baseString := fmt.Sprintf("%s:%s:%s", version, timeStamp, body)

	signature := getSignature([]byte(baseString), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)

	if err != nil {
		return false, fmt.Errorf("hex.DecodeString(%v): %v", trimmed, err)
	}

	return hmac.Equal(signature, signatureInHeader), nil
}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

// Arbitrarily trusting requests time stamped less than 5 minutes ago.
func checkTimestamp(timeStamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timeStamp, 0))

	return t.Minutes() <= 5, t
}
