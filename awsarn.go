package awsarn

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

var (
	// ErrMalformed is returned when the ARN appears to be invalid.
	ErrMalformed = errors.New("malformed ARN")

	// ErrVariablesNotSupported is returned when the ARN contains policy
	// variables.
	ErrVariablesNotSupported = errors.New("policy variables are not supported")
)

// Components encapsulate the individual pieces of an AWS ARN.
type Components struct {
	ARN               string
	Partition         string
	Service           string
	Region            string
	AccountID         string
	ResourceType      string
	Resource          string
	ResourceDelimiter string
}

// Validate checks if an input ARN string conforms to a format which can be
// parsed by this package.
func Validate(arn string) error {
	pieces := strings.SplitN(arn, ":", 6)
	return validate(arn, pieces)
}

func validate(arn string, pieces []string) error {
	if strings.Contains(arn, "${") {
		return ErrVariablesNotSupported
	}
	if len(pieces) < 6 {
		return ErrMalformed
	}
	return nil
}

// Parse accepts and ARN string and attempts to break it into constiuent parts.
func Parse(arn string) (*Components, error) {
	pieces := strings.SplitN(arn, ":", 6)

	if err := validate(arn, pieces); err != nil {
		return nil, err
	}

	components := &Components{
		ARN:       pieces[0],
		Partition: pieces[1],
		Service:   pieces[2],
		Region:    pieces[3],
		AccountID: pieces[4],
	}
	if n := strings.Count(pieces[5], ":"); n > 0 {
		components.ResourceDelimiter = ":"
		resourceParts := strings.SplitN(pieces[5], ":", 2)
		components.ResourceType = resourceParts[0]
		components.Resource = resourceParts[1]
	} else {
		if m := strings.Count(pieces[5], "/"); m == 0 {
			components.Resource = pieces[5]
		} else {
			components.ResourceDelimiter = "/"
			resourceParts := strings.SplitN(pieces[5], "/", 2)
			components.ResourceType = resourceParts[0]
			components.Resource = resourceParts[1]
		}
	}
	return components, nil
}

// String rebuilds the original input ARN to res
func (c Components) String() string {
	front := []string{
		c.ARN,
		c.Partition,
		c.Service,
		c.Region,
		c.AccountID,
	}
	s := strings.Join(front, ":") + ":" + c.ResourceChunk()
	return s
}

// ResourceChunk returns the tail-end of the ARN, containing resource
// specifications.
func (c Components) ResourceChunk() string {
	return c.ResourceType + c.ResourceDelimiter + c.Resource
}

// SupersetOf returns true if c is a superset of the other passed components.
func (c Components) SupersetOf(other *Components) bool {
	if reflect.DeepEqual(c, other) {
		return true
	}

	type pair struct {
		a string
		b string
	}

	pairs := []pair{
		{c.ARN, other.ARN},
		{c.Partition, other.Partition},
		{c.Service, other.Service},
		{c.Region, other.Region},
		{c.AccountID, other.AccountID},
		{c.ResourceChunk(), other.ResourceChunk()},
	}

	for _, p := range pairs {
		if !toRegExp(p.a).MatchString(p.b) {
			return false
		}
	}
	return true
}

// toRegexp takes an AWS ARN resource component and converts it to a
// go-compatible regular expression.  The '*' and '?' characters have special
// wildcard meanings in this context.
func toRegExp(s string) *regexp.Regexp {
	// Escape the input string.
	s = regexp.QuoteMeta(s)
	// Unescape special * and ? characters because they have special meanings.
	s = strings.Replace(s, "\\*", ".*", -1)
	s = strings.Replace(s, "\\?", ".", -1)
	expr := regexp.MustCompile(s)
	return expr
}
