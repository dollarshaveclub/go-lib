package awsservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/route53"
)

var awsRegion = "us-west-2"

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

type Route53RecordDefinition struct {
	ZoneID string
	Name   string
	Value  string
	Type   string
	TTL    int64
}

type AWSLoadBalancerService interface {
	CreateLoadBalancer(*LoadBalancerDefinition) (string, error)
	DeleteLoadBalancer(string) error
	RegisterInstances(string, []string) error
	DeregisterInstances(string, []string) error
}

type AWSRoute53Service interface {
	CreateDNSRecord(*Route53RecordDefinition) error
	DeleteDNSRecord(*Route53RecordDefinition) error
}

type AWSService interface {
	AWSLoadBalancerService
	AWSRoute53Service
}

type LimitedRoute53API interface {
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
}

type LimitedELBAPI interface {
	CreateLoadBalancer(*elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error)
	DeleteLoadBalancer(*elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error)
	RegisterInstancesWithLoadBalancer(*elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	DeregisterInstancesFromLoadBalancer(*elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error)
}

type RealAWSService struct {
	elbc *elb.ELB
	r53c *route53.Route53
}

// Testing types
type AWSActionLog struct {
	Action        string
	NotableParams map[string]string
}

type TestingAWSService struct {
	Log []AWSActionLog
}

// NewStaticAWSService uses the static credential provider (pass in access key ID and secret key)
func NewStaticAWSService(id string, secret string) AWSService {
	s := session.New(&aws.Config{Credentials: credentials.NewStaticCredentials(id, secret, ""), Region: &awsRegion})

	return &RealAWSService{
		elbc: elb.New(s),
		r53c: route53.New(s),
	}
}

// NewAWSService uses the default Environment credential store
func NewAWSService() AWSService {
	s := session.New(&aws.Config{Region: &awsRegion})

	return &RealAWSService{
		elbc: elb.New(s),
		r53c: route53.New(s),
	}
}

// Stupid AWS SDK...
func stringSlicetoStringPointerSlice(s []string) []*string {
	o := []*string{}
	for _, str := range s {
		nstr := str
		o = append(o, &nstr)
	}
	return o
}

func instanceIDSlice(ids []string) []*elb.Instance {
	instances := []*elb.Instance{}
	for _, id := range ids {
		curid := id
		instances = append(instances, &elb.Instance{
			InstanceId: &curid,
		})
	}
	return instances
}
