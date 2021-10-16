package controllers_test

import (
	"fmt"
	"testing"

	"github.com/markliederbach/qrkdns/pkg/controllers"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	successCommand cli.Command = cli.Command{
		Name: "successful",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func TestQrkDNSApp(t *testing.T) {
	// disable help text for tests
	cli.AppHelpTemplate = ""

	tests := []testRunner{
		{
			testCase: "qrkdns app runs successfully with no args",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				app := controllers.NewQrkDNSApp("version123", []*cli.Command{&successCommand})
				err := app.Run([]string{"qrkdns"})
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			testCase: "qrkdns app returns error for bad log level",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)
				app := controllers.NewQrkDNSApp("version123", []*cli.Command{&successCommand})
				err := app.Run([]string{"qrkdns", fmt.Sprintf("--%v", controllers.LogLevelFlag), "FOO"})
				_, expectedErr := logrus.ParseLevel("FOO")
				g.Expect(err).To(MatchError(expectedErr))
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
