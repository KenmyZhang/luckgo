package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"luckgo/model"
)

var _ = Describe("Model", func() {
	BeforeEach(func() {
	})

	Describe("sqlSupplier", func() {
		Context("Check start and close", func() {
			It("Total MasterDbConnections", func() {
				count := model.Srv.SqlSupplier.TotalMasterDbConnections()
				Expect(count).To(Equal(1))
			})
			It("Total ReadDbConnections", func() {
			count := model.Srv.SqlSupplier.TotalReadDbConnections()
				Expect(count).To(Equal(0))
			})
			It("Total SearchDbConnections", func() {
			count := model.Srv.SqlSupplier.TotalSearchDbConnections()
				Expect(count).To(Equal(0))
			})
		})
	})
})