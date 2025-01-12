package main

import (
	"aws/pkg/deleteEC2"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Ask user for the name of the instances they want to print,
// then fetch and print their IDs.
func PrintByName(ec2client *ec2.Client) error {
	// prompt user for the name of the instances they want to print
	fmt.Println("You chose to print the IDs of instances that have a certain name.")
	fmt.Print("Please enter the name: ")
	var name string
	_, err := fmt.Scan(&name)
	if err != nil {
		return fmt.Errorf("error reading user input: %w", err)
	}
	// finds and print instance IDs by searching the ones that have the corresponding tag (Name)
	instanceIDs, err := deleteEC2.FindInstanceIDsByTag(ec2client, "Name", name, true)
	if err != nil {
		return err
	}
	if len(instanceIDs) == 0 {
		fmt.Printf("No instance of name \"%s\" was found.\n", name)
	}
	return nil
}

// Prompts user for the IDs of the instances they want to delete
// and deletes them.
func DeleteByInstanceIDs(ec2client *ec2.Client) error {
	// prompt user for the IDs of the instances they want to delete
	fmt.Println("> Please enter the IDs of the instances you want to delete, separated by a space (ex \"i-07aeed4133f5057a6 i-0b7993f98975e0f47\"). Please note that this action is permanent.")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading user input: %w", err)
	}
	instanceIDstring := scanner.Text()
	// stores result in a list
	instanceIDs := strings.Fields(instanceIDstring)

	// delete the instances
	err := deleteEC2.DeleteInstances(ec2client, instanceIDs)
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}

// Prompts user for the name of the instances they want to delete.
// and deletes them. Note: multiple instances can share the same name.
func DeleteByName(ec2client *ec2.Client) error {
	// prompt user for the name of the instances they want to delete
	fmt.Println("This action will delete all instances that have the name that you'll provide. Please note that this action is irreversible.")
	fmt.Print("Please enter the name: ")
	var name string
	_, err := fmt.Scan(&name)
	if err != nil {
		return fmt.Errorf("error reading user input: %w", err)
	}

	// find IDs of corresponding instances
	// (filtering by the tag "Name")
	instanceIDs, err := deleteEC2.FindInstanceIDsByTag(ec2client, "Name", name, false)
	if err != nil {
		return err
	}
	if len(instanceIDs) == 0 {
		fmt.Printf("No instance of name \"%s\" was found.\n", name)
		return nil
	}

	// delete these instances
	err = deleteEC2.DeleteInstances(ec2client, instanceIDs)
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}

func main() {
	// loads AWS user configuration from the files ~/.aws/config (to retrieve the AWS region)
	// and ~/.aws/credentials (to retrieve the user AWS access key)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// creates a service client to perform actions on EC2
	ec2client := ec2.NewFromConfig(cfg)

	// menu
	fmt.Println("This program can delete EC2 instances on your account. Possible actions:")
	fmt.Println("1- Print IDs of all instances on the account (not recommended if you have a lot of instances)")
	fmt.Println("2- Print instance IDs, filtering the instances by name")
	fmt.Println("3- Delete instances by giving the instance IDs")
	fmt.Println("4- Delete instances by giving the name")
	fmt.Println("5- Delete all instances on the account")
	fmt.Println("6- Exit program")

	// get user input
	for {
		fmt.Print("\n> Enter a number: ")

		var answer string
		_, err = fmt.Scan(&answer)
		if err != nil {
			// if getting user input failed, we stay in the loop and try again
			log.Printf("error reading user input: %v", err)
			continue
		}

		switch answer {
		case "1":
			instanceIDs, err := deleteEC2.FindAllInstanceID(ec2client, true)
			if err != nil {
				log.Println(err)
				continue
			}
			if len(instanceIDs) == 0 {
				fmt.Println("No instance found.")
			}
		case "2":
			err = PrintByName(ec2client)
			if err != nil {
				log.Println(err)
				continue
			}
		case "3":
			err = DeleteByInstanceIDs(ec2client)
			if err != nil {
				log.Println(err)
				continue
			}
		case "4":
			err = DeleteByName(ec2client)
			if err != nil {
				log.Println(err)
				continue
			}
		case "5":
			err = deleteEC2.DeleteAllInstances(ec2client)
			if err != nil {
				log.Println(err)
				continue
			}
		case "6":
			return
		default:
			fmt.Print("Invalid output, please answer a number between 1 and 6.")
		}

	}
}
