package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testProvides(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when the buildpack is run with pack build", func() {
		var (
			image  occam.Image
			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("Provides a process", func() {
			var err error

			source, err = occam.Source(filepath.Join("testdata", "fake_app"))
			Expect(err).ToNot(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					settings.Buildpacks.Procfile.Online,
					settings.Buildpacks.BuildPlan.Online,
				).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)

			Expect(logs).To(ContainLines(
				"Paketo Buildpack for Procfile 0.0.0",
				"  https://github.com/paketo-buildpacks/procfile",
				"  Process types:",
				ContainSubstring("procfile: echo Procfile Command"),
				ContainSubstring("web:      sleep 3600"),
			))
		})
	})
}
