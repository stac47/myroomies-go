package libmyroomies

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/stac47/myroomies/pkg/models"
)

type Client struct {
	url        string
	httpClient *http.Client
	debug      bool
	login      string
	password   string
}

// Create a myroomies client
func NewClient(myroomiesUrl, login, password string, debug bool) (Client, error) {
	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	c := Client{
		url: myroomiesUrl,
		httpClient: &http.Client{
			Transport: tr,
		},
		debug:    debug,
		login:    login,
		password: password,
	}

	return c, nil
}

func (c *Client) ExpenseList() ([]models.Expense, error) {
	expenses := make([]models.Expense, 0)
	err := c.commonHttpGet("/expenses", &expenses)
	return expenses, err
}

func (c *Client) ExpenseCreate(expense models.Expense) (id string, err error) {
	target := fmt.Sprintf("%s%s", c.url, "/expenses")
	expenseStr, err := json.Marshal(expense)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(expenseStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.login, c.password)
	c.debugRequest(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	c.debugResponse(resp)
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		err = fmt.Errorf("Incorrect status code (should be %d): %d", 202, resp.StatusCode)
		return
	}
	location := resp.Header.Get("Location")
	if location == "" {
		err = errors.New("Created resource id cannot be found")
		return
	} else {
		splitted := strings.Split(location, "/")
		if len(splitted) < 1 {
			msg := fmt.Sprintf("Invalid location header: %s", location)
			err = errors.New(msg)
			return
		}
		id = splitted[len(splitted)-1]
	}

	return
}

func (c *Client) ExpenseUpdate(id string, expense *models.Expense) error {
	return errors.New("Not yet implemented")
}

func (c *Client) ExpenseDelete(id string) error {
	target := fmt.Sprintf("%s/%s/%s", c.url, "expenses", id)
	req, err := http.NewRequest("DELETE", target, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.login, c.password)
	c.debugRequest(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	c.debugResponse(resp)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Incorrect status code (should be %d): %d",
			http.StatusOK, resp.StatusCode)
	}
	fmt.Printf("Deleted expense: %s\n", id)
	return nil
}

func (c *Client) commonHttpGet(uri string, obj interface{}) error {
	target := fmt.Sprintf("%s%s", c.url, uri)
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.login, c.password)
	c.debugRequest(req)
	resp, err := c.httpClient.Do(req)
	c.debugResponse(resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Bad HTTP response")
	}
	if err = json.NewDecoder(resp.Body).Decode(obj); err != nil {
		return err
	}
	return nil
}

func (c *Client) debugResponse(resp *http.Response) {
	if c.debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Warnf("Debug response failed: %s", err)
		}
		fmt.Printf("Response:\n%s\n", dump)
	}
}

func (c *Client) debugRequest(req *http.Request) {
	if c.debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Warnf("Debug request failed: %s", err)
		}
		fmt.Printf("Request:\n%s\n", dump)
	}
}
