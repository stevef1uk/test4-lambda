package main

import (
	"errors"
	"fmt"
	//request2 "github.com/aws/aws-sdk-go/aws/request"
	"github.com/stevef1uk/test4/data"
	mylambda "github.com/stevef1uk/test4/lambda"
	"github.com/stevef1uk/test4/restapi/operations"
	"io/ioutil"
	//"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
	initialised = false;
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	params := operations.GetVerysimpleParams{}
	params.ID = 1
	ret := data.SearchNew( params )
	fmt.Printf("Called SearchNew with result = %s", ret)


	return events.APIGatewayProxyResponse{
		Body:        ret ,
		StatusCode: 200,
	}, nil
}

func main() {

	fmt.Println("In Main about to call GetFile" )
	if initialised == false {
		initialised = true
		mylambda.GetFiles( "sjfisher2", "us-east-1")
		data.SetUp()
	}
	lambda.Start(handler)
	data.Stop()
}
