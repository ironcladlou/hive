package hive

import (
	g "github.com/onsi/ginkgo/v2"
	compat_otp "github.com/openshift/origin/test/extended/util/compat_otp"
	e2e "k8s.io/kubernetes/test/e2e/framework"
)

var _ = g.BeforeSuite(func() {
	if err := compat_otp.InitTest(false); err != nil {
		e2e.Failf("Failed to initialize test framework: %v", err)
	}
	e2e.AfterReadingAllFlags(compat_otp.TestContext)
})
