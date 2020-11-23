package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
)

const (
	awsTag     = "awsbi-module"
	moduleName = "awsbi-module"
	awsRegion  = "eu-central-1"
	sshKeyName = "vms_rsa"
	retries    = 30
)

var (
	imageTag                  string
	awsAccessKey              string
	awsSecretKey              string
	k8sHostPath               string
	k8sVolPath                string
	stateFilePath             = "shared/state.yml"
	sharedAbsoluteFilePath, _ = filepath.Abs("./shared")
	mountDir                  = sharedAbsoluteFilePath + ":/shared"
	dockerExecPath, _         = exec.LookPath("docker")
)

func TestMain(m *testing.M) {
	cleanupDiskTestStructure()
	cleanupAWSResources()
	setup()
	log.Println("Run tests")
	exitVal := m.Run()
	log.Println("Finish test")
	cleanupDiskTestStructure()
	cleanupAWSResources()
	os.Exit(exitVal)
}

func TestOnInitWithDefaultsShouldCreateProperFileAndFolder(t *testing.T) {
	// given
	expectedFileContentRegexp := "kind: state\nawsbi:\n  status: initialized"

	// when
	_, stderr := runDocker(t, "init", "M_NAME="+moduleName)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	data, err := ioutil.ReadFile(stateFilePath)

	if err != nil {
		t.Fatal("Cannot read state file: ", stateFilePath, err)
	}

	fileContent := string(data)

	matched, err := regexp.MatchString(expectedFileContentRegexp, fileContent)
	if err != nil {
		t.Fatal("There was an error matching expression: ", err)
	}

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedFileContentRegexp, "\nbut found:\n", fileContent)
	}

}

func TestOnPlanWithDefaultsShouldDisplayPlan(t *testing.T) {
	// given
	expectedOutputRegexp := ".*Plan: 14 to add, 0 to change, 0 to destroy.*"

	// when
	stdout, stderr := runDocker(t, "plan", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, err := regexp.MatchString(expectedOutputRegexp, outStr)
	if err != nil {
		t.Fatal("There was an error matching expression: ", err)
	}

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}

}

func TestOnApplyShouldCreateEnvironment(t *testing.T) {
	// given
	expectedOutputRegexp := ".*Apply complete! Resources: 14 added, 0 changed, 0 destroyed.*"

	// when
	stdout, stderr := runDocker(t, "apply", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, err := regexp.MatchString(expectedOutputRegexp, outStr)
	if err != nil {
		t.Fatal("There was an error matching expression: ", err)
	}

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}

	checkNumberOfVms(t)
}

// checks if the proper number of ec2s has been created
func checkNumberOfVms(t *testing.T) {
	// given
	instancesNumber := 1

	newSession, errSession := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if errSession != nil {
		t.Fatal("Cannot get session.", errSession)
	}

	// when

	// ec2
	ec2Client := ec2.New(newSession)

	ec2DescInp := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(awsTag),
				},
			},
		},
	}

	ec2Result, err := ec2Client.DescribeInstances(ec2DescInp)
	if err != nil {
		t.Fatal("There was an error. ", err)
	}

	// then
	if len(ec2Result.Reservations[0].Instances) != 1 {
		t.Error("Expected ", instancesNumber, "instance, got ", len(ec2Result.Reservations[0].Instances))
	}

}

func TestOnDestroyPlanShouldDisplayDestroyPlan(t *testing.T) {
	// given
	expectedOutputRegexp := "Plan: 0 to add, 0 to change, 14 to destroy"

	// when
	stdout, stderr := runDocker(t, "plan-destroy", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, err := regexp.MatchString(expectedOutputRegexp, outStr)
	if err != nil {
		t.Fatal("There was an error matching expression: ", err)
	}

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}
}

func TestOnDestroyShouldDestroyEnvironment(t *testing.T) {
	// given
	expectedOutputRegexp := "Apply complete! Resources: 0 added, 0 changed, 14 destroyed."

	// when
	stdout, stderr := runDocker(t, "destroy", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, err := regexp.MatchString(expectedOutputRegexp, outStr)
	if err != nil {
		t.Fatal("There was an error matching expression: ", err)
	}

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}
}

