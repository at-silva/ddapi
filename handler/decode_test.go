package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/at-silva/ddapi/handler/handlerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeRequest", func() {

	var (
		fakeNext *handlerfakes.FakeHandler
		recorder *httptest.ResponseRecorder
		dhandler http.Handler
	)

	BeforeEach(func() {
		fakeNext = new(handlerfakes.FakeHandler)
		recorder = httptest.NewRecorder()
		dhandler = DecodeRequest(fakeNext)
	})

	It("should decode a request into the context", func() {
		body := `
		{
			"sql": "insert into product(name) values(:name)",
			"sqlSignature": "valid-sql-signature",
			"params": "{\"name\": \"Product 1\"}",
			"paramsSchema": "{\"type\":\"object\", \"required\": [\"name\"], \"properties\": {\"name\": {\"type\": \"string\"}}}",
			"paramsSchemaSignature": "valid-params-signature"
		}`

		ctx := context.Background()

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", strings.NewReader(body))
		Expect(err).ShouldNot(HaveOccurred())

		dhandler.ServeHTTP(recorder, r)
		_, req := fakeNext.ServeHTTPArgsForCall(0)

		Expect(req.Context().Value(DecodedRequest)).Should(Equal(request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          "valid-sql-signature",
			Params:                `{"name": "Product 1"}`,
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: "valid-params-signature",
		}))

		Expect(req.Context().Value(DecodedParams)).Should(Equal(map[string]interface{}{
			"name": "Product 1",
		}))
	})

	It("should return InternalServerError when it can't read the body", func() {
		ctx := context.Background()
		fakeReader := new(handlerfakes.FakeReader)
		fakeReader.ReadReturns(0, io.ErrUnexpectedEOF)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", fakeReader)
		Expect(err).ShouldNot(HaveOccurred())

		dhandler.ServeHTTP(recorder, r)
		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{"error":"could not read body: unexpected EOF"}`))
	})

	It("should return BadRequest when it can't unmarshal the request", func() {
		body := ""

		ctx := context.Background()

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", strings.NewReader(body))
		Expect(err).ShouldNot(HaveOccurred())

		dhandler.ServeHTTP(recorder, r)
		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{"error":"could not unmarshal body: unexpected end of JSON input"}`))
	})

	It("should return BadRequest when it can't unmarshal the params", func() {
		body := `
		{
			"sql": "insert into product(name) values(:name)",
			"sqlSignature": "valid-sql-signature",
			"params": "",
			"paramsSchema": "{\"type\":\"object\", \"required\": [\"name\"], \"properties\": {\"name\": {\"type\": \"string\"}}}",
			"paramsSchemaSignature": "valid-params-signature"
		}`

		ctx := context.Background()

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", strings.NewReader(body))
		Expect(err).ShouldNot(HaveOccurred())

		dhandler.ServeHTTP(recorder, r)
		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{"error":"could not unmarshal params: unexpected end of JSON input"}`))
	})

})
