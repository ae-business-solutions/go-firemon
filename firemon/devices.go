package firemon

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Device struct {
	ID                   int      `json:"id"`
	DomainID             int      `json:"domainId"`
	Name                 string   `json:"name"`
	ManagementIP         string   `json:"managementIp"`
	DataCollectorID      int      `json:"dataCollectorId"`
	SecurityConcernIndex float64  `json:"securityConcernIndex"`
	Licenses             []string `json:"licenses"`
	State                string   `json:"state"`
}

type Devices struct {
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Count    int      `json:"count"`
	Results  []Device `json:"results"`
}

func (c *Client) GetDevices() ([]Device, error) {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/device?pageSize=100", c.BaseURL, c.Domain)
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
	var data Devices
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Results, nil
}
