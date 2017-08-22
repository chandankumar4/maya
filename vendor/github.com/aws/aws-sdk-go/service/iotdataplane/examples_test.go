// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package iotdataplane_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleIoTDataPlane_DeleteThingShadow() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := iotdataplane.New(sess)

	params := &iotdataplane.DeleteThingShadowInput{
		ThingName: aws.String("ThingName"), // Required
	}
	resp, err := svc.DeleteThingShadow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleIoTDataPlane_GetThingShadow() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := iotdataplane.New(sess)

	params := &iotdataplane.GetThingShadowInput{
		ThingName: aws.String("ThingName"), // Required
	}
	resp, err := svc.GetThingShadow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleIoTDataPlane_Publish() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := iotdataplane.New(sess)

	params := &iotdataplane.PublishInput{
		Topic:   aws.String("Topic"), // Required
		Payload: []byte("PAYLOAD"),
		Qos:     aws.Int64(1),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleIoTDataPlane_UpdateThingShadow() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := iotdataplane.New(sess)

	params := &iotdataplane.UpdateThingShadowInput{
		Payload:   []byte("PAYLOAD"),       // Required
		ThingName: aws.String("ThingName"), // Required
	}
	resp, err := svc.UpdateThingShadow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
