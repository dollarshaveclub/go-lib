package awsservice

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWS SDK pointer-itis
var True = true
var False = false

type InstancesDefinition struct {
	AMI           string
	Subnet        string
	SecurityGroup string
	Keypair       string
	Type          string
	UserData      []byte
	Count         int
	RootSizeGB    int
}

func (aws *RealAWSService) RunInstances(idef *InstancesDefinition) ([]string, error) {
	count := int64(idef.Count)
	rs := int64(20)
	vt := "gp2"
	rdn := "/dev/xvda"
	ud := base64.StdEncoding.EncodeToString(idef.UserData)
	if idef.RootSizeGB != 0 {
		rs = int64(idef.RootSizeGB)
	}
	bdm := ec2.BlockDeviceMapping{
		DeviceName: &rdn,
		Ebs: &ec2.EbsBlockDevice{
			DeleteOnTermination: &True,
			Encrypted:           &False,
			VolumeSize:          &rs,
			VolumeType:          &vt,
		},
	}
	ri := ec2.RunInstancesInput{
		ImageId:             &idef.AMI,
		MinCount:            &count,
		MaxCount:            &count,
		KeyName:             &idef.Keypair,
		InstanceType:        &idef.Type,
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{&bdm},
		SecurityGroupIds:    []*string{&idef.SecurityGroup},
		SubnetId:            &idef.Subnet,
		UserData:            &ud,
	}
	r, err := aws.ec2.RunInstances(&ri)
	if err != nil {
		return []string{}, err
	}
	instances := []string{}
	for _, inst := range r.Instances {
		instances = append(instances, *(inst.InstanceId))
	}
	return instances, nil
}

func (aws *RealAWSService) StartInstances(ids []string) error {
	si := ec2.StartInstancesInput{
		InstanceIds: stringSlicetoStringPointerSlice(ids),
	}
	_, err := aws.ec2.StartInstances(&si)
	return err
}

func (aws *RealAWSService) StopInstances(ids []string) error {
	si := ec2.StopInstancesInput{
		InstanceIds: stringSlicetoStringPointerSlice(ids),
	}
	_, err := aws.ec2.StopInstances(&si)
	return err
}

func (aws *RealAWSService) FindInstancesByTag(n string, v string) ([]string, error) {
	fn := fmt.Sprintf("tag:%v", n)
	f := ec2.Filter{
		Name:   &fn,
		Values: []*string{&v},
	}
	dii := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{&f},
	}
	instances := []string{}
	r, err := aws.ec2.DescribeInstances(&dii)
	if err != nil {
		return instances, err
	}
	for _, rev := range r.Reservations {
		for _, inst := range rev.Instances {
			instances = append(instances, *(inst.InstanceId))
		}
	}
	return instances, nil
}

func (aws *RealAWSService) TagInstances(ids []string, n string, v string) error {
	tag := ec2.Tag{
		Key:   &n,
		Value: &v,
	}
	cti := ec2.CreateTagsInput{
		Tags:      []*ec2.Tag{&tag},
		Resources: stringSlicetoStringPointerSlice(ids),
	}
	_, err := aws.ec2.CreateTags(&cti)
	return err
}

func (aws *RealAWSService) DeleteTag(ids []string, n string) error {
	tag := ec2.Tag{
		Key: &n,
	}
	dti := ec2.DeleteTagsInput{
		Tags:      []*ec2.Tag{&tag},
		Resources: stringSlicetoStringPointerSlice(ids),
	}
	_, err := aws.ec2.DeleteTags(&dti)
	return err
}
