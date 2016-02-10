package blackbox

import (
	"os/exec"
	"strings"

	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/say"
)

func CF(dir string, args ...string) *gexec.Session {
	say.Println(0, say.Green("cf %s", strings.Join(args, " ")))
	cf := exec.Command("cf", args...)
	cf.Dir = dir
	session, err := gexec.Start(cf, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
