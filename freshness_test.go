package vizzini_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Freshness", func() {
	Describe("Creating a fresh domain", func() {
		Context("with no TTL", func() {
			It("should create a fresh domain that never disappears", func() {
				Ω(bbsClient.UpsertDomain(domain, 0)).Should(Succeed())
				Consistently(bbsClient.Domains, 3).Should(ContainElement(domain))
				bbsClient.UpsertDomain(domain, 1*time.Second) //to clear it out
			})
		})

		Context("with a TTL", func() {
			It("should create a fresh domain that eventually disappears", func() {
				Ω(bbsClient.UpsertDomain(domain, 2*time.Second)).Should(Succeed())

				Ω(bbsClient.Domains()).Should(ContainElement(domain))
				Eventually(bbsClient.Domains, 5).ShouldNot(ContainElement(domain))
			})
		})

		Context("with no domain", func() {
			It("should error", func() {
				Ω(bbsClient.UpsertDomain("", 0)).ShouldNot(Succeed())
			})
		})
	})
})
