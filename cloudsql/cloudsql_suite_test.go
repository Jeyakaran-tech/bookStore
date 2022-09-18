package cloudsql_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCloudsql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudsql Suite")
}