// initializes test with creation of key pair and checks if variables need to run tests are setup
func setup() {
	log.Println("Initialize test")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	if len(awsAccessKey) == 0 {
		log.Fatalf("expected non-empty AWS_ACCESS_KEY_ID environment variable")
	}
	awsAccessKey = "M_AWS_ACCESS_KEY=" + awsAccessKey

	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if len(awsSecretKey) == 0 {
		log.Fatalf("expected non-empty AWS_SECRET_ACCESS_KEY environment variable")
	}
	awsSecretKey = "M_AWS_SECRET_KEY=" + awsSecretKey

	imageTag = os.Getenv("AWSBI_IMAGE_TAG")
	if len(imageTag) == 0 {
		log.Fatalf("expected non-empty AWSBI_IMAGE_TAG environment variable")
	}

	k8sHostPath  = os.Getenv("K8S_HOST_PATH")
	k8sVolPath   = os.Getenv("K8S_VOL_PATH")

	if ( len(k8sHostPath) != 0 && len(k8sVolPath) != 0) {
		sharedAbsoluteFilePath = k8sVolPath
		mountDir               = k8sHostPath + ":/shared"
		stateFilePath          = k8sVolPath + "/state.yml"
	}

	err := os.MkdirAll(sharedAbsoluteFilePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generating Keys")
	err = generateRsaKeyPair(sharedAbsoluteFilePath, sshKeyName)
	if err != nil {
		log.Fatal(err)
	}
}

// cleans up artifacts from test build from disk
func cleanupDiskTestStructure() {
	log.Println("Starting cleanup.")
	err := os.RemoveAll(sharedAbsoluteFilePath + "/*")
	if err != nil {
		log.Fatal("Cannot remove data folder. ", err)
	}
	log.Println("Cleanup finished.")
}

// cleans up AWS resources if module couldn't clean up resources properly during the test
func cleanupAWSResources() {

	newSession, errSession := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if errSession != nil {
		log.Fatal("Cannot get session.", errSession)
	}

	rgClient := resourcegroups.New(newSession)

	rgName := "rg-" + moduleName
	kpName := "kp-" + moduleName
	eipName := "eip-" + moduleName

	rgResourcesList, errResourcesList := rgClient.ListGroupResources(&resourcegroups.ListGroupResourcesInput{
		GroupName: aws.String(rgName),
	})

	if errResourcesList != nil {
		if aerr, ok := errResourcesList.(awserr.Error); ok {
			log.Println(aerr.Code())
			if aerr.Code() == "NotFoundException" {
				log.Println("Resource group: ", rgName, " not found.")
			} else {
				log.Fatal("Resource group: Cannot get list of resources: ", errResourcesList)
			}
		} else {
			log.Fatal("Resource group: There was an error: ", errResourcesList.Error())
		}
	}

	resourcesTypesToRemove := []string{"Instance", "SecurityGroup", "NatGateway", "EIP", "InternetGateway", "Subnet", "RouteTable", "VPC"}

	for _, resourcesTypeToRemove := range resourcesTypesToRemove {

		filtered := make([]*resourcegroups.ResourceIdentifier, 0)
		for _, element := range rgResourcesList.ResourceIdentifiers {
			s := strings.Split(*element.ResourceType, ":")
			if s[4] == resourcesTypeToRemove {
				filtered = append(filtered, element)
			}

		}

		switch resourcesTypeToRemove {
		case "Instance":
			log.Println("Instance.")
			removeEc2s(newSession, filtered)
		case "EIP":
			log.Println("Releasing public EIPs.")
			releaseAddresses(newSession, eipName)
		case "RouteTable":
			log.Println("RouteTable.")
			removeRouteTables(newSession, filtered)
		case "InternetGateway":
			log.Println("InternetGateway.")
			removeInternetGateway(newSession, filtered)
		case "SecurityGroup":
			log.Println("SecurityGroup.")
			removeSecurityGroup(newSession, filtered)
		case "NatGateway":
			log.Println("NatGateway.")
			removeNatGateways(newSession, filtered)
		case "Subnet":
			log.Println("Subnet.")
			removeSubnet(newSession, filtered)
		case "VPC":
			log.Println("VPC.")
			removeVpc(newSession, filtered)
		}
	}

	removeResourceGroup(newSession, rgName)

	removeKeyPair(newSession, kpName)

}

// run docker with image tag and mounts storage from mountDir with imageTag and other parameters
func runDocker(t *testing.T, params ...string) (bytes.Buffer, bytes.Buffer) {
	var stdout, stderr bytes.Buffer

	commandWithParams := []string{dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag}

	commandWithParams = append(commandWithParams, params...)

	command := &exec.Cmd{
		Path:   commandWithParams[0],
		Args:   commandWithParams,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	if err := command.Run(); err != nil {
		t.Fatal("There was an error running command:", err)
	}

	t.Log("Stdout: ", string(stdout.Bytes()))
	t.Log("Stderr: ", string(stderr.Bytes()))

	return stdout, stderr
}

// generateRsaKeyPair function generates RSA public and private keys and returns
// error if the operation was unsuccessful.
func generateRsaKeyPair(directory, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, name+".pub"), publicKeyBytes, 0644)
}

