package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var settings struct {
	Buildpacks struct {
		BuildPlan struct {
			Online string
		}
		Procfile struct {
			Online string
		}
	}
	Buildpack struct {
		ID   string
		Name string
	}
	Config struct {
		BuildPlan string `json:"build-plan"`
	}
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())
	defer file.Close()

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())

	file, err = os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&settings.Buildpack)
	Expect(err).NotTo(HaveOccurred())

	root, err := filepath.Abs("./..")
	Expect(err).NotTo(HaveOccurred())

	buildpackStore := occam.NewBuildpackStore()

	settings.Buildpacks.Procfile.Online, err = buildpackStore.Get.
		WithVersion("0.0.0").
		Execute(root)
	Expect(err).ToNot(HaveOccurred())

	settings.Buildpacks.BuildPlan.Online, err = buildpackStore.Get.
		Execute(settings.Config.BuildPlan)
	Expect(err).ToNot(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("Provides", testProvides)
	suite.Run(t)
}
