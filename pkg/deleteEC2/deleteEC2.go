package deleteEC2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Permanently deletes the EC2 instance of ID given in parameter.
func DeleteInstance(ec2client *ec2.Client, instanceId string) error {
	// indicate instance ID in parameter of the termination request
	terminateInstanceInput := &ec2.TerminateInstancesInput{InstanceIds: []string{instanceId}}
	// terminates the instance
	_, err := ec2client.TerminateInstances(context.TODO(), terminateInstanceInput)
	if err != nil {
		return fmt.Errorf("couldn't delete instance %s: %w", instanceId, err)
	}
	fmt.Printf("Instance %s successfully deleted.\n", instanceId)
	return err
}

// Permanently deletes the EC2 instances of IDs given in parameter (in a list).
func DeleteInstances(ec2client *ec2.Client, instanceIdList []string) error {
	success := 0
	fails := 0
	for _, instance := range instanceIdList {
		err := DeleteInstance(ec2client, instance)
		// if we couldn't delete one instance,
		// print the error and keep going
		if err != nil {
			log.Printf("%s\n", err)
			fails++
		} else {
			success++
		}
	}
	// at the end, if some instances coudln't be deleted,
	// we return an error informing how many failed
	if fails > 0 {
		return fmt.Errorf("error deleting multiple instances: %d instances were successfully deleted, and %d instances couldn't be deleted", success, fails)
	}
	fmt.Printf("%d instances were successfully deleted.\n", success)
	return nil
}

// Returns a list containing the ID of all instances owned.
// If print=true, instances ID are also printed.
func FindAllInstanceID(ec2client *ec2.Client, print bool) ([]string, error) {
	var instanceIDs []string

	// describe all instances owned
	describeInstanceOutput, err := ec2client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return instanceIDs, fmt.Errorf("fetching info on instances failed: %w", err)
	}

	// add each non-terminated/non-terminating instance found
	// to the list
	for _, reservation := range describeInstanceOutput.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name != "shutting-down" && instance.State.Name != "terminated" {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
				if print {
					fmt.Println(*instance.InstanceId)
				}
			}
		}
	}
	// returns the list
	return instanceIDs, nil
}

// Returns a list containing the ID of all the instances
// that have the tag "key=value", and not all owned instances.
// If print=true, instances ID are also printed.
func FindInstanceIDsByTag(ec2client *ec2.Client, tagKey string, tagValue string, print bool) ([]string, error) {
	var instanceIDs []string

	// creates a filter for the describe instances request
	// to filter the instances that have the tag "key=value"
	tag := "tag:" + tagKey
	filters := []types.Filter{
		{
			Name:   &tag,
			Values: []string{tagValue},
		},
	}

	describeInstanceInput := &ec2.DescribeInstancesInput{Filters: filters}

	// make the request
	describeInstanceOutput, err := ec2client.DescribeInstances(context.TODO(), describeInstanceInput)
	if err != nil {
		return instanceIDs, fmt.Errorf("fetching info on instances with tag %s=%s failed: %w", tagKey, tagValue, err)
	}

	// add each non-terminated/non-terminating instance found
	// to the list
	for _, reservation := range describeInstanceOutput.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name != "shutting-down" && instance.State.Name != "terminated" {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
				if print {
					fmt.Println(*instance.InstanceId)
				}
			}
		}
	}
	// returns the list
	return instanceIDs, nil
}

// Permanently delete all the EC2 instances owned by the user.
// This first asks the user for confirmation.
func DeleteAllInstances(ec2client *ec2.Client) error {
	// asks user for confirmation before deleting all the instances
	fmt.Println("> You asked to delete **ALL** EC2 instances owned on your account. This action is non-reversible. Proceed? (Y/N): ")
	validAnswer := false
	// prompt loop (until user gives an valid answer)
	for !validAnswer {
		var answer string
		_, err := fmt.Scan(&answer)
		if err != nil {
			return fmt.Errorf("error reading user input: %w", err)
		}
		answer = strings.ToUpper(answer)

		switch answer {
		case "Y":
			validAnswer = true
		case "N":
			fmt.Println("Action aborted.")
			return nil
		default:
			fmt.Println("Invalid output, please answer \"Y\" or \"N\": ")
		}
	}

	// Fetches the IDs of all the instances owned by user
	instancesIDs, err := FindAllInstanceID(ec2client, false)
	if err != nil {
		return err
	}

	// Terminates all instances found
	if len(instancesIDs) == 0 {
		fmt.Println("No instance found.")
	} else {
		for _, id := range instancesIDs {
			err = DeleteInstance(ec2client, id)
			if err != nil {
				return err
			}
		}
		fmt.Println("All instances successfully deleted.")
	}

	fmt.Println("Done")
	return nil
}