// removes ec2s using AWS session based on resource identifiers that belong to resource group
func removeEc2s(session *session.Session, ec2sToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, ec2ToRemove := range ec2sToRemove {

		ec2ToRemoveID := strings.Split(*ec2ToRemove.ResourceArn, "/")[1]
		log.Println("EC2: Removing instance with ID: ", ec2ToRemoveID)

		ec2DescInp := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{&ec2ToRemoveID},
		}

		outDesc, errDesc := ec2Client.DescribeInstances(ec2DescInp)
		if errDesc != nil {
			log.Fatalf("EC2: Describe error: %s", errDesc)
		}
		log.Printf("EC2: Describe output: %s", outDesc)

		if outDesc.Reservations != nil {

			instancesToTerminateInp := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{&ec2ToRemoveID},
			}

			outputTerm, errTerm := ec2Client.TerminateInstances(instancesToTerminateInp)
			if errTerm != nil {
				log.Fatalf("EC2: Terminate error: %s", outputTerm)
			}
			log.Printf("EC2: Terminate output: %s", outputTerm)

			errWait := ec2Client.WaitUntilInstanceTerminated(ec2DescInp)
			if errWait != nil {
				log.Fatalf("EC2: Waiting for termination error: %s", errWait)
			}
		}

	}

}

// removes route tables using AWS session based on resource identifiers that belong to resource group
func removeRouteTables(session *session.Session, rtsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, rtToRemove := range rtsToRemove {
		rtIDToRemove := strings.Split(*rtToRemove.ResourceArn, "/")[1]
		log.Println("RouteTable: rtIDToRemove: ", rtIDToRemove)

		rtToDeleteInp := &ec2.DeleteRouteTableInput{
			RouteTableId: &rtIDToRemove,
		}

		output, err := ec2Client.DeleteRouteTable(rtToDeleteInp)

		if err != nil {
			log.Fatal("RouteTable: Deleting route table error: ", err)
		}

		log.Println("RouteTable: Deleting route table: ", output)
	}

}

// removes security groups using AWS session based on resource identifiers that belong to resource group
func removeSecurityGroup(session *session.Session, sgsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, sgToRemove := range sgsToRemove {
		sgIDToRemove := strings.Split(*sgToRemove.ResourceArn, "/")[1]
		log.Println("Security Group: sgIdToRemove: ", sgIDToRemove)

		secGrpInp := &ec2.DeleteSecurityGroupInput{GroupId: &sgIDToRemove}

		output, err := ec2Client.DeleteSecurityGroup(secGrpInp)
		if err != nil {
			log.Fatal("Security Group: Deleting security group error: ", err)
		}

		log.Println("Security Group: Deleting security group: ", output)
	}

}

// removes internet gateways using AWS session based on resource identifiers that belong to resource group
func removeInternetGateway(session *session.Session, igsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, igToRemove := range igsToRemove {
		igIDToRemove := strings.Split(*igToRemove.ResourceArn, "/")[1]
		log.Println("Internet Gateway: igIdToRemove: ", igIDToRemove)

		igDescribeInp := &ec2.DescribeInternetGatewaysInput{
			InternetGatewayIds: []*string{&igIDToRemove},
		}

		descOut, descErr := ec2Client.DescribeInternetGateways(igDescribeInp)

		if descErr != nil {
			log.Fatal("Internet Gateway: Describing internet gateway error: ", descErr)
		}
		log.Println("Internet Gateway: Describing internet gateway: ", descOut)
		vpcID := *descOut.InternetGateways[0].Attachments[0].VpcId

		igDetachInp := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: &igIDToRemove,
			VpcId:             &vpcID,
		}

		detachOut, detachErr := ec2Client.DetachInternetGateway(igDetachInp)
		if detachErr != nil {
			log.Fatal("Internet Gateway: Detaching internet gateway error: ", detachErr)
		}
		log.Println("Internet Gateway: Detaching internet gateway: ", detachOut)

		igDeleteInp := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: &igIDToRemove,
		}

		delOut, delErr := ec2Client.DeleteInternetGateway(igDeleteInp)
		if delErr != nil {
			log.Fatal("Internet Gateway: Deleting internet gateway error: ", delErr)
		}
		log.Println("Internet Gateway: Deleting internet gateway: ", delOut)
	}
}

