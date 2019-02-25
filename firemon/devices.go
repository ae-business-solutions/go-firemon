package firemon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Device struct {
	ID           int    `json:"id"`
	DomainID     int    `json:"domainId"`
	Name         string `json:"name"`
	ManagementIP string `json:"managementIp"`
	Parents      *struct {
		ID           int    `json:"id"`
		DomainID     int    `json:"domainId"`
		Name         string `json:"name"`
		ManagementIP string `json:"managementIp"`
		Vendor       string `json:"vendor"`
		DeviceType   string `json:"deviceType"`
		State        string `json:"state"`
	} `json:"parents"`
	Children             []interface{} `json:"children"`
	DataCollectorID      int           `json:"dataCollectorId"`
	CentralSyslogID      int           `json:"centralSyslogId,omitempty"`
	SyslogMatchName      string        `json:"syslogMatchName,omitempty"`
	SecurityConcernIndex float64       `json:"securityConcernIndex"`
	Licenses             []string      `json:"licenses"`
	DevicePack           *struct {
		Type       string `json:"type"`
		ID         int    `json:"id"`
		ArtifactID string `json:"artifactId"`
		GroupID    string `json:"groupId"`
		Version    string `json:"version"`
		Artifacts  []struct {
			Name     string `json:"name"`
			Checksum string `json:"checksum"`
		} `json:"artifacts"`
		DeviceName       string `json:"deviceName"`
		DeviceType       string `json:"deviceType"`
		Vendor           string `json:"vendor"`
		CollectionConfig struct {
			ID                   int    `json:"id"`
			Name                 string `json:"name"`
			DevicePackID         int    `json:"devicePackId"`
			DevicePackVendor     string `json:"devicePackVendor"`
			DevicePackDeviceType string `json:"devicePackDeviceType"`
			DevicePackDeviceName string `json:"devicePackDeviceName"`
			DevicePackGroupID    string `json:"devicePackGroupId"`
			DevicePackArtifactID string `json:"devicePackArtifactId"`
			ChangePattern        string `json:"changePattern"`
			ChangeCriterion      []struct {
				Pattern         string `json:"pattern"`
				TimeoutSeconds  int    `json:"timeoutSeconds"`
				ContinueMatch   bool   `json:"continueMatch"`
				ParentUserName  string `json:"parentUserName"`
				RetrieveOnMatch bool   `json:"retrieveOnMatch"`
			} `json:"changeCriterion"`
			UsagePattern   string `json:"usagePattern"`
			UsageCriterion []struct {
				Pattern       string        `json:"pattern"`
				Fields        []interface{} `json:"fields"`
				DynamicFields []interface{} `json:"dynamicFields"`
			} `json:"usageCriterion"`
			CreatedDate            string   `json:"createdDate"`
			LastModifiedDate       string   `json:"lastModifiedDate"`
			CreatedBy              string   `json:"createdBy"`
			LastModifiedBy         string   `json:"lastModifiedBy"`
			UsageKeys              []string `json:"usageKeys"`
			ActivatedForDevicePack bool     `json:"activatedForDevicePack"`
		} `json:"collectionConfig"`
		BehaviorTranslator      string        `json:"behaviorTranslator"`
		Normalization           bool          `json:"normalization"`
		Usage                   bool          `json:"usage"`
		Change                  bool          `json:"change"`
		UsageSyslog             bool          `json:"usageSyslog"`
		ChangeSyslog            bool          `json:"changeSyslog"`
		Active                  bool          `json:"active"`
		SupportsDiff            bool          `json:"supportsDiff"`
		SupportsManualRetrieval bool          `json:"supportsManualRetrieval"`
		ImplicitDrop            bool          `json:"implicitDrop"`
		DiffDynamicRoutes       bool          `json:"diffDynamicRoutes"`
		Automation              bool          `json:"automation"`
		LookupNoIntfRoutes      bool          `json:"lookupNoIntfRoutes"`
		AutomationCli           bool          `json:"automationCli"`
		SSH                     bool          `json:"ssh"`
		SharedNetworks          bool          `json:"sharedNetworks"`
		SharedServices          bool          `json:"sharedServices"`
		SupportedTypes          []string      `json:"supportedTypes"`
		DiffIgnorePatterns      []string      `json:"diffIgnorePatterns"`
		ConvertableTo           []interface{} `json:"convertableTo"`
	} `json:"devicePack"`
	GpcDirtyDate         string `json:"gpcDirtyDate"`
	GpcComputeDate       string `json:"gpcComputeDate"`
	GpcImplementDate     string `json:"gpcImplementDate"`
	State                string `json:"state"`
	ExtendedSettingsJSON *struct {
		SSHPort                       int    `json:"sshPort,omitempty"`
		Password                      string `json:"password,omitempty"`
		RestPort                      int    `json:"restPort,omitempty"`
		Username                      string `json:"username,omitempty"`
		Connected                     bool   `json:"connected"`
		SupportsFQDN                  bool   `json:"supportsFQDN"`
		LoggingPlugin                 string `json:"loggingPlugin,omitempty"`
		RetrievalMethod               string `json:"retrievalMethod,omitempty"`
		RetrievalPlugin               string `json:"retrievalPlugin,omitempty"`
		MonitoringPlugin              string `json:"monitoringPlugin,omitempty"`
		ResetSSHKeyValue              bool   `json:"resetSSHKeyValue"`
		LogUpdateInterval             int    `json:"logUpdateInterval,omitempty"`
		BatchConfigRetrieval          bool   `json:"batchConfigRetrieval"`
		LogMonitoringEnabled          bool   `json:"logMonitoringEnabled"`
		LogRecordCacheTimeout         int    `json:"logRecordCacheTimeout,omitempty"`
		SkipUserFileRetrieval         bool   `json:"skipUserFileRetrieval"`
		ChangeMonitoringEnabled       bool   `json:"changeMonitoringEnabled"`
		SuppressFQDNCapabilities      bool   `json:"suppressFQDNCapabilities"`
		ScheduledRetrievalEnabled     bool   `json:"scheduledRetrievalEnabled"`
		ScheduledRetrievalInterval    int    `json:"scheduledRetrievalInterval,omitempty"`
		SkipDynamicBlockListRetrieval bool   `json:"skipDynamicBlockListRetrieval"`
	} `json:"extendedSettingsJson"`
	Cluster *struct {
		ID             int    `json:"id"`
		DomainID       int    `json:"domainId"`
		Name           string `json:"name"`
		ActiveDeviceID int    `json:"activeDeviceId"`
	} `json:"cluster,omitempty"`
	Editable  bool   `json:"editable"`
	GpcStatus string `json:"gpcStatus"`
}

type Devices struct {
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Count    int      `json:"count"`
	Results  []Device `json:"results"`
}

func (c *Client) GetDevices() ([]Device, error) {
	var devices []Device
	page := 0
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/device?pageSize=100&page=%d", c.BaseURL, c.Domain, page)
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
		if data.Count == 0 {
			break
		}
		devices = append(devices, data.Results...)
		page++
	}

	return devices, nil
}

func (c *Client) GetDevicesByName(pattern string) ([]Device, error) {
	var devices []Device
	page := 0
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/device?pageSize=100&page=%d", c.BaseURL, c.Domain, page)
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
		if data.Count == 0 {
			break
		}
		for _, device := range data.Results {
			match, _ := regexp.MatchString(pattern, device.Name)
			if match {
				devices = append(devices, device)
			}
		}
		page++
	}

	return devices, nil
}

func (c *Client) UpdateDevice(device Device) error {
	url := fmt.Sprintf("https://%s/securitymanager/api/domain/%d/device/%d", c.BaseURL, c.Domain, device.ID)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Using custom version of json.Marshal that does NOT escape HTML symbols (<>&)
	body, err := JSONMarshal(device)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
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
	if 204 != resp.StatusCode {
		return fmt.Errorf("%s", body)
	}
	return nil
}
