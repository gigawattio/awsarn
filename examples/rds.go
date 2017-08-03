package main

import (
	"fmt"

	"github.com/gigawattio/awsarn"
)

const arn = "arn:aws:rds:region:account-id:db:db-instance-name"

func main() {
	components, err := awsarn.Parse(arn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%# v\n", *components)
	eq := components.String() == arn
	fmt.Printf("Reconstruction: %v, equal=%v\n", components.String(), eq)
}