// removes natgateways using AWS session based on resource identifiers that belong to resource group
func removeNatGateways(session *session.Session, ngsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, ngToRemove := range ngsToRemove {
		ngIDToRemove := strings.Split(*ngToRemove.ResourceArn, "/")[1]
		log.Println("Nat Gateway: ngIdToRemove: ", ngIDToRemove)
		removeSingleNatGatewayWithRetries(ec2Client, ngIDToRemove)
	}

}

// remove single nat gateway based on nat gateway ID using ec2 client
func removeSingleNatGatewayWithRetries(ec2Client *ec2.EC2, ngIDToRemove string) {

	found := true

	for retry := 0; retry <= retries && found; retry++ {

		found = describeNatGateway(ec2Client, ngIDToRemove)

		if found == false {
			continue
		}

		found = removeNatGateway(ec2Client, ngIDToRemove)

		if found == false {
			continue
		}

		waitForNatGatewayDelete(ec2Client, ngIDToRemove)

		log.Println("Nat Gateway: Deleting NAT Gateway. ", ngIDToRemove, " Retry: ", retry)
		time.Sleep(5 * time.Second)
	}
}

// describe Nat Gateway based on ec2 client and Nat Gateway id, returns false if Nat Gateway not found or deleted
func describeNatGateway(ec2Client *ec2.EC2, ngIDToDescribe string) bool {
	descInp := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{&ngIDToDescribe},
	}

	outDesc, errDesc := ec2Client.DescribeNatGateways(descInp)
	if errDesc != nil {
		log.Println(errDesc)
		if aerr, ok := errDesc.(awserr.Error); ok {
			if aerr.Code() == "NatGatewayNotFound" {
				log.Println("Nat Gateway: Nat Gateway not found.")
				return false
			} else {
				log.Fatal("Nat Gateway: Describe error: ", errDesc)
			}
		} else {
			log.Fatal("Nat Gateway: There was an error: ", errDesc.Error())
		}
	}
	log.Printf("Nat Gateway: Describe output: %s", outDesc)

	if outDesc.NatGateways == nil || *outDesc.NatGateways[0].State == "deleted" {
		log.Print("Nat Gateway: Element not found or has been already deleted.")
		return false
	}
	return true
}

// appropriate Nat Gateway delete method based on ec2 client and Nat Gateway id, returns false if Nat Gateway not found
func removeNatGateway(ec2Client *ec2.EC2, ngIDToRemove string) bool {
	ngDelInp := &ec2.DeleteNatGatewayInput{
		NatGatewayId: &ngIDToRemove,
	}

	_, err := ec2Client.DeleteNatGateway(ngDelInp)

	if err != nil {
		log.Println("Nat Gateway: Error: ", err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NatGatewayNotFound" {
				log.Print("Nat Gateway: Element not found.", err)
				return false
			}
			if aerr.Code() != "ResourceNotReady" {
				log.Fatal("Nat Gateway: Deleting NAT Gateway: ", err)
			}
		} else {
			log.Fatal("Nat Gateway: Deleting NAT Gateway: ", err.Error())
		}

	}
	return true
}

// wait for Nat Gateway to be removed based on ec2 client and Nat Gateway id
func waitForNatGatewayDelete(ec2Client *ec2.EC2, ngIDToWait string) {
	descInp := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{&ngIDToWait},
	}

	errWait := ec2Client.WaitUntilNatGatewayAvailable(descInp)
	if errWait != nil {
		if aerr, ok := errWait.(awserr.Error); ok {
			if aerr.Code() != "ResourceNotReady" {
				log.Fatal("Nat Gateway: Wait error: ", errWait)
			}
		} else {
			log.Fatal("Nat Gateway: There was an error: ", errWait.Error())
		}
	}
}

