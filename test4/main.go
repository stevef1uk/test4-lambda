package main

import (
	//"context"
	"errors"
	"fmt"
	"strconv"

	//request2 "github.com/aws/aws-sdk-go/aws/request"
	"github.com/stevef1uk/test4/data"
	"github.com/stevef1uk/test4/models"
	mylambda "github.com/stevef1uk/test4/lambda"
	"github.com/stevef1uk/test4/restapi/operations"
	//"github.com/stevef1uk/test4/restapi/operations/verysimple"
	"io/ioutil"
	"encoding/json"
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

	switch request.HTTPMethod {
	case "GET":
		params := operations.GetVerysimpleParams{}
		fmt.Printf("Request ID = %s", request.QueryStringParameters["ID"])
		i, _ := strconv.Atoi( request.QueryStringParameters["ID"])
		params.ID = int32(i)
		ret := data.SearchNew( params )
		fmt.Printf("Called SearchNew with result = %s", ret)
		return events.APIGatewayProxyResponse{
			Body:        ret ,
			StatusCode: 200,
		}, nil
	case "POST":
		fmt.Printf("In POST Body = %s", request.Body)
		params := models.Verysimple{}
		//params = models.UnmarshalBinary( []byte request.Body )
		b := [] byte (request.Body)
		_ = json.Unmarshal( b, &params )
		ret, err := data.InsertNew( &params )
		status_code := 200
		if err {
			status_code = 201
		}
		return events.APIGatewayProxyResponse{
			Body:        ret ,
			StatusCode: status_code,
		}, nil

	default:
		fmt.Printf("Methid called %s not handled!", request.HTTPMethod )
		return events.APIGatewayProxyResponse{
			Body:        "Not implemented",
			StatusCode: 405,
		}, nil
	}
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
