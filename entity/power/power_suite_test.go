package power_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPower(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Power Suite")
}
