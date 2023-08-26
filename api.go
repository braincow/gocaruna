package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	AuthorizationURL  = "https://plus.caruna.fi/api/authorization/login"
	AuthenticationURL = "https://authentication2.caruna.fi"
	PortalURL         = "https://authentication2.caruna.fi/portal"
	TokenURL          = "https://plus.caruna.fi/api/authorization/token"
	RedirectURL       = "https://plus.caruna.fi"
	CustomersURL      = "https://plus.caruna.fi/api/customers"
	LogoutURL         = "https://authentication2.caruna.fi/portal/logout"
)

type CarunaAPIClient struct {
	client http.Client
	login  LoginInfo
}

func NewCarunaClient() (*CarunaAPIClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("while building a new client: %v", err)
	}

	client := http.Client{
		Jar: jar,
	}

	return &CarunaAPIClient{client: client}, nil
}

func (api *CarunaAPIClient) LoginInfo() LoginInfo {
	return api.login
}

func (api *CarunaAPIClient) CustomerInfo(customerId string) (*CustomerInfo, error) {
	// Construct the URL
	fullURL := fmt.Sprintf("%s/%s/info", CustomersURL, customerId)

	// Prepare the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+api.login.Token)

	// Send the request using the client
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	infoBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var customerInfo CustomerInfo
	if err := json.Unmarshal(infoBody, &customerInfo); err != nil {
		return nil, err
	}

	return &customerInfo, nil
}

func (api *CarunaAPIClient) ConsumedHours(customer string, meteringPoint string, date time.Time) ([]ConsumerCost, error) {
	// Construct the URL
	fullURL := fmt.Sprintf(
		"%s/%s/assets/%s/energy?year=%d&month=%d&day=%d&timespan=daily",
		CustomersURL, customer, meteringPoint, date.Year(), date.Month(), date.Day())

	// Prepare the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+api.login.Token)

	// Send the request using the client
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	hoursBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var costs []ConsumerCost
	if err := json.Unmarshal(hoursBody, &costs); err != nil {
		return nil, err
	}

	return costs, nil
}

func (api *CarunaAPIClient) MeteringPoints(customer string) ([]MeteringPoint, error) {
	// Base URL
	fullURL := CustomersURL + "/" + customer + "/assets"

	// Prepare request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+api.login.Token)

	// Make request using the client
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	meteringPointsBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//log.Println(string(meteringPointsBody))

	var meteringPoints []MeteringPoint
	if err := json.Unmarshal(meteringPointsBody, &meteringPoints); err != nil {
		return nil, err
	}

	return meteringPoints, nil
}

func (api *CarunaAPIClient) Logout() error {
	_, err := api.client.Get(LogoutURL)
	if err != nil {
		return err
	}
	api.login = LoginInfo{}
	return nil
}

func (api *CarunaAPIClient) Login(username, password string) error {
	loginReq, _ := json.Marshal(map[string]string{
		"redirectAfterLogin": RedirectURL,
		"language":           "fi",
	})
	resp, err := api.client.Post(AuthorizationURL, "application/json", strings.NewReader(string(loginReq)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var loginResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return err
	}
	loginRedirectURL := loginResp["loginRedirectUrl"]
	resp, err = api.client.Get(loginRedirectURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	postURL, _ := doc.Find("meta").Attr("content")
	postURL = postURL[6:]
	resp, err = api.client.Get(AuthenticationURL + postURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	action, _ := doc.Find("form").Attr("action")
	action = action[1:][:17] + "IBehaviorListener.0-userIDPanel-usernameLogin-loginWithUserID"
	formData := url.Values{}
	doc.Find("input[type='hidden']").Each(func(index int, item *goquery.Selection) {
		name, _ := item.Attr("name")
		value, exists := item.Attr("value")
		if !exists {
			value = ""
		}
		formData.Set(name, value)
	})
	formData.Set("ttqusername", username)
	formData.Set("userPassword", password)
	submitName, _ := doc.Find("input[type='submit']").Attr("name")
	formData.Set(submitName, "1")
	req, err := http.NewRequest("POST", PortalURL+action, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Wicket-Ajax", "true")
	req.Header.Set("Wicket-Ajax-BaseURL", ".")
	req.Header.Set("Wicket-FocusedElementId", "loginWithUserID5")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", AuthenticationURL)
	req.Header.Set("Referer", PortalURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r := regexp.MustCompile(`CDATA\[.(.*?)\]\]`)
	matches := r.FindAllStringSubmatch(string(bodyBytes), -1)
	resp, err = api.client.Get(PortalURL + matches[0][1])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	metaContent, _ := doc.Find("meta").Attr("content")
	newUrl := metaContent[6:]
	resp, err = api.client.Get(newUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	action, _ = doc.Find("form").Attr("action")
	formData = url.Values{}
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		formData.Add(name, value)
	})
	resp, err = api.client.PostForm(action, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	parsedURL, _ := url.Parse(resp.Request.URL.String())
	resp, err = api.client.PostForm(TokenURL, parsedURL.Query())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	loginInfoBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var loginInfo LoginInfo
	if err := json.Unmarshal(loginInfoBody, &loginInfo); err != nil {
		return err
	}

	api.login = loginInfo

	return nil
}

// eof
