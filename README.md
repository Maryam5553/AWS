# AWS

This code demonstrates how to use AWS SDK for Go v2 to perform programatic actions on AWS EC2 (ie, launch an instance, create keys on AWS, from code only).

## Prerequisites:

### Go prerequisites
- Go: the v2 SDK requires a minimum version of Go 1.21 (and this code was written with version 1.22.4).

- After cloning repository please `go mod tidy` to retrieve the required Go dependencies.

### AWS prerequisites

Of course, you will need a AWS account.

To execute the code, you need to create a user with permission to access AWS EC2.

TODO explain how to configure user

## Usage

Here is an example on how to test these programs.

Functions in launchEC2.go demonstrate how to configure and launch EC2 instances. From the directory cmd/, execute `go run launchEC2_test/main.go` to launch an instance with default values (you can change the default valurs in the file launchEC2_test/main.go)

Functions in deleteEC2.go demonstrate how to retrieve information on instances, filter by them by tag, and delete them. From the directory cmd/, execute `go run deleteEC2_test/main.go`. You will be able to delete the instances of your choice. To delete instances created with the program launchEC2_test/main.go, simply select "4" (delete instances by giving the name) and enter "myEC2instance" (default name given in the previous program).