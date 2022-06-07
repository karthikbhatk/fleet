package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	flags "github.com/jessevdk/go-flags"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

type OptionsStruct struct {
	LambdaExecutionEnv string `long:"lambda-execution-environment" env:"AWS_EXECUTION_ENV"`
	LifecycleTable     string `long:"dynamodb-lifecycle-table" env:"DYNAMODB_LIFECYCLE_TABLE" required:"true"`
	InstanceID         string `long:"instance-id" env:"INSTANCE_ID"`
}

var options = OptionsStruct{}

type LifecycleRecord struct {
	ID    string
	State string
}

func runCmd(args []string) error {
	cmd := exec.Cmd{
		Path:   "/build/terraform",
		Dir:    "/build/deploy_terraform",
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Args:   append([]string{"/build/terraform"}, args...),
	}
	log.Printf("%+v\n", cmd)
	return cmd.Run()
}

func initTerraform() error {
	err := runCmd([]string{
		"init",
		"-backend-config=backend.conf",
	})
	return err
}

func runTerraform(workspace string) error {
	err := runCmd([]string{
		"workspace",
		"new",
		workspace,
	})
	if err != nil {
		return err
	}
	err = runCmd([]string{
		"destroy",
		"-auto-approve",
		"-no-color",
	})
	return err
}

func handler(ctx context.Context, name NullEvent) error {
	if err := initTerraform(); err != nil {
		return err
	}
	if err := runTerraform(fmt.Sprintf("t%s", uuid.New().String()[:8]), redisDatabase); err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Get config from environment
	parser := flags.NewParser(&options, flags.Default)
	if _, err = parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		} else {
			log.Fatal(err)
		}
	}
	if options.LambdaExecutionEnv == "AWS_Lambda_go1.x" {
		lambda.Start(handler)
	} else {
		if err = handler(context.Background(), NullEvent{}); err != nil {
			log.Fatal(err)
		}
	}
}
