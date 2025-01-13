# AWS

This code demonstrates how to use AWS SDK for Go v2 to perform programmatic actions on AWS EC2 and S3 (ie, launch and delete EC2 instances, create EC2 key pairs, list S3 buckets, from code only).

## Prerequisites

### Go prerequisites
- Go: the v2 SDK requires a minimum version of Go 1.21 (and this code was written with version 1.22.4).

- After cloning repository please `go mod tidy` to retrieve the required Go dependencies.

### AWS prerequisites

Of course, you will need a AWS account.

To execute the code, you need to create a user with permission to access AWS EC2 and S3. Here is how to do it:

- In the AWS IAM console, go to [Users](https://us-east-1.console.aws.amazon.com/iam/home?region=us-east-1#/users). Create a user and attach the policies [AmazonEC2FullAccess](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AmazonEC2FullAccess.html) and [AmazonS3FullAccess](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AmazonS3FullAccess.html).

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

### EC2

- From the directory cmd/, execute `go run launchEC2_test/main.go` to launch an instance with default values (you can change the default values in the file launchEC2_test/main.go).

    Please note that the access key associated to your instance needs to be in folder cmd/ to work (and if you choose to create one, it will automatically be downloaded to folder cmd/).

- From the directory cmd/, execute `go run deleteEC2_test/main.go`. You will be able to delete the instances of your choice. 

    To delete instances created with the program launchEC2_test/main.go, simply select "4" (delete instances by giving the name) and enter "myEC2instance" (default name given in the previous program).

### S3

- From the directory cmd/, execute `go run s3_test/main.go` to perform operations on S3.