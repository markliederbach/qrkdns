package main

import (
	"testing"
	// . "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func TestMain(t *testing.T) {

	tests := []testRunner{
		{
			testCase: "main runs successfully",
			runner: func(tt *testing.T) {
				// TODO
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
