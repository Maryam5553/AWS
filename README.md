# AWS

This code demonstrates how to use AWS SDK for Go v2 to perform programmatic actions on AWS EC2 (ie, launch an instance, create EC2 key pairs, delete an instance, from code only).

## Prerequisites

### Go prerequisites
- Go: the v2 SDK requires a minimum version of Go 1.21 (and this code was written with version 1.22.4).

- After cloning repository please `go mod tidy` to retrieve the required Go dependencies.

### AWS prerequisites

Of course, you will need a AWS account.

To execute the code, you need to create a user with permission to access AWS EC2. Here is how to do it:

- In the AWS IAM console, go to [Users](https://us-east-1.console.aws.amazon.com/iam/home?region=us-east-1#/users). Create a user and attach the policy [AmazonEC2FullAccess](https://us-east-1.console.aws.amazon.com/iam/home?region=us-east-1#/policies/details/arn%3Aaws%3Aiam%3A%3Aaws%3Apolicy%2FAmazonEC2FullAccess?section=permissions).

- Create an access key: click on the user, select "Create Access Key". At the end, download or write down the Access Key ID and Secret Access Key.

- Create folder `~/.aws/` and create two files with the following contents:

    - `~/.aws/credentials`:

    ```
    [default]
    aws_access_key_id = YOUR_ACCESS_KEY_ID
    aws_secret_access_key = YOUR_SECRET_ACCESS_KEY
    ```

    - `~/.aws/config`:

    ```
    [default]
    region = YOUR_AWS_REGION
    ```
    (where your [AWS region](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html) is the code us-east-1, eu-west-3 etc.)

## Usage

Here is an example on how to test these programs.

- From the directory cmd/, execute `go run launchEC2_test/main.go` to launch an instance with default values (you can change the default values in the file launchEC2_test/main.go).

    Please note that the access key associated to your instance needs to be in folder cmd/ to work (and if you choose to create one, it will automatically be downloaded to folder cmd/).

- From the directory cmd/, execute `go run deleteEC2_test/main.go`. You will be able to delete the instances of your choice. 

    To delete instances created with the program launchEC2_test/main.go, simply select "4" (delete instances by giving the name) and enter "myEC2instance" (default name given in the previous program).