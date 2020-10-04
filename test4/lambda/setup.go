package lambda

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Function to create a file in /tmp
func CreateFile(  fileName string ) (*os.File, bool) {


	fullFileName := "/tmp/" + fileName

	_, err := os.Stat(fullFileName)
	if os.IsNotExist(err) {

		file, err := os.Create(fullFileName)
		if err == nil {
			return file, true
		}
	}
	return nil, false
}


// snippet-end:[s3.go.customHttpClient.import]

// HTTPClientSettings defines the HTTP setting for clients
// snippet-start:[s3.go.customHttpClient_struct]
type HTTPClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
}

// snippet-end:[s3.go.customHttpClient_struct]

// NewHTTPClientWithSettings creates an HTTP client with some custom settings
// Inputs:
//     httpSettings contains some custom HTTP settings for the client
// Output:
//     If success, an HTTP client
//     Otherwise, ???
// snippet-start:[s3.go.customHttpClient_client]
func NewHTTPClientWithSettings(httpSettings HTTPClientSettings) (*http.Client, error) {
	var client http.Client
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			DualStack: true,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
	}

	// So client makes HTTP/2 requests
	err := http2.ConfigureTransport(tr)
	if err != nil {
		return &client, err
	}

	return &http.Client{
		Transport: tr,
	}, nil
}

// Copy file from S3 to /tmp
func getFile ( bucketName string, fileName string)  {
	var body io.ReadCloser
	var err error

	destFile, created := CreateFile(  fileName  )
	if ! created {
		return
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Println("Got error closing file " + err.Error())
		}
	}()


	// Creating a SDK session using the custom HTTP client
	// and use that session to create S3 client.
	// snippet-start:[s3.go.customHttpClient_session]
	httpClient, err := NewHTTPClientWithSettings(HTTPClientSettings{
		Connect:          5 * time.Second,
		ExpectContinue:   1 * time.Second,
		IdleConn:         90 * time.Second,
		ConnKeepAlive:    30 * time.Second,
		MaxAllIdleConns:  100,
		MaxHostIdleConns: 10,
		ResponseHeader:   5 * time.Second,
		TLSHandshake:     5 * time.Second,
	})
	if err != nil {
		fmt.Println("Got an error creating custom HTTP client:")
		fmt.Println(err)
		return
	}

	sess := session.Must(session.NewSession(&aws.Config{
		HTTPClient: httpClient,
	}))

	svc := s3.New(sess)

	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
	})
	if err != nil {
		fmt.Println("Got error calling GetObject:")
		fmt.Println(err.Error())
		return
	}
	body = obj.Body



	// Convert body from IO.ReadCloser to string:
	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(body)
	if err != nil {
		fmt.Println("Got an error reading body of object:")
		fmt.Println(err)
		return
	}

	//body = obj.Body
	//newBytes := buf.String()
	//fmt.Println("Read file, contents = " + newBytes )
	if _, err := destFile.Write(buf.Bytes()); err != nil {
		fmt.Println("Got error writing to file " + err.Error())
	}
}

func GetFiles( bucketName string, regionName string )  {
	fmt.Println("Getting files  from region " +  bucketName + " from region " + regionName )
	os.Setenv("AWS_REGION", regionName )
	getFile( bucketName, "ca.crt" )
	getFile( bucketName, "cert.pfx" )
	getFile( bucketName, "key" )
	getFile( bucketName, "cert" )
	getFile( bucketName, "identity.jks" )
	getFile( bucketName, "trustStore.jks" )
	getFile( bucketName, "config.json" )
}
