package firemon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DeviceGroup struct {
	ID                   *int     `json:"id,omitempty"`
	DomainID             int      `json:"domainId"`
	Name                 string   `json:"name"`
	Description          string   `json:"description,omitempty"`
	ParentID             *int     `json:"parentId,omitempty"`
	SecurityConcernIndex *float32 `json:"securityConcernIndex,omitempty"`
	Analysis             bool     `json:"analysis"`
	ChildDeviceGroups    *int     `json:"childDeviceGroups,omitempty"`
	ChildDevices         *int     `json:"childDevices,omitempty"`
}

type DeviceGroups struct {
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Count    int           `json:"count"`
	Results  []DeviceGroup `json:"results"`
}

func (c *Client) GetDeviceGroups() ([]DeviceGroup, error) {
	var devicegroups []DeviceGroup
	page := 0
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/devicegroup?pageSize=100&page=%d", c.BaseURL, c.Domain, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("Content-Type", "application/json")
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
			return nil, fmt.Errorf("Error: Status: %d, Body: %s", resp.StatusCode, body)
		}
		var data DeviceGroups
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		if data.Count == 0 {
			break
		}
		devicegroups = append(devicegroups, data.Results...)
		page++
	}

	return devicegroups, nil
}

func (c *Client) GetDeviceGroupDevices(devicegroupid int) ([]Device, error) {
	var devices []Device
	page := 0
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/devicegroup/%d/device?pageSize=100", c.BaseURL, c.Domain, devicegroupid)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("Content-Type", "application/json")
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
			return nil, fmt.Errorf("Error: Status: %d, Body: %s", resp.StatusCode, body)
		}
		var data Devices
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		if data.Count == 0 {
			break
		}
		devices = append(devices, data.Results...)
		page++
	}

	return devices, nil
}

func (c *Client) AddDeviceToDeviceGroup(devicegroupid, deviceid int) error {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/devicegroup/%d/device/%d", c.BaseURL, c.Domain, devicegroupid, deviceid)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if 204 != resp.StatusCode { // API returned 204 on success; still successful even if device was already in the device group
		return fmt.Errorf("Error adding device (ID: %d) to device group (ID: %d) [Status: %d]", deviceid, devicegroupid, resp.StatusCode)
	}
	return nil
}

func (c *Client) DeleteDeviceFromDeviceGroup(devicegroupid, deviceid int) error {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/devicegroup/%d/device/%d", c.BaseURL, c.Domain, devicegroupid, deviceid)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if 204 != resp.StatusCode { // API returned 204 on success; still successful even if device is not present in the device group
		return fmt.Errorf("Error removing device (ID: %d) from device group (ID: %d) [Status: %d]", deviceid, devicegroupid, resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateDeviceGroup(name, description string) error {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/devicegroup", c.BaseURL, c.Domain)
	devicegroup := DeviceGroup{
		DomainID:    c.Domain,
		Name:        name,
		Description: description,
		Analysis:    true,
	}
	body, err := JSONMarshal(devicegroup)
	if err != nil {
		return err
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	// API returns 400 when Device Group already exists
	if resp.StatusCode != 200 || resp.StatusCode != 400 {
		return fmt.Errorf("Error (Status: %d), Body: %s", resp.StatusCode, body)
	}
	return nil
}
