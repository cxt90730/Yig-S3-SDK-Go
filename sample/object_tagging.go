package sample

import "fmt"

//Not Implemented
func ObjectTaggingSample() {
	DeleteTestBucketAndObject()

	defer DeleteTestBucketAndObject()

	// TODO 1. Set Tagging of object
	// TODO 2. Get Tagging of object
	// TODO 3. Delete Tagging of object
	fmt.Printf("ObjectTaggingSample Run Success !\n\n")
}
