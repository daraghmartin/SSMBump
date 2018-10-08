package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"os"
	"strconv"
	"strings"
)

const newVersion string = "0.0.0"

func main() {
	keyName := os.Args[1]
	version := getSSMVersion(keyName)
	fmt.Println(version)
}

func createNewParam(keyName string, value string, ssmc *ssm.SSM) bool {
	_, puterr := ssmc.PutParameter(&ssm.PutParameterInput{
		Name:  &keyName,
		Type:  aws.String("String"),
		Value: aws.String(newVersion),
	})

	if puterr != nil {
		fmt.Println("Error creating new param: ", puterr)
		os.Exit(2)
	}

	return true
}

func updateParam(keyName string, value string, overwrite bool, ssmc *ssm.SSM) bool {
	_, puterr := ssmc.PutParameter(&ssm.PutParameterInput{
		Name:      &keyName,
		Type:      aws.String("String"),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
	})

	if puterr != nil {
		fmt.Println("Error updating param: ", puterr)
		os.Exit(1)
	}

	return true
}

func bump(version string) string {
	s := strings.Split(version, ".")
	major, minor, patch := s[0], s[1], s[2]
	i, _ := strconv.Atoi(patch)
	return fmt.Sprintf("%s.%s.%d", major, minor, i+1)
}

func getSSMVersion(keyName string) string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ssmc := ssm.New(sess)

	withDecryption := false

	param, err := ssmc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyName,
		WithDecryption: &withDecryption,
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "ParameterNotFound":
				createNewParam(keyName, newVersion, ssmc)
				return newVersion
			default:
				fmt.Println(awsErr)
				os.Exit(3)
			}
		}
	}

	version := bump(*param.Parameter.Value)
	updateParam(keyName, version, true, ssmc)
	return version
}
