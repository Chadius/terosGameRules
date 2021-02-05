package squaddie_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSquaddie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Squaddie Suite")
}
