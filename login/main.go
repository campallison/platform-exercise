package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	fenderAuth "github.com/campallison/platform-exercise"
)

func main() {
	lambda.Start(fenderAuth.LoginHandler)
}
