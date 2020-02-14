package cluster_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/testutils"
	cli_test "github.com/solo-io/mesh-projects/cli/pkg/test"
	cli_util "github.com/solo-io/mesh-projects/cli/pkg/util"
)

var _ = Describe("Cluster Root Cmd", func() {
	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("complains if it is invoked without a subcommand", func() {
		output, err := cli_test.MockMeshctl{MockController: ctrl, Ctx: context.TODO()}.Invoke("cluster --kubeconfig foo")
		Expect(output).To(BeEmpty())

		nonTerminalCommandErrorBuilder := cli_util.NonTerminalCommand("cluster")
		nonTerminalErr := nonTerminalCommandErrorBuilder(nil, nil)
		Expect(err).To(testutils.HaveInErrorChain(nonTerminalErr))
	})
})
