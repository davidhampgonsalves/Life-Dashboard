package main

import (
	"encoding/base64"
	"net/http"

	"davidhampgonsalves/lifedashboard/pkg"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func apiResponse(status int, imgData string) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers:         map[string]string{"Content-Type": "application/octet-stream"},
		IsBase64Encoded: true,
		Body:            imgData,
		StatusCode:      status,
	}
	return &resp, nil
}

func generateImg(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	imgBytes := pkg.Generate()
	pngAsBase64 := base64.StdEncoding.EncodeToString(imgBytes)

	return apiResponse(http.StatusOK, pngAsBase64)
}

func main() {
	lambda.Start(generateImg)
}
