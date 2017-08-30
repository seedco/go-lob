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
	"strconv"
	"strings"

	"github.com/op/go-logging"
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

type Lob interface {
	// Checks
	CreateCheck(*CreateCheckRequest) (*Check, error)
	GetCheck(string) (*Check, error)
	CancelCheck(string) (*CancelCheckResponse, error)
	ListChecks(int, int) (*ListChecksResponse, error)
	// Addresses
	CreateAddress(*Address) (*Address, error)
	GetAddress(string) (*Address, error)
	DeleteAddress(string) error
	ListAddresses(int, int) (*ListAddressesResponse, error)
	VerifyAddress(*Address) (*AddressVerificationResponse, error)
	// NamedObject
	GetStates() (*NamedObjectList, error)
	GetCountries() (*NamedObjectList, error)
	// Bank Accounts
	CreateBankAccount(*CreateBankAccountRequest) (*BankAccount, error)
	GetBankAccount(string) (*BankAccount, error)
	ListBankAccounts(int, int) (*ListBankAccountsResponse, error)
}

// Lob represents information on how to connect to the lob.com API.
type lob struct {
	BaseAPI   string
	APIKey    string
	UserAgent string
}

// Base URL and API version for Lob.
const (
	BaseAPI    = "https://api.lob.com/v1/"
	APIVersion = "2016-01-19"
)

// NewLob creates an object that can be used to connect to the lob.com API.
func NewLob(baseAPI, apiKey, userAgent string) *lob {
	return &lob{
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
		switch x := fv.(type) {
		case *string:
			if x != nil {
				params[name] = *x
			}
		case string:
			if x != "" {
				params[name] = x
			}
		case int:
			if x != 0 {
				params[name] = strconv.Itoa(x)
			}
		case *bool:
			if x != nil {
				params[name] = fmt.Sprintf("%v", *x)
			}
		case int64:
			if x != 0 {
				params[name] = strconv.FormatInt(x, 10)
			}
		case float64:
			params[name] = fmt.Sprintf("%.2f", x)
		case []string:
			if len(x) > 0 {
				params[name] = strings.Join(x, " ")
			}
		case map[string]string:
			for mapkey, mapvalue := range x {
				params[name+"["+mapkey+"]"] = mapvalue
			}
		default:
			// ignore
			panic(fmt.Errorf("Unknown field type: " + value.Field(i).Type().String()))
		}
	}
	return params
}

// Get performs a GET request to the lob API.
func (l *lob) get(endpoint string, params map[string]string, returnValue interface{}) error {
	fullURL := l.BaseAPI + endpoint + queryParams(params)
	log.Debugf("Lob GET %s", fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logStackTrace(err)
		return err
	}

	req.SetBasicAuth(l.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", l.UserAgent)

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
func (l *lob) post(endpoint string, params map[string]string, returnValue interface{}) error {
	fullURL := l.BaseAPI + endpoint
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

	req.SetBasicAuth(l.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", l.UserAgent)

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
func (l *lob) delete(endpoint string, returnValue interface{}) error {
	fullURL := l.BaseAPI + endpoint
	log.Debugf("Lob DELETE %s", fullURL)

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		logStackTrace(err)
		return err
	}

	req.SetBasicAuth(l.APIKey, "")
	req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", l.UserAgent)

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
