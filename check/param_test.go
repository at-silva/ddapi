package check_test

import (
	"time"

	"github.com/at-silva/ddapi/check"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParamsChecker", func() {

	Describe("JSONSchema", func() {

		var (
			pm map[string]interface{}
			s  string
		)

		BeforeEach(func() {
			s = `
{
	"type":"object", 
	"required": ["name"], 
	"properties": {
		"id": {
			"type": "integer"
		},
		"name": {
			"type": "string"
		}
	}
}`
			pm = map[string]interface{}{
				"id":   42,
				"name": "product 1",
			}
		})

		It("should be idempotent", func() {
			for i := 0; i < 10; i++ {
				Expect(check.JSONSchema(pm, s)).Should(Succeed())
				time.Sleep(100 * time.Millisecond)
			}
		})

		It("should return an error if schema is empty", func() {
			Expect(check.JSONSchema(pm, "")).ShouldNot(Succeed())
		})

		It("should return an error if params is empty", func() {
			Expect(check.JSONSchema(nil, s)).ShouldNot(Succeed())
		})

		It("should return an error if the parameters are invalid", func() {
			pm["id"] = ""
			Expect(check.JSONSchema(pm, s)).ShouldNot(Succeed())
		})

		It("should return an error if the schema is invalid", func() {
			s = "invalid schema"
			Expect(check.JSONSchema(pm, s)).ShouldNot(Succeed())
		})

		Measure("it should run in less than half a millisecond (500 microseconds)", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				Expect(check.JSONSchema(pm, s)).Should(Succeed())
			})

			Expect(runtime.Microseconds()).Should(BeNumerically("<", 500))
			b.RecordValue("runtime (in ms)", float64(runtime.Microseconds()))
		}, 100)

	})

})