// removes subnets using AWS session based on resource identifiers that belong to resource group
func removeSubnet(session *session.Session, subnetsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, subnetToRemove := range subnetsToRemove {
		subnetIDToRemove := strings.Split(*subnetToRemove.ResourceArn, "/")[1]
		log.Println("Subnet: subnetIdToRemove: ", subnetIDToRemove)

		subnetInp := &ec2.DeleteSubnetInput{
			SubnetId: &subnetIDToRemove,
		}

		output, err := ec2Client.DeleteSubnet(subnetInp)
		if err != nil {
			log.Fatal("Subnet: Deleting subnet error: ", err)
		}
		log.Println("Subnet: Deleting subnet: ", output)
	}

}

// removes vpcs using AWS session based on resource identifiers that belong to resource group
func removeVpc(session *session.Session, vpcsToRemove []*resourcegroups.ResourceIdentifier) {

	ec2Client := ec2.New(session)

	for _, vpcToRemove := range vpcsToRemove {
		vpcIDToRemove := strings.Split(*vpcToRemove.ResourceArn, "/")[1]
		log.Println("VPC: vpcIdToRemove: ", vpcIDToRemove)

		vpcToDeleteInp := &ec2.DeleteVpcInput{
			VpcId: &vpcIDToRemove,
		}

		output, err := ec2Client.DeleteVpc(vpcToDeleteInp)
		if err != nil {
			log.Println("VPC: Delete VPC error: ", err)
		}
		log.Println("VPC: Delete VPC: ", output)
	}
}

// removes key pairs using AWS session based on resource identifiers that belong to resource group
func removeKeyPair(session *session.Session, kpName string) {

	ec2Client := ec2.New(session)

	removeKeyInp := &ec2.DeleteKeyPairInput{
		KeyName: &kpName,
	}

	output, err := ec2Client.DeleteKeyPair(removeKeyInp)
	if err != nil {
		log.Fatal("Key Pair: Deleting key pair error: ", err)
	}
	log.Println("Key Pair: Deleting key pair: ", output)
}

// removes ec2s using AWS session based on resource identifiers that belong to resource group
func releaseAddresses(session *session.Session, eipName string) {

	ec2Client := ec2.New(session)

	describeEips, err := ec2Client.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		log.Fatal("EIP: Cannot get EIP list.", err)
	}

	for _, eip := range describeEips.Addresses {

		tagPresent := false
		if eip.Tags != nil {
			tagPresent = checkIfTagPresent(eipName, eip.Tags)
			if tagPresent {
				log.Printf("EIP: Releasing EIP with AllocationId: %s", *eip.AllocationId)

				eipToReleaseInp := &ec2.ReleaseAddressInput{
					AllocationId: eip.AllocationId,
				}

				found := true

				for retry := 0; retry <= retries && found; retry++ {
					_, err := ec2Client.ReleaseAddress(eipToReleaseInp)
					if err != nil {
						if aerr, ok := err.(awserr.Error); ok {
							if aerr.Code() == "InvalidAllocationID.NotFound" {
								log.Print("EIP: Element not found.", err)
								found = false
								continue
							}
							if aerr.Code() != "AuthFailure" && aerr.Code() != "InvalidAllocationID.NotFound" {
								log.Fatal("EIP: Releasing EIP error: ", err)
							}
						} else {
							log.Fatal("EIP: There was an error: ", err.Error())
						}
					}
					log.Println("EIP: Releasing EIP. Retry: ", retry)
					time.Sleep(5 * time.Second)
				}
			}
		}

	}

}

// removes resource groups using AWS session based on name
func removeResourceGroup(session *session.Session, rgToRemoveName string) {

	rgClient := resourcegroups.New(session)

	log.Println("Resource Group: Removing resource group: ", rgToRemoveName)
	rgDelInp := resourcegroups.DeleteGroupInput{
		GroupName: aws.String(rgToRemoveName),
	}
	rgDelOut, rgDelErr := rgClient.DeleteGroup(&rgDelInp)
	if rgDelErr != nil {
		if aerr, ok := rgDelErr.(awserr.Error); ok {
			if aerr.Code() == "NotFoundException" {
				log.Println("Resource Group: Resource group not found. ")
			} else {
				log.Fatal("Resource Group: Deleting resource group error: ", rgDelErr)
			}
		} else {
			log.Fatal("Resource Group: There was an error: ", rgDelErr.Error())
		}

	} else {
		log.Println("Resource Group: Deleting resource group: ", rgDelOut)
	}

}

// checks if the tag is present in tag list
func checkIfTagPresent(toFind string, tags []*ec2.Tag) bool {

	for _, tag := range tags {
		if *tag.Value == toFind {
			return true
		}
	}
	return false
}
