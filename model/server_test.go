package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"luckgo/model"
	"time"
)

var _ = Describe("Model", func() {
	BeforeEach(func() {
	})

	Describe("server", func() {
		Context("Check start and close", func() {
			It("Should be successful", func() {
				srv := model.NewServer()
				go func() {
					time.Sleep(6*time.Second)
					srv.Server.Close()
				}()
				srv.SqlSupplier = model.NewSqlSupplier()
				err := srv.Start()
				Expect(err.Error()).To(Equal("http: Server closed"))
			})
		})
	})
})