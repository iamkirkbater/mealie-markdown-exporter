package apitoken_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApitoken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apitoken Suite")
}
