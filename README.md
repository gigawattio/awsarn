# awsarn

[![Documentation](https://godoc.org/github.com/gigawattio/awsarn?status.svg)](https://godoc.org/github.com/gigawattio/awsarn)
[![Build Status](https://travis-ci.org/gigawattio/awsarn.svg?branch=master)](https://travis-ci.org/gigawattio/awsarn)
[![Report Card](https://goreportcard.com/badge/github.com/gigawattio/awsarn)](https://goreportcard.com/report/github.com/gigawattio/awsarn)

### About

[awsarn](https://github.com/gigawattio/awsarn) is an ARN parser.

More specifically, this is a Go (golang) library for validating, parsing, and comparing AWS ARN resource identifier strings.

This package also provides the capability of determining if one ARN is a superset of another.  This is useful for safely eliminating redundant ARNs from a set.

Created by [Jay Taylor](https://jaytaylor.com/) and used by [Gigawatt](https://gigawatt.io/).

#### ARN Vocabulary

The AWS documentation uses [two subtly different sets of vocabulary](#further-reading) when discussing the internal workings of ARNs:

Variant #1

```
arn:partition:service:region:account-id:resource
arn:partition:service:region:account-id:resourcetype/resource
arn:partition:service:region:account-id:resourcetype:resource
```

Variant #2

```
arn:partition:service:region:namespace:relative-id
```

This package uses the vocabulary of variant #1, that is:

* arn
* partition
* service
* region
* account-id
* resource, resourcetype/resource, resourcetype:resource

#### Wildcards

The documentation is ambiguous about which components of an ARN allow wildcards like `*` and `?`.  This package uses the loosest possible interpretation, which means wildcards are allowed in any and all parts of ARNs.

### Requirements

* Go version 1.1 or newer

### Example usage

Parse an AWS ARN for an RDS database:

[examples/rds.go](examples/rds.go)

```go
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
```

Output:

```shell
awsarn.Components{
    ARN: "arn",
    Partition: "aws",
    Service: "rds",
    Region: "region",
    AccountID: "account-id",
    ResourceType: "db",
    Resource: "db-instance-name"
    ResourceDelimiter: ":"
}
Reconstruction: arn:aws:rds:region:account-id:db:db-instance-name, equal=true
```

Also may be worth checking out the [unit-tests](awsarn_test.go), too!

### Running the test suite

    go test -v ./...
    echo $?

if `echo $?` produces a 0, that's a clean exit status and means the tests succeeded.  Anything else indicates one or more failed tests.

### Terminology

* ARN: Amazon Resource Name; used for identifying, specifying, and referencing resources
* AWS: Amazon Web Services; Cloud provider

### Components of an ARN

Piece by piece:

arn:partition:service:region:account-id:resourcetype/resource

**arn**

> This should always be the string "arn", indicating the start of an ARN.

**partition**

> The partition that the resource is in. For standard AWS regions, the
> partition is aws. If you have resources in other partitions, the partition
> is aws-partitionname. For example, the partition for resources in the China
> (Beijing) region is aws-cn.

**service**

>The service namespace that identifies the AWS product (for example, Amazon
>S3, IAM, or Amazon RDS). For a list of namespaces, see AWS Service Namespaces.

**region**

> The region the resource resides in. Note that the ARNs for some resources do
> not require a region, so this component might be omitted.

**account**

> The ID of the AWS account that owns the resource, without the hyphens. For
> example, 123456789012. Note that the ARNs for some resources don't require
> an account number, so this component might be omitted.

**resource, resourcetype:resource, or resourcetype/resource**

> The content of this part of the ARN varies by service. It often includes an
> indicator of the type of resource—for example, an IAM user or Amazon RDS
> database —followed by a slash (/) or a colon (:), followed by the resource
> name itself. Some services allows paths for resource names, as described in
> Paths in ARNs.

### Further reading

* [Amazon AWS ARN documentation](http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html)
* [Amazon AWS S3 ARN format](http://docs.aws.amazon.com/AmazonS3/latest/dev/s3-arn-format.html)

### License

Permissive MIT license, see the [LICENSE](LICENSE) file for more information.
