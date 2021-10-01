package config_test

import (
	"testing"

	"github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/config"
	. "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func TestFile(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "loads config",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := mocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				conf, err := config.Load()
				g.Expect(err).NotTo(HaveOccurred())

				expectedConf := config.Config{
					CloudFlareAccountID: "foo",
					CloudFlareAPIToken:  "bar",
					LogLevel:            "INFO",
				}

				g.Expect(conf).To(Equal(expectedConf))
			},
		},
		{
			testCase: "returns loading error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := mocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"CLOUDFLARE_ACCOUNT_ID": "", // required
						"CLOUDFLARE_API_TOKEN":  "", // required
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				_, err = config.Load()
				g.Expect(err).To(HaveOccurred())
			},
		},
		{
			testCase: "returns log level error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := mocks.MockEnv{}
				err := env.Load(
					map[string]string{
						"CLOUDFLARE_ACCOUNT_ID": "foo",
						"CLOUDFLARE_API_TOKEN":  "bar",
						"LOG_LEVEL":             "baz", // not valid
					},
				)
				g.Expect(err).NotTo(HaveOccurred())
				defer env.Restore()

				_, err = config.Load()
				g.Expect(err).To(HaveOccurred())
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
