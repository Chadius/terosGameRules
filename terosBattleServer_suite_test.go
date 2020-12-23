package terosbattleserver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTerosBattleServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerosBattleServer Suite")
}
