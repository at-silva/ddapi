package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/at-silva/ddapi/handler"
	"github.com/at-silva/ddapi/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	var (
		fakeNext *handlerfakes.FakeHandler
		recorder *httptest.ResponseRecorder
		vhandler http.Handler
	)

	BeforeEach(func() {
		fakeNext = new(handlerfakes.FakeHandler)
		recorder = httptest.NewRecorder()
		vhandler = handler.ValidateFormRequest(fakeNext)
	})

	DescribeTable("ValidateFormRequest", func(missing, sql, sqlSignature, params, paramsSchema, paramsSchemaSignature string) {
		r, err := http.NewRequest(http.MethodPost, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		r.PostForm = url.Values{
			"sql":                   {sql},
			"sqlSignature":          {sqlSignature},
			"params":                {params},
			"paramsSchema":          {paramsSchema},
			"paramsSchemaSignature": {paramsSchemaSignature},
		}

		vhandler.ServeHTTP(recorder, r)
		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(fmt.Sprintf(`{"error":"invalid request: %s is required"}`, missing)))

	},
		Entry("sql is required", "sql", "", "sqlSignature", "params", "paramsSchema", "paramsSchemaSignature"),
		Entry("sqlSignature is required", "sqlSignature", "sql", "", "params", "paramsSchema", "paramsSchemaSignature"),
		Entry("params is required", "params", "sql", "sqlSignature", "", "paramsSchema", "paramsSchemaSignature"),
		Entry("paramsSchema is required", "paramsSchema", "sql", "sqlSignature", "params", "", "paramsSchemaSignature"),
		Entry("paramsSchemaSignature is required", "paramsSchemaSignature", "sql", "sqlSignature", "params", "paramsSchema", ""),
	)

	Describe("ValidateFormRequest", func() {

		It("should call the next handler", func() {
			r, err := http.NewRequest(http.MethodPost, "/exec", nil)
			Expect(err).ShouldNot(HaveOccurred())

			r.PostForm = url.Values{
				"sql":                   {"insert into product(name) values(:name)"},
				"sqlSignature":          {"valid-sql-signature"},
				"params":                {`{"name": "Product 1"}`},
				"paramsSchema":          {`{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`},
				"paramsSchemaSignature": {"valid-params-signature"},
			}

			vhandler.ServeHTTP(recorder, r)
			_, req := fakeNext.ServeHTTPArgsForCall(0)
			Expect(r).Should(Equal(req))
		})

	})
})
