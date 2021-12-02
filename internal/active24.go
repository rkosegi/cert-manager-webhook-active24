package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/klog"
)

type Config struct {
	ApiKey     string
	ApiUrl     string
	DomainName string
}

type DnsRecord struct {
	HashId     string `json:"hashId,omitempty"`
	NameServer string `json:"nameServer,omitempty"`
	Type       string `json:"type,omitempty"`
	Name       string `json:"name"`
	Text       string `json:"text"`
	Ttl        int    `json:"ttl"`
}

type ApiClient struct {
	config Config
}

//Find TXT record by name
func (a *ApiClient) FindTxtRecord(name string, text string) (*DnsRecord, error) {
	klog.V(4).Infof("Find TXT record, name=%s, text=%s", name, text)
	records, err := a.FetchDnsRecords()
	if err != nil {
		return nil, err
	}
	for _, record := range *records {
		klog.V(9).Infof("record=%v", record)
		if record.Name == name && record.Type == "TXT" && record.Text == text {
			return &record, nil
		}
	}
	return nil, nil
}

func (a *ApiClient) callApi(method string, uri string, data interface{}) (*http.Response, error) {
	var err error
	var req *http.Request
	url := fmt.Sprintf("%s/%s/%s", a.config.ApiUrl, a.config.DomainName, uri)

	klog.V(4).Infof("API request: method=%s, url=%s, data=%v", method, url, data)

	if data != nil {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.config.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	klog.V(4).Infof("API response: status=%d", resp.StatusCode)
	return resp, err
}

// Get all DNS records
func (a *ApiClient) FetchDnsRecords() (*[]DnsRecord, error) {
	klog.V(4).Infof("Client: Fetch DNS records, domain=%s", a.config.DomainName)
	resp, err := a.callApi("GET", "records/v1", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code returned by API: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var records []DnsRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

// Update existing DNS TXT record
func (a *ApiClient) UpdateTxtRecord(hashId string, name string, text string, ttl int) (int, error) {
	klog.V(4).Infof("Update DNS record, domain=%s, name=%s, text=%s, ttl=%d, hashId=%s",
		a.config.DomainName, name, text, ttl, hashId)
	resp, err := a.callApi("PUT", "txt/v1", DnsRecord{
		HashId: hashId,
		Name:   name,
		Text:   text,
		Ttl:    ttl,
	})
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, err
}

// Create new DNS TXT record
func (a *ApiClient) NewTxtRecord(name string, text string, ttl int) (int, error) {
	klog.V(4).Infof("Create DNS record, domain=%s, name=%s, text=%s, ttl=%d",
		a.config.DomainName, name, text, ttl)
	resp, err := a.callApi("POST", "txt/v1", DnsRecord{
		Name: name,
		Text: text,
		Ttl:  ttl,
	})
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func (a *ApiClient) DeleteTxtRecord(hashId string) (int, error) {
	klog.V(4).Infof("Delete DNS record, domain=%s, hashId=%s", a.config.DomainName, hashId)

	resp, err := a.callApi("DELETE", fmt.Sprintf("%s/v1", hashId), nil)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func NewApiClient(config Config) *ApiClient {
	return &ApiClient{
		config: config,
	}
}
