package sample

import "fmt"

//Not Implemented
func ArchiveSample() {
	DeleteTestBucketAndObject()

	defer DeleteTestBucketAndObject()

	// TODO : ArchiveSample archives sample

	fmt.Printf("ArchiveSample Run Success !\n\n")
}
