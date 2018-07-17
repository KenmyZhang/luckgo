package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"luckgo/model"
	"fmt"
	"os"
)

var _ = Describe("Model", func() {
	BeforeEach(func() {
	})

	Describe("utils", func() {
		Context("Check Error", func() {
			It("Should be successful", func() {
				err := model.Err{Code:200, Result:"success"}
				Expect(err.Error()).To(Equal("code:200, result:success"))
			})

			It("Should be successful", func() {
				var err model.Err
				str := err.Error()
				Expect(str).To(Equal("code:0, result:"))
			})
		})

		Context("Check NewInternalServerError", func() {
			It("Should be successful", func() {
				err := model.NewInternalServerError("here", "")
				fmt.Println(err.Error())
				Expect(err.Error()).To(Equal("code:5000, result:Internal Server Error,where:here"))
			})
		})

		Context("Check FindFile", func() {
			It("Should be successful", func() {
				path := model.FindFile("")
				Expect(path).To(Equal(""))
			})
		})

		Context("Check GetFile", func() {
			It("Should be successful", func() {
				str := "./test/sentry.log"
				file, err := model.GetFile(str)
				if err != nil {
					Fail(err.Error())
				}

				os.RemoveAll(str)
				Expect(file.Name()).To(Equal(str))
			})
		})

	})
})