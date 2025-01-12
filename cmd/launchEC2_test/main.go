package main

import (
	"aws/pkg/launchEC2"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	securityGroup_name = "mySecurityGroup"
	ec2key_name        = "myEC2key" // default extension: ".pem"
	instance_name      = "myEC2instance"
	instance_type      = "t2.micro"              // smallest kind of EC2 instance (available in the free tier)
	AMIid              = "ami-0fda19674ff597992" // amazon linux AMI
)

func main() {
	// loads AWS user configuration from the files ~/.aws/config (to retrieve the AWS region)
	// and ~/.aws/credentials (to retrieve the user AWS access key)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// creates a service client to perform actions on EC2
	ec2client := ec2.NewFromConfig(cfg)

	// creates a security group to define authorized traffic rules to the instance
	err = launchEC2.ConfigureSecurityGroup(ec2client, securityGroup_name)
	if err != nil {
		log.Fatal(err)
	}

	// if the EC2 access key doesn't exist, creates and downloads one.
	// the access key will be used to connect to the instance with SSH
	err = launchEC2.ConfigureAccessKey(ec2client, ec2key_name)
	if err != nil {
		log.Fatal(err)
	}

	// launches the instance with the given parameters.
	// retrieves the instance ID of the instance created.
	instanceID, err := launchEC2.LaunchInstance(ec2client, instance_type, AMIid, securityGroup_name, ec2key_name, instance_name)
	if err != nil {
		log.Fatal(err)
	}

	// fetches the public IP of the newly created instance.
	// public IP takes a couple seconds/minutes to get assigned,
	// and after retrieving it, it can also take a couple seconds/minutes
	// to be accessible.
	_, err = launchEC2.GetPublicIP(ec2client, instanceID)
	if err != nil {
		log.Fatal(err)
	}
}
