package handlers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"github.com/MiG-21/go-sso/internal/web/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Healthcheck", func() {
	It("HealthPingHandler", func() {
		resp, err := app.Test(httptest.NewRequest("GET", "/v1/healthcheck/ping", nil))
		Expect(err).NotTo(HaveOccurred())

		actual := types.HealthCheckPing{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &actual)
		Expect(err).NotTo(HaveOccurred())
		expected := types.HealthCheckPing{Ping: "PONG"}
		Expect(expected).To(Equal(actual))
	})

	It("HealthPingHandler", func() {
		resp, err := app.Test(httptest.NewRequest("GET", "/v1/healthcheck/info", nil))
		Expect(err).NotTo(HaveOccurred())

		actual := types.HealthCheckInfo{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &actual)
		Expect(err).NotTo(HaveOccurred())
		expected := types.HealthCheckInfo{
			AppName:     "SomeApp",
			AppVersion:  "SomeVersion",
			ClusterName: "SomeCluster",
			Git: types.HealthCheckInfoGit{
				Hash: "SomeGitHash",
				Ref:  "SomeGitBranch",
				Url:  "SomeGitUrl",
			},
		}
		Expect(expected).To(Equal(actual))
	})
})
