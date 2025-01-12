package launchEC2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

// Creates a security group allowing traffic to the instance.
func ConfigureSecurityGroup(ec2client *ec2.Client, securityGroupName string) error {
	// information about the security group
	desc := "allow SSH access"
	securityGroupInput := ec2.CreateSecurityGroupInput{
		Description: &desc,
		GroupName:   &securityGroupName,
	}

	// creates the new security group
	_, err := ec2client.CreateSecurityGroup(context.TODO(), &securityGroupInput)

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		// if a security group by this name already exists,
		// I don't do anything (to adapt to the desired usage)
		if apiErr.ErrorCode() == "InvalidGroup.Duplicate" {
			fmt.Printf("Security group %s already exists.\n", securityGroupName)
			return nil
		}
	} else if err != nil {
		return fmt.Errorf("error creating security group %s: %w", securityGroupName, err)
	}

	fmt.Printf("Security Group %s created.\n", securityGroupName)

	/* Add inbound rules to the security group.
	   This will allow the instance to receive traffic from the
	   specified IP range on the specified port.

	   Let's add 2 inbound rules for this example: */

	/* first rule: allow SSH traffic from anywhere. This is necessary if we
	   want to connect to the instance with SSH. We could also specify a more
	   restrictive IP range to be more secure. */
	sourceAll := "0.0.0.0/0"
	tcp := "tcp"
	ssh := int32(22)

	inboundRuleSSH := types.IpPermission{
		FromPort:   &ssh,
		IpProtocol: &tcp,
		IpRanges:   []types.IpRange{{CidrIp: &sourceAll}},
		ToPort:     &ssh,
	}

	// second rule: let's open TCP traffic on port 8080
	port := int32(8080)

	inboundRuleTCP := types.IpPermission{
		FromPort:   &port,
		IpProtocol: &tcp,
		IpRanges:   []types.IpRange{{CidrIp: &sourceAll}},
		ToPort:     &port,
	}

	// adds the 2 rules to the security group previously created.
	inboundRulesInput := ec2.AuthorizeSecurityGroupIngressInput{
		GroupName:     &securityGroupName,
		IpPermissions: []types.IpPermission{inboundRuleSSH, inboundRuleTCP},
	}

	_, err = ec2client.AuthorizeSecurityGroupIngress(context.TODO(), &inboundRulesInput)
	if err != nil {
		return fmt.Errorf("error adding inbound rules to security group %s: %w", securityGroupName, err)
	}

	fmt.Printf("Done configuring security group %s.\n", securityGroupName)

	return nil
}

// If the EC2 access key doesn't exist: creates and downloads one.
// This key will be used to connect with SSH to the instance.
func ConfigureAccessKey(ec2client *ec2.Client, ec2KeyName string) error {
	// first let's check if the desired access key already exists
	describeOutput, err := ec2client.DescribeKeyPairs(context.TODO(), &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return fmt.Errorf("error fetching key pairs info: %w", err)
	}
	exists := false
	for _, keyPair := range describeOutput.KeyPairs {
		if *keyPair.KeyName == ec2KeyName {
			exists = true
			break
		}
	}
	if exists {
		// if it exists, exit.
		fmt.Printf("EC2 key %s already exists.\n", ec2KeyName)
		return nil
	}

	// if the key pair doesn't exist,
	// prompt user to ask if they want to create the key pair.
	fmt.Printf("> EC2 key \"%s\" doesn't exist. Do you want to create it? (Y/N): ", ec2KeyName)
	validAnswer := false
	// prompt loop (until user gives an valid answer)
	for !validAnswer {
		var answer string
		_, err = fmt.Scan(&answer)
		if err != nil {
			return fmt.Errorf("error reading user input: %w", err)
		}
		answer = strings.ToUpper(answer)

		switch answer {
		case "Y":
			validAnswer = true
		case "N":
			/* if user chooses to not create the key,
			   returns an error as the key was unabled to be
			   configure (to adapt to desired usage) */
			return fmt.Errorf("key %s not created", ec2KeyName)
		default:
			fmt.Println("Invalid output, please answer \"Y\" or \"N\": ")
		}
	}

	/* if EC2 key pair doesn't exist and user decides to create it,
	   we create the key. AWS will store the public key and we will
	   retrieve the private key (note: if we don't write the private
	   key in a file now) */
	createKeyPairInput := ec2.CreateKeyPairInput{
		KeyName: &ec2KeyName, // default: format=pem, type=rsa
	}
	key, err := ec2client.CreateKeyPair(context.TODO(), &createKeyPairInput)
	if err != nil {
		return fmt.Errorf("error: creating key pair %s failed: %w", ec2KeyName, err)
	}
	fmt.Printf("Key pair \"%s\" successfully created on AWS.\n", ec2KeyName)

	// write the private key in a file and restrict permissions
	err = os.WriteFile(ec2KeyName+".pem", []byte(*key.KeyMaterial), 0400)
	if err != nil {
		return fmt.Errorf("couldn't create file \"%s.pem\"! The key was created but not downloaded. aws error: %w", ec2KeyName, err)
	}

	fmt.Printf("Private key downloaded in file %s.pem.\n", ec2KeyName)
	return nil
}

