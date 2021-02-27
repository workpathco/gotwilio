// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"strconv"
)

// https://www.twilio.com/docs/proxy/api/participants

// ProxyPhoneNumberRequest shape for PhoneNumber request
type ProxyPhoneNumberRequest struct {
	Sid             string // optional
	PhoneNumber     string // optional
	IsReserved      bool   //optional
}

// ProxyPhoneNumber shape for PhoneNumber
type ProxyPhoneNumber struct {
	Sid                string      `json:"sid"`
	ServiceSid         string      `json:"service_sid"`
	PhoneNumber        string      `json:"phone_number"`
	DateUpdated        time.Time   `json:"date_updated"`
	FriendlyName       interface{} `json:"friendly_name"`
	IsoCountry         string      `json:"iso_country"`
	AccountSid         string      `json:"account_sid"`
	URL                string      `json:"url"`
	DateCreated        time.Time   `json:"date_created"`
	Capabilities       string      `json:"capabilities"`
	IsReserved         string      `json:"is_reserved"`
	InUser             string      `json:"in_use"`
}

// ProxyPhoneNumberList shape for list of PhoneNumbers
type ProxyPhoneNumberList struct {
	PhoneNumbers []ProxyPhoneNumber `json:"participants"`
	Meta         Meta               `json:"meta"`
}

// AddPhoneNumber adds a PhoneNumber to a Service
func (service *ProxyService) AddPhoneNumber(req ProxyPhoneNumberRequest) (response ProxyPhoneNumber, exception *Exception, err error) {

	twilioURL := fmt.Sprintf("%s/%s/%s/%s", ProxyBaseUrl, "Services", service.Sid, "PhoneNumbers")

	res, err := service.twilio.post(PhoneNumberFormValues(req), twilioURL)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err

}

// ListPhoneNumbers reads a list of PhoneNumbers on Service
func (service *ProxyService) ListPhoneNumbers() (response []ProxyPhoneNumber, exception *Exception, err error) {

	twilioURL := fmt.Sprintf("%s/%s/%s/%s", ProxyBaseUrl, "Services", service.Sid, "PhoneNumbers")

	res, err := service.twilio.get(twilioURL)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	list := ProxyPhoneNumberList{}
	err = json.Unmarshal(responseBody, &list)
	return list.PhoneNumbers, exception, err

}

// GetPhoneNumber finds a PhoneNumber by SID
func (service *ProxyService) GetPhoneNumber(phoneNumberSid string) (response ProxyPhoneNumber, exception *Exception, err error) {

	twilioURL := fmt.Sprintf("%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", service.Sid, "PhoneNumbers", phoneNumberSid)

	res, err := service.twilio.get(twilioURL)
	if err != nil {
		return response, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return response, exception, err
	}

	err = json.Unmarshal(responseBody, &response)
	return response, exception, err

}

// DeletePhoneNumber deletes PhoneNumber from ProxyService
func (service *ProxyService) DeletePhoneNumber(phoneNumberSid string) (exception *Exception, err error) {

	twilioURL := fmt.Sprintf("%s/%s/%s/%s/%s", ProxyBaseUrl, "Services", service.Sid, "PhoneNumbers", phoneNumberSid)

	res, err := service.twilio.delete(twilioURL)
	if err != nil {
		return exception, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, nil
}

// PhoneNumberFormValues - Form values initialization
func PhoneNumberFormValues(req ProxyPhoneNumberRequest) url.Values {
	formValues := url.Values{}

	if req.Sid != "" {
		formValues.Set("Sid", req.Sid)
	}

	if req.PhoneNumber != "" {
		formValues.Set("PhoneNumber", req.PhoneNumber)
	}

	formValues.Set("IsReserved", strconv.FormatBool(req.IsReserved))

	return formValues
}
