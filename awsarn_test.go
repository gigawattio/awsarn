package awsarn

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{
			input:    "arn:aws:s3:::my_corporate_bucket/exampleobject.png",
			expected: nil,
		},
	}
	for i, test := range tests {
		if expected, actual := test.expected, Validate(test.input); actual != expected {
			t.Errorf("[i=%v] Expected input=%q to produce value=%v but actual=%v", i, test.input, expected, actual)
		}
	}
}

func TestParse(t *testing.T) {
	type expected struct {
		components *Components
		err        error
	}
	tests := []struct {
		input    string
		expected expected
	}{
		{
			input: "arn:partition:service:region:account-id:resource",
			expected: expected{
				components: &Components{
					ARN:       "arn",
					Partition: "partition",
					Service:   "service",
					Region:    "region",
					AccountID: "account-id",
					Resource:  "resource",
				},
			},
		},
		{
			input: "arn:partition:service:region:account-id:resource-type/resource",
			expected: expected{
				components: &Components{
					ARN:               "arn",
					Partition:         "partition",
					Service:           "service",
					Region:            "region",
					AccountID:         "account-id",
					ResourceType:      "resource-type",
					Resource:          "resource",
					ResourceDelimiter: "/",
				},
			},
		},
		{
			input: "arn:partition:service:region:account-id:resource-type:resource",
			expected: expected{
				components: &Components{
					ARN:               "arn",
					Partition:         "partition",
					Service:           "service",
					Region:            "region",
					AccountID:         "account-id",
					ResourceType:      "resource-type",
					Resource:          "resource",
					ResourceDelimiter: ":",
				},
			},
		},
		{
			input: "arn:partition:service:region:",
			expected: expected{
				err: ErrMalformed,
			},
		},
	}
	for i, test := range tests {
		c, err := Parse(test.input)
		if expected, actual := test.expected.err, err; actual != expected {
			t.Errorf("[i=%v] Expected err=%v but actual=%v", i, expected, actual)
		}
		if expected, actual := test.expected.components, c; !reflect.DeepEqual(actual, expected) {
			t.Errorf("[i=%v] Expected components=%# v", i, expected)
			t.Errorf("[i=%v]              actual=%# v", i, actual)
		}
	}
}

func TestSupersetOf(t *testing.T) {
	tests := []struct {
		a            string
		b            string
		aSupersetOfB bool
		bSupersetOfA bool
	}{
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company",
			aSupersetOfB: true,
			bSupersetOfA: true,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: true,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.*",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: true,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.???????",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: true,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.???????/*",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: true,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company:*/*",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: false,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/*/*",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: true,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company:mesos-dcname-platform-environment/*",
			b:            "arn:aws:s3:::service.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: false,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service1.dcname-platform-environment.us-west-2.company",
			b:            "arn:aws:s3:::service2.dcname-platform-environment.us-west-2.company",
			aSupersetOfB: false,
			bSupersetOfA: false,
		},
		{
			a:            "arn:aws:s3:::service1.dcname-platform-environment.us-west-2.company",
			b:            "arn:aws:s3:::service2.dcname-platform-environment.us-west-2.company/mesos-dcname-platform-environment/*",
			aSupersetOfB: false,
			bSupersetOfA: false,
		},
	}
	for i, test := range tests {
		a, err := Parse(test.a)
		if err != nil {
			t.Errorf("[i=%v] Unexpected error parsing %q: %s", i, test.a, err)
			continue
		}
		b, err := Parse(test.b)
		if err != nil {
			t.Errorf("[i=%v] Unexpected error parsing %q: %s", i, test.b, err)
			continue
		}

		if expected, actual := test.aSupersetOfB, a.SupersetOf(b); actual != expected {
			t.Errorf("[i=%v] Expected (a ⊃ b) superset evaluation=%v but actual=%v;\n\n\ta = %q\n\n\tb = %q", i, expected, actual, test.a, test.b)
		}
		if expected, actual := test.bSupersetOfA, b.SupersetOf(a); actual != expected {
			t.Errorf("[i=%v] Expected (b ⊃ a) superset evaluation=%v but actual=%v;\n\n\tb = %q\n\n\ta = %q", i, expected, actual, test.b, test.a)
		}
	}
}

func TestString(t *testing.T) {
	arns := []string{
		":::::",
		"::::::",
		":::::/",
		"::::::foo",
		":::::/foo",
		"arn:aws:s3:::my_corporate_bucket/exampleobject.png",
	}
	for i, arn := range arns {
		components, err := Parse(arn)
		if err != nil {
			t.Errorf("[i=%v] Unexpected error parsing input=%q: %s", i, arn, err)
			continue
		}
		if expected, actual := arn, components.String(); actual != expected {
			t.Errorf("[i=%v] Expected input ARN=%q to produce same output string but actual=%q", i, expected, actual)
		}
	}
}

func TestToRegExp(t *testing.T) {
	tests := []struct {
		input     string
		expectOut string
	}{
		{
			input:     `foo-bar/main.dir/*`,
			expectOut: `foo-bar/main\.dir/.*`,
		},
		{
			input:     `foo-bar/main.dir/?`,
			expectOut: `foo-bar/main\.dir/.`,
		},
		{
			input:     `?/foo-bar/?/*`,
			expectOut: `./foo-bar/./.*`,
		},
		{
			input:     `foo-bar/?/(main.dir/*`,
			expectOut: `foo-bar/./\(main\.dir/.*`,
		},
	}
	for i, test := range tests {
		if expected, actual := test.expectOut, toRegExp(test.input).String(); actual != expected {
			t.Errorf("[i=%v] Expected input=%q production=%q but actual=%q", i, test.input, expected, actual)
		}
	}
}
