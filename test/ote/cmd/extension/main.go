package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/openshift-eng/openshift-tests-extension/pkg/cmd"
	e "github.com/openshift-eng/openshift-tests-extension/pkg/extension"
	et "github.com/openshift-eng/openshift-tests-extension/pkg/extension/extensiontests"
	g "github.com/openshift-eng/openshift-tests-extension/pkg/ginkgo"

	// Import test packages to register Ginkgo tests
	_ "github.com/openshift/hive/test/ote/hive"
)

var platformFileSelectors = map[string]string{
	"hive_aws.go":     "aws",
	"hive_azure.go":   "azure",
	"hive_gcp.go":     "gcp",
	"hive_vsphere.go": "vsphere",
}

func main() {
	registry := e.NewRegistry()

	ext := e.NewExtension("openshift", "payload", "hive")

	ext.AddSuite(e.Suite{
		Name:    "openshift/hive",
		Parents: []string{"openshift/conformance/parallel"},
	})

	specs, err := g.BuildExtensionTestSpecsFromOpenShiftGinkgoSuite()
	if err != nil {
		panic(fmt.Sprintf("couldn't build extension test specs from ginkgo: %+v", err.Error()))
	}

	applyEnvironmentSelectors(specs)

	ext.AddSpecs(specs)
	registry.Register(ext)

	root := &cobra.Command{
		Long: "OpenShift Hive OTE Extension",
	}

	root.AddCommand(cmd.DefaultExtensionCommands(registry)...)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func applyEnvironmentSelectors(specs et.ExtensionTestSpecs) {
	specs.Walk(func(spec *et.ExtensionTestSpec) {
		for _, cl := range spec.CodeLocations {
			for file, platform := range platformFileSelectors {
				if strings.Contains(cl, file) {
					spec.Include(et.PlatformEquals(platform))
					return
				}
			}
		}
	})
}
