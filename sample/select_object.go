package sample

import "fmt"

//Not Implemented
func SelectObjectSample() {
	DeleteTestBucketAndObject()

	defer DeleteTestBucketAndObject()

	// TODO : SelectObjectSample shows how to get data from csv/json object by sql
	fmt.Printf("SelectObjectSample Run Success !\n\n")
}
