package vizzini_test

import (
	"github.com/cloudfoundry-incubator/bbs/models"
	. "github.com/cloudfoundry-incubator/vizzini/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Privileged", func() {
	var task *models.TaskDefinition
	var runUser string
	var containerPrivileged bool

	JustBeforeEach(func() {
		task = Task()
		task.Privileged = containerPrivileged
		task.Action = models.WrapAction(&models.RunAction{
			Path: "bash",
			Args: []string{"-c", "id > /tmp/bar; echo h > /proc/sysrq-trigger ; echo have_real_root=$? >> /tmp/bar"},
			User: runUser,
		})

		Expect(bbsClient.DesireTask(guid, domain, task)).To(Succeed())
		Eventually(TaskGetter(guid)).Should(HaveTaskState(models.Task_Completed))
	})

	AfterEach(func() {
		Expect(bbsClient.ResolvingTask(guid)).To(Succeed())
		Expect(bbsClient.DeleteTask(guid)).To(Succeed())
	})

	Context("with a privileged container", func() {
		BeforeEach(func() {
			containerPrivileged = true
		})

		Context("when running a privileged action", func() {
			BeforeEach(func() {
				runUser = "root"
			})

			It("should run as root", func() {
				completedTask, err := bbsClient.TaskByGuid(guid)
				Expect(err).NotTo(HaveOccurred())
				Expect(completedTask.Result).To(ContainSubstring("uid=0(root)"), "If this fails, then your executor may not be configured to allow privileged actions")
				Expect(completedTask.Result).To(MatchRegexp(`groups=.*0\(root\)`))
				Expect(completedTask.Result).To(ContainSubstring("have_real_root=0"))
			})
		})

		Context("when running a non-privileged action", func() {
			BeforeEach(func() {
				runUser = "vcap"
			})

			It("should run as non-root", func() {
				completedTask, err := bbsClient.TaskByGuid(guid)
				Expect(err).NotTo(HaveOccurred())
				Expect(completedTask.Result).To(MatchRegexp(`uid=\d{4,5}\(vcap\)`))
				Expect(completedTask.Result).To(MatchRegexp(`groups=.*0\(root\)`))
				Expect(completedTask.Result).To(ContainSubstring("have_real_root=1"))
			})
		})
	})

	Context("with an unprivileged container", func() {
		BeforeEach(func() {
			containerPrivileged = false
		})

		Context("when running a privileged action", func() {
			BeforeEach(func() {
				runUser = "root"
			})

			It("should run as namespaced root", func() {
				completedTask, err := bbsClient.TaskByGuid(guid)
				Expect(err).NotTo(HaveOccurred())
				Expect(completedTask.Result).To(ContainSubstring("uid=0(root)"), "If this fails, then your executor may not be configured to allow privileged actions")
				Expect(completedTask.Result).To(ContainSubstring("have_real_root=1"))
			})
		})

		Context("when running a non-privileged action", func() {
			BeforeEach(func() {
				runUser = "vcap"
			})

			It("should run as non-root", func() {
				completedTask, err := bbsClient.TaskByGuid(guid)
				Expect(err).NotTo(HaveOccurred())
				Expect(completedTask.Result).To(MatchRegexp(`uid=\d{4,5}\(vcap\)`))
				Expect(completedTask.Result).To(ContainSubstring("have_real_root=1"))
			})
		})
	})
})
