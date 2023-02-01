/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"github.com/rkosegi/active24-go/active24"
	"k8s.io/klog/v2"
)

type Config struct {
	ApiKey     string
	ApiUrl     string
	DomainName string
}

type ApiClient struct {
	dns active24.DnsRecordActions
	dom string
}

// FindTxtRecord Find TXT record by name and content
func (a *ApiClient) FindTxtRecord(name string, text string) (*active24.DnsRecord, error) {
	klog.V(4).Infof("FindTxtRecord: name=%s, text=%s", name, text)

	records, err := a.dns.List()
	if err != nil {
		klog.V(1).ErrorS(err.Error(), "invalid API response", "code", err.Response().Status)
		return nil, err.Error()
	}
	for _, record := range records {
		klog.V(9).Infof("record=%v", record)
		if record.Name == name && *record.Type == "TXT" && *record.Text == text {
			return &record, nil
		}
	}
	return nil, nil
}

// UpdateTxtRecord Update existing DNS TXT record
func (a *ApiClient) UpdateTxtRecord(hashId string, name string, text string, ttl int) error {
	klog.V(4).Infof("UpdateTxtRecord: domain=%s, name=%s, text=%s, ttl=%d, hashId=%s",
		a.dom, name, text, ttl, hashId)
	err := a.dns.Update(active24.DnsRecordTypeTXT, &active24.DnsRecord{
		HashId: &hashId,
		Name:   name,
		Text:   &text,
		Ttl:    ttl,
	})
	if err != nil {
		klog.V(1).ErrorS(err.Error(), "invalid API response", "code", err.Response().Status)
		return err.Error()
	}
	return nil
}

// NewTxtRecord Create new DNS TXT record
func (a *ApiClient) NewTxtRecord(name string, text string, ttl int) error {
	klog.V(4).Infof("NewTxtRecord: domain=%s, name=%s, text=%s, ttl=%d",
		a.dom, name, text, ttl)
	err := a.dns.Create(active24.DnsRecordTypeTXT, &active24.DnsRecord{
		Ttl:  ttl,
		Name: name,
		Text: &text,
	})
	if err != nil {
		klog.V(1).ErrorS(err.Error(), "invalid API response", "code", err.Response().Status)
		return err.Error()
	}
	return nil
}

// DeleteTxtRecord Delete existing DNS record
func (a *ApiClient) DeleteTxtRecord(hashId string) error {
	klog.V(4).Infof("DeleteTxtRecord: domain=%s, hashId=%s", a.dom, hashId)
	err := a.dns.Delete(hashId)
	if err != nil {
		klog.V(1).ErrorS(err.Error(), "invalid API response", "code", err.Response().Status)
		return err.Error()
	}
	return nil
}

func NewApiClient(config Config) *ApiClient {
	opts := make([]active24.Option, 0)
	if len(config.ApiUrl) > 0 {
		opts = append(opts, active24.ApiEndpoint(config.ApiUrl))
	}
	return &ApiClient{
		dns: active24.New(config.ApiKey, opts...).Dns().With(config.DomainName),
		dom: config.DomainName,
	}
}
