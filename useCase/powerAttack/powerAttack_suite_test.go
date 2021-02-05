package powerAttack_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPowerAttack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PowerAttack Suite")
}
