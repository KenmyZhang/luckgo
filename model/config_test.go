package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"luckgo/model"
)

var _ = Describe("Model", func() {
	BeforeEach(func() {
	})

	Describe("config", func() {
		Context("Check config default", func() {
			It("Should be successful", func() {
				cfg := &model.Config{}
				cfg.SetDefaults()
				Expect(*cfg.SqlSettings.DriverName).To(Equal(model.DATABASE_DRIVER_MYSQL))
				Expect(*cfg.SqlSettings.DataSource).To(Equal(model.SQL_SETTINGS_DEFAULT_DATA_SOURCE))
				Expect(*cfg.SqlSettings.MaxIdleConns).To(Equal(20))
				Expect(*cfg.SqlSettings.MaxOpenConns).To(Equal(300))
				Expect(*cfg.SqlSettings.QueryTimeout).To(Equal(30))
				Expect(cfg.ServiceSettings.ListenAddress).To(Equal(model.SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS))
				Expect(*cfg.LogSettings.ConsoleJson).To(Equal(true))
				Expect(*cfg.LogSettings.FileJson).To(Equal(true))
			})
		})
	})
})