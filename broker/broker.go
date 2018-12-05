package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

type TestServiceBroker struct {
	name string
}

func New() (*TestServiceBroker, error) {
	return &TestServiceBroker{
		name: "this is a test broker",
	}, nil
}

// Services does something with services
func (b *TestServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:                   "1",
			Name:                 "yolo",
			Description:          "this is a service",
			Bindable:             true,
			InstancesRetrievable: true,
			BindingsRetrievable:  true,
			Tags: []string{
				"sure",
				"idk",
			},
			PlanUpdatable: false,
			Plans:         []brokerapi.ServicePlan{},
		},
	}, nil
}

func (b *TestServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}
	fmt.Println("in the provision")
	return spec, nil
}

func (b *TestServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fmt.Println("in the deprovision")
	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *TestServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	fmt.Println("in the GetInstance")
	return brokerapi.GetInstanceDetailsSpec{}, nil
}

func (b *TestServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	fmt.Println("in the bind")
	return brokerapi.Binding{}, nil
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
