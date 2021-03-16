package powerusage_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPowerUsage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PowerUsage Suite")
}
