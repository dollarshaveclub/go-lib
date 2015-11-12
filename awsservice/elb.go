package awsservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elb"
)

type ELBListener struct {
	InstancePort         int64
	LoadBalancerPort     int64
	LoadBalancerProtocol string
	InstanceProtocol     string
	CertificateID        string
}

type LoadBalancerDefinition struct {
	Listeners      []ELBListener
	Name           string
	SecurityGroups []string
	Scheme         string
	Subnets        []string
}

func (aws *RealAWSService) CreateLoadBalancer(lbd *LoadBalancerDefinition) (string, error) {
	listeners := []*elb.Listener{}
	for _, l := range lbd.Listeners {
		ip := l.InstancePort // allocate new objects so pointers in struct are unique
		lbp := l.LoadBalancerPort
		pr := l.LoadBalancerProtocol
		ipr := l.InstanceProtocol
		cid := l.CertificateID
		listeners = append(listeners, &elb.Listener{
			InstancePort:     &ip,
			LoadBalancerPort: &lbp,
			Protocol:         &pr,
			InstanceProtocol: &ipr,
			SSLCertificateId: &cid,
		})
	}
	o, err := aws.elbc.CreateLoadBalancer(&elb.CreateLoadBalancerInput{
		Listeners:        listeners,
		LoadBalancerName: &lbd.Name,
		SecurityGroups:   stringSlicetoStringPointerSlice(lbd.SecurityGroups),
		Subnets:          stringSlicetoStringPointerSlice(lbd.Subnets),
	})
	if err != nil {
		return "", err
	}
	return *o.DNSName, nil
}

func (aws *RealAWSService) DeleteLoadBalancer(n string) error {
	_, err := aws.elbc.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: &n,
	})
	return err
}

func (aws *RealAWSService) RegisterInstances(n string, ids []string) error {
	_, err := aws.elbc.RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        instanceIDSlice(ids),
		LoadBalancerName: &n,
	})
	return err
}

func (aws *RealAWSService) DeregisterInstances(n string, ids []string) error {
	_, err := aws.elbc.DeregisterInstancesFromLoadBalancer(&elb.DeregisterInstancesFromLoadBalancerInput{
		Instances:        instanceIDSlice(ids),
		LoadBalancerName: &n,
	})
	return err
}

// Testing mocks

func (aws *TestingAWSService) CreateLoadBalancer(lbd *LoadBalancerDefinition) (string, error) {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "CreateLoadBalancer",
		NotableParams: map[string]string{
			"name":            lbd.Name,
			"security_groups": fmt.Sprintf("%v", lbd.SecurityGroups),
			"scheme":          lbd.Scheme,
			"subnets":         fmt.Sprintf("%v", lbd.Subnets),
			"listeners":       fmt.Sprintf("%v", lbd.Listeners),
		},
	})
	return "", nil
}

func (aws *TestingAWSService) DeleteLoadBalancer(n string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "DeleteLoadBalancer",
		NotableParams: map[string]string{
			"name": n,
		},
	})
	return nil
}

func (aws *TestingAWSService) RegisterInstances(n string, ids []string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "RegisterInstances",
		NotableParams: map[string]string{
			"ids": fmt.Sprintf("%v", ids),
		},
	})
	return nil
}

func (aws *TestingAWSService) DeregisterInstances(n string, ids []string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "DeregisterInstances",
		NotableParams: map[string]string{
			"ids": fmt.Sprintf("%v", ids),
		},
	})
	return nil
}
