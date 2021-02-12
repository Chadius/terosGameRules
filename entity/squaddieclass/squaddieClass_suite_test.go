package squaddieclass_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSquaddieClass(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SquaddieClass Suite")
}
