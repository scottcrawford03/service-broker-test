package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/oauth2"
)

const (
	APIToken string = "your_token_here"
)

type TokenSource struct {
	AccessToken string
}

type ClusterConnectionData struct {
	Uri      string `json:"uri"`
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSL      string `json:"ssl"`
}

type ClusterData struct {
	Id             string                `json:"id"`
	ConnectionData ClusterConnectionData `json:"connection"`
}
type ClusterResponse struct {
	Cluster ClusterData `json:"cluster"`
}

type TestServiceBroker struct {
	name        string
	instanceMap map[string]*TestServiceInstance
}

type TestServiceInstance struct {
	internalID string
}

func New() (*TestServiceBroker, error) {
	return &TestServiceBroker{
		name:        "this is a test broker",
		instanceMap: make(map[string]*TestServiceInstance),
	}, nil
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

var testplans = []brokerapi.ServicePlan{
	brokerapi.ServicePlan{
		ID:          "1",
		Name:        "Small",
		Description: "Smallest single node cluster running in nyc3",
	},
	brokerapi.ServicePlan{
		ID:          "2",
		Name:        "Medium",
		Description: "Medium single node cluster running in nyc3",
	},
	brokerapi.ServicePlan{
		ID:          "3",
		Name:        "Large",
		Description: "Large single node cluster running in nyc3",
	},
}

type PlanParams struct {
	Name   string
	Region string
	Size   string
}

var plansToParams = map[string]PlanParams{
	"1": PlanParams{
		Name:   "first",
		Region: "nyc3",
		Size:   "1-1-8",
	},
	"2": PlanParams{
		Name:   "second",
		Region: "nyc3",
		Size:   "1-2-30",
	},
	"3": PlanParams{
		Name:   "third",
		Region: "nyc3",
		Size:   "2-4-80",
	},
}

// Services does something with services
func (b *TestServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:                   "1",
			Name:                 "DBaaS",
			Description:          "Database as a service",
			Bindable:             true,
			InstancesRetrievable: true,
			BindingsRetrievable:  true,
			Tags: []string{
				"dbaas",
				"sql",
			},
			PlanUpdatable: false,
			Plans:         testplans,
		},
	}, nil
}

type ProvisionParams struct {
	Name string `json:"name"`
}

func (b *TestServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	fmt.Println("in Provisoin")
	fmt.Println("instanceID", instanceID)

	tokenSource := &TokenSource{
		AccessToken: APIToken,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	url := "https://api.digitalocean.com/v2/databases"
	fmt.Println("URL:>", url)

	fmt.Println("details:", details.PlanID)

	var pp ProvisionParams
	json.Unmarshal(details.GetRawParameters(), &pp)

	var jsonStr PlanParams
	for i, v := range plansToParams {
		if i == details.PlanID {
			jsonStr = v
		}
	}

	fmt.Println("jsonStr:", jsonStr)
	body, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	resp, err := oauthClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return brokerapi.ProvisionedServiceSpec{}, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	var clusterResponse ClusterResponse
	json.Unmarshal(respBody, &clusterResponse)

	var instance TestServiceInstance
	instance.internalID = clusterResponse.Cluster.Id
	b.instanceMap[instanceID] = &instance

	return brokerapi.ProvisionedServiceSpec{
		DashboardURL: "yolo",
		IsAsync:      false,
	}, nil
}

func (b *TestServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fmt.Println("in deprovisoin")
	fmt.Println("instance map", b.instanceMap)

	tokenSource := &TokenSource{
		AccessToken: APIToken,
	}

	instance := b.instanceMap[instanceID]
	url := "https://api.digitalocean.com/v2/databases/" + instance.internalID
	fmt.Println("URL:>", url)

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Fetch Request
	resp, err := oauthClient.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response Body:", resp.Body)

	return brokerapi.DeprovisionServiceSpec{
		OperationData: "Deprovisioned it",
		IsAsync:       false,
	}, nil
}

func (b *TestServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	fmt.Println("in the GetInstance")
	return brokerapi.GetInstanceDetailsSpec{}, nil
}

func (b *TestServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	fmt.Println("in the bind")

	tokenSource := &TokenSource{
		AccessToken: APIToken,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	instance := b.instanceMap[instanceID]
	url := "https://api.digitalocean.com/v2/databases/" + instance.internalID
	fmt.Println("URL:>", url)

	resp, err := oauthClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return brokerapi.Binding{}, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var clusterResponse ClusterResponse
	json.Unmarshal(respBody, &clusterResponse)

	return brokerapi.Binding{
		Credentials: clusterResponse.Cluster.ConnectionData,
	}, nil
}

func (b *TestServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	fmt.Println("in the unbind")
	return brokerapi.UnbindSpec{}, nil
}

func (b *TestServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	fmt.Println("in the getbinding")
	return brokerapi.GetBindingSpec{}, nil
}

func (b *TestServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	fmt.Println("in the update")
	return brokerapi.UpdateServiceSpec{}, nil
}

func (b *TestServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	fmt.Println("in the lastoperation")
	return brokerapi.LastOperation{}, nil
}

func (b *TestServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	fmt.Println("in the lastbindingoperation")
	return brokerapi.LastOperation{}, nil
}
