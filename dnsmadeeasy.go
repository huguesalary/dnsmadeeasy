// package dnsmadeeasy provides a client for the DNS Made Easy API
// The API Reference is currently available at: http://www.dnsmadeeasy.com/wp-content/uploads/2014/07/API-Docv2.pdf
package dnsmadeeasy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	apiUrlFmt        = "https://api.dnsmadeeasy.com/V%s/dns/managed/"
	sandboxApiUrlFmt = "https://api.sandbox.dnsmadeeasy.com/V%s/dns/managed/"
)

// A Client is the basic type of this package. It provides methods for interacting with the DNS Made Easy API.
type Client struct {
	APIKey    string
	APISecret string
	Version   string
	Url       string
}

// DomainList represents a list of Domains
type DomainList struct {
	Domains      []Domain `json:"data"`
	Page         int
	TotalPages   int
	TotalRecords int
}

// Domain represents a DNS Made Easy domain
type Domain struct {
	ActiveThirdParties  []string
	Created             int64
	DelegateNameServers []string
	FolderId            int
	GtdEnabled          bool
	Id                  uint32
	Name                string
	NameServers         []struct {
		Fqdn string
		Ipv4 net.IP
		Ipv6 net.IP
	}
	PendingActionId int
	ProcessMulti    bool
	Updated         int64
}

// DomainRecords holds a list of Records associated to a specific Domain
type DomainRecords struct {
	Records      []*Record `json:"data"`
	Domain       Domain
	Page         int
	TotalPages   int
	TotalRecords int
}

// A Record represents a DNS record (A, CNAME, AAAA, etc.)
type Record struct {
	DynamicDns  bool     `json:",omitempty"`
	Failed      bool     `json:",omitempty"`
	Failover    bool     `json:",omitempty"`
	GtdLocation string   `json:"gtdLocation"`
	HardLink    bool     `json:",omitempty"`
	Id          uint32   `json:"id,omitempty"`
	Monitor     bool     `json:",omitempty"`
	Name        string   `json:"name"`
	Source      int      `json:",omitempty"`
	SourceId    int      `json:",omitempty"`
	Ttl         int      `json:"ttl"`
	Type        string   `json:"type"`
	Value       string   `json:"value"`
	Error       []string `json:"error,omitempty"`
}

// Represents an API Error. Code corresponds to the HTTP Status Code returned by the API.
// Messages is a list of error messages returned by the DNS Made Easy API
type APIError struct {
	Code     int      `json:"-"`
	Messages []string `json:"error"`
}

func (a *APIError) Error() string {
	return fmt.Sprintf("API Error. Code:%d Message:%s", a.Code, strings.Join(a.Messages, " "))
}

// NewClient returns an instance of Client ready to be used for communication with DNS Made Easy's API
func NewClient(key, secret string) *Client {
	return &Client{key, secret, "2.0", fmt.Sprintf(apiUrlFmt, "2.0")}
}

// CreateDomains creates all the domains passed in its argument. When multiple domains are created at once, the API only returns domain ids.
func (c *Client) CreateDomains(domainNames []string) ([]uint32, error) {

	var data map[string][]string
	data = make(map[string][]string)
	data["names"] = domainNames

	body, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(body)
	r, err := http.NewRequest("POST", c.Url, buf)

	if err != nil {
		return nil, err
	}

	var ids []uint32

	return ids, c.request(r, &ids)
}

func (c *Client) CreateDomain(domainName string) (*Domain, error) {

	var data map[string][]string
	data = make(map[string][]string)
	data["names"] = append(data["name"], domainName)

	body, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(body)

	r, err := http.NewRequest("POST", c.Url, buf)

	if err != nil {
		return nil, err
	}

	domain := &Domain{}

	return domain, c.request(r, domain)
}

func (c *Client) DeleteDomain(id uint32) error {
	return c.DeleteDomains([]uint32{id})
}

func (c *Client) DeleteDomains(ids []uint32) error {

	body, err := json.Marshal(ids)

	if err != nil {
		return err
	}

	buf := bytes.NewReader(body)
	r, err := http.NewRequest("DELETE", c.Url, buf)
	if err != nil {
		return err
	}

	return c.request(r, nil)
}

func (c *Client) UpdDomains(ids []uint32, fields map[string]interface{}) error {

	var data map[string]interface{}
	data = make(map[string]interface{})

	data["ids"] = ids
	for k, v := range fields {
		data[k] = v
	}

	body, err := json.Marshal(data)
	buf := bytes.NewReader(body)
	r, err := http.NewRequest("PUT", c.Url, buf)

	if err != nil {
		return err
	}

	return c.request(r, nil)
}

func (c *Client) GetDomains() (*DomainList, error) {

	r, err := http.NewRequest("GET", c.Url, nil)

	if err != nil {
		return nil, err
	}

	domainList := &DomainList{}
	err = c.request(r, domainList)

	return domainList, err
}

func (c *Client) GetDomainById(id uint32) (*Domain, error) {

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%d", c.Url, id), nil)

	if err != nil {
		return nil, err
	}

	domain := &Domain{}
	err = c.request(r, domain)

	return domain, err
}

func (c *Client) GetDomainByName(name string) (*Domain, error) {

	r, err := http.NewRequest("GET", fmt.Sprintf("%sname?domainname=%s", c.Url, name), nil)

	if err != nil {
		return nil, err
	}

	domain := &Domain{}
	err = c.request(r, domain)

	return domain, err
}

func (c *Client) GetDomainRecords(id uint32) ([]*Record, error) {

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%d/records", c.Url, id), nil)

	if err != nil {
		return nil, err
	}

	domainRecords := &DomainRecords{}
	err = c.request(r, domainRecords)

	return domainRecords.Records, err
}

func (c *Client) AddRecord(domainId uint32, record *Record) error {

	body, err := json.Marshal(record)

	if err != nil {
		return err
	}

	buf := bytes.NewReader(body)
	r, err := http.NewRequest("POST", fmt.Sprintf("%s%d/records", c.Url, domainId), buf)
	if err != nil {
		return nil
	}

	return c.request(r, record)
}

func (c *Client) DelRecord(domainId, recordId uint32) error {
	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s%d/records/%d", c.Url, domainId, recordId), nil)
	if err != nil {
		return err
	}

	return c.request(r, nil)
}

func (c *Client) UpdRecord(domainId uint32, record *Record) error {

	body, err := json.Marshal(record)

	if err != nil {
		return err
	}

	buf := bytes.NewReader(body)

	r, err := http.NewRequest("PUT", fmt.Sprintf("%s%d/records/%d", c.Url, domainId, record.Id), buf)
	if err != nil {
		return err
	}

	return c.request(r, record)
}

// request takes care of adding the authentication headers, making the request and unmarshaling the JSON response
func (c *Client) request(r *http.Request, object interface{}) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	r.Header.Add("x-dnsme-apikey", c.APIKey)
	requestDate := time.Now().UTC().Format(time.RFC1123)
	r.Header.Add("x-dnsme-requestdate", requestDate)

	h := hmac.New(sha1.New, []byte(c.APISecret))
	h.Write([]byte(requestDate))
	r.Header.Add("x-dnsme-hmac", fmt.Sprintf("%x", h.Sum(nil)))

	r.Header.Add("Accept", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		apierr := &APIError{}
		json.Unmarshal(buf.Bytes(), apierr)
		apierr.Code = resp.StatusCode

		return apierr
	}

	if buf.Len() > 0 {
		return json.Unmarshal(buf.Bytes(), object)
	}

	return nil
}
