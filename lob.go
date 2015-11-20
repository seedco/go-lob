package lob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"github.com/seedco/go-logging"
)

var log = logging.MustGetLogger("lob")

// LogStackTrace logs a stack trace for the given error.
func logStackTrace(err error) {
	buf := make([]byte, 0, 16384)
	n := runtime.Stack(buf, false)
	if err != nil {
		log.Errorf("Non-nil error %s; stack trace %s", err.Error(), buf[:n])
	} else {
		log.Errorf("Nil error; stack trace %s", buf[:n])
	}
}

// Lob represents information on how to connect to the lob.com API.
type Lob struct {
	BaseAPI   string
	APIKey    string
	UserAgent string
}

// Base URL and API version for Lob.
const (
	BaseAPI    = "https://api.lob.com/v1/"
	APIVersion = "2014-12-18"
)

// MetricsSet is the bundle of metrics associated
// with each Lob method.
type MetricsSet struct {
	CreateCheck       *MetricsBundle
	GetCheck          *MetricsBundle
	ListChecks        *MetricsBundle
	CreateBankAccount *MetricsBundle
	GetBankAccount    *MetricsBundle
	ListBankAccounts  *MetricsBundle
	VerifyAddress     *MetricsBundle
	CreateAddress     *MetricsBundle
	GetAddress        *MetricsBundle
	DeleteAddress     *MetricsBundle
	ListAddresses     *MetricsBundle
	GetStates         *MetricsBundle
	GetCountries      *MetricsBundle
}

// Metrics is the set of metrics for this API.
// It is shared across all instances.
var Metrics *MetricsSet

func init() {
	Metrics = &MetricsSet{
		CreateCheck:       NewMetricsBundle("check_create"),
		GetCheck:          NewMetricsBundle("check_get"),
		ListChecks:        NewMetricsBundle("check_list"),
		CreateBankAccount: NewMetricsBundle("bank_account_create"),
		GetBankAccount:    NewMetricsBundle("bank_account_get"),
		ListBankAccounts:  NewMetricsBundle("bank_account_list"),
		VerifyAddress:     NewMetricsBundle("address_verify"),
		CreateAddress:     NewMetricsBundle("address_create"),
		GetAddress:        NewMetricsBundle("address_get"),
		DeleteAddress:     NewMetricsBundle("address_delete"),
		ListAddresses:     NewMetricsBundle("address_list"),
		GetStates:         NewMetricsBundle("states_list"),
		GetCountries:      NewMetricsBundle("countries_list"),
	}
}

// NewLob creates an object that can be used to connect to the lob.com API.
func NewLob(baseAPI, apiKey, userAgent string) *Lob {
	return &Lob{
		BaseAPI:   baseAPI,
		APIKey:    apiKey,
		UserAgent: userAgent,
	}
}

func queryParams(params map[string]string) string {
	if params == nil {
		return ""
	}
	pieces := make([]string, 0, len(params))
	for k, v := range params {
		pieces = append(pieces, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	return "?" + strings.Join(pieces, "&")
}

// Use JSON tag information to create a form values map.
func json2form(v interface{}) map[string]string {
	value := reflect.ValueOf(v)
	t := value.Type()
	params := make(map[string]string)
	for i := 0; i < value.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("json")
		fv := value.Field(i).Interface()
		if fv == nil {
			continue
		}
		switch fv.(type) {
		case string:
			if fv.(string) != "" {
				params[name] = fv.(string)
			}
		case []string:
			if len(fv.([]string)) > 0 {
				params[name] = strings.Join(fv.([]string), " ")
			}
		default:
			// ignore
			log.Debugf("Unknown field type: " + value.Field(i).Type().String())
		}
	}
	return params
}

// Get performs a GET request to the Lob API.
func (lob *Lob) Get(endpoint string, params map[string]string, returnValue interface{}) error {
	fullURL := lob.BaseAPI + endpoint + queryParams(params)
	log.Debugf("Lob GET %s", fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logStackTrace(err)
		return err
	}

	req.SetBasicAuth(lob.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", lob.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logStackTrace(err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logStackTrace(err)
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non-200 status code %d returned from %s with body %s", resp.StatusCode, fullURL, data)
		logStackTrace(err)
		json.Unmarshal(data, returnValue) // try, anyway -- in case the caller wants error info
		return err
	}

	return json.Unmarshal(data, returnValue)
}

// Post performs a POST request to the Lob API.
func (lob *Lob) Post(endpoint string, params map[string]string, returnValue interface{}) error {
	fullURL := lob.BaseAPI + endpoint
	log.Debugf("Lob POST %s", fullURL)

	var body io.Reader
	if params != nil {
		form := url.Values(make(map[string][]string))
		for k, v := range params {
			form.Add(k, v)
		}
		bodyString := form.Encode()
		body = bytes.NewBuffer([]byte(bodyString))
	}

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		logStackTrace(err)
		return err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.SetBasicAuth(lob.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", lob.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logStackTrace(err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logStackTrace(err)
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non-200 status code %d returned from %s with body %s", resp.StatusCode, fullURL, data)
		logStackTrace(err)
		json.Unmarshal(data, returnValue) // try, anyway -- in case the caller wants error info
		return err
	}

	return json.Unmarshal(data, returnValue)
}

// Delete performs a DELETE request to the Lob API.
func (lob *Lob) Delete(endpoint string, returnValue interface{}) error {
	fullURL := lob.BaseAPI + endpoint
	log.Debugf("Lob DELETE %s", fullURL)

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		logStackTrace(err)
		return err
	}

	req.SetBasicAuth(lob.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", lob.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logStackTrace(err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logStackTrace(err)
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non-200 status code %d returned from %s with body %s", resp.StatusCode, fullURL, data)
		logStackTrace(err)
		json.Unmarshal(data, returnValue) // try, anyway -- in case the caller wants error info
		return err
	}

	return json.Unmarshal(data, returnValue)
}