/*
Launches a EC2 Instance of type and AMI (Amazon Machine Image) given in parameters.
It will also be associated with the access key, security group, and name given in parameters.
If successfull, this function returns the ID of the instance created. This will be used
to describe the instance later
*/
func LaunchInstance(ec2client *ec2.Client, instanceType string, AMI_id string, securityGroupName string, ec2KeyName string, instanceName string) (string, error) {
	/* Creates a tag to name the instance.
	   The name is simply a tag called "Name".
	   This code would work the same to create any tag.
	   Note: multiple instances can share the same name. */
	tagKey := "Name"
	tag := types.Tag{
		Key:   &tagKey,
		Value: &instanceName,
	}
	tagSpecification := types.TagSpecification{
		ResourceType: "instance",
		Tags:         []types.Tag{tag},
	}

	// launches the instance
	nbInstance := int32(1)
	runInstanceInput := ec2.RunInstancesInput{
		MaxCount:          &nbInstance, // how many instances to launch
		MinCount:          &nbInstance,
		KeyName:           &ec2KeyName,                                // associate access key
		ImageId:           &AMI_id,                                    // give AMI ID
		InstanceType:      types.InstanceType(instanceType),           // give EC2 instance type
		SecurityGroups:    []string{securityGroupName},                // associate security group
		TagSpecifications: []types.TagSpecification{tagSpecification}, // associate Name
	}
	instanceOutput, err := ec2client.RunInstances(context.TODO(), &runInstanceInput)
	if err != nil {
		return "", fmt.Errorf("failed to launch instance: %w", err)
	}
	// print the instance ID. This can be used to fetch additionnal info on the instance
	fmt.Printf("New instance successfully launched, with the following attributes:\n")
	fmt.Printf(" - id: %s\n", *instanceOutput.Instances[0].InstanceId)
	fmt.Printf(" - name: %s\n", instanceName)
	fmt.Printf(" - access key: %s\n", ec2KeyName)
	fmt.Printf(" - type, AMI: %s, %s\n", instanceType, AMI_id)
	fmt.Printf(" - security group: %s\n", securityGroupName)

	return *instanceOutput.Instances[0].InstanceId, nil
}

/*
Retrieves the public IP of the instance, with the instance ID given in parameters.
A EC2 instance takes a couple seconds/minutes to fully start,
so we're polling until the instance has an IP.
It might also take a couple seconds/minutes after that to be able to reach the IP.
The IP can be used to log in to the instance.
*/
func GetPublicIP(ec2client *ec2.Client, instanceId string) (string, error) {
	fmt.Printf("Waiting for the instance %s's public IP...\n", instanceId)

	total_wait := 120  // 2 min wait, arbitrary.
	wait_interval := 1 // 1 second wait between each try.
	tries := 0         // counter to keep track of the number of tries

	// start polling.
	for {
		// abandon after 2 min
		if tries > (total_wait / wait_interval) {
			return "", fmt.Errorf("failed to get instance %s public ip after %d min", instanceId, total_wait/60)
		}
		// increase tries counter
		tries += 1
		time.Sleep(time.Duration(wait_interval) * time.Second)

		// makes a AWS request to fetch info on the instance of ID given in parameter.
		describeInstanceInput := &ec2.DescribeInstancesInput{
			InstanceIds: []string{instanceId},
		}
		describeInstanceOutput, err := ec2client.DescribeInstances(context.TODO(), describeInstanceInput)

		if err != nil {
			return "", fmt.Errorf("failed to fetch info on instance of ID %s: %w", instanceId, err)
		}
		// if AWS sent empty results:
		if len(describeInstanceOutput.Reservations) == 0 {
			return "", fmt.Errorf("couldn't find instance of instance ID %s", instanceId)
		} else if len(describeInstanceOutput.Reservations[0].Instances) == 0 {
			return "", fmt.Errorf("couldn't find instance of instance ID %s", instanceId)
		}

		// checks if field "PublicIpAddress" is empty.
		// if not empty, it means the instance has a public IP and we can finish here.
		if describeInstanceOutput.Reservations[0].Instances[0].PublicIpAddress != nil {
			publicIp := *describeInstanceOutput.Reservations[0].Instances[0].PublicIpAddress
			fmt.Printf("Public IP successfully retrieved: %s\n", publicIp)
			return publicIp, nil
		}
	}
}
