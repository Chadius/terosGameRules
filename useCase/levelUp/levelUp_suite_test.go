package levelup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLevelUp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LevelUp Suite")
}
