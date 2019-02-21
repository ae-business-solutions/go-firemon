package firemon

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CentralSyslogServer struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	DomainID            int    `json:"domainId"`
	IP                  string `json:"ip"`
	CentralSyslogConfig struct {
		ID          int    `json:"id"`
		DomainID    int    `json:"domainId"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Patterns    []struct {
			Description string `json:"description"`
			Pattern     string `json:"pattern"`
		} `json:"patterns"`
		CreatedDate      string `json:"createdDate"`
		LastModifiedDate string `json:"lastModifiedDate"`
		CreatedBy        string `json:"createdBy"`
		LastModifiedBy   string `json:"lastModifiedBy"`
	} `json:"centralSyslogConfig"`
}

type CentralSyslogServers struct {
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
	Count    int                   `json:"count"`
	Results  []CentralSyslogServer `json:"results"`
}

func (c *Client) GetCentralSyslogServers() ([]CentralSyslogServer, error) {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/central-syslog?pageSize=100", c.BaseURL, c.Domain)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	// Extract JSON
	var data CentralSyslogServers
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}
