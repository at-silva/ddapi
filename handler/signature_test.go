package handler

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/at-silva/ddapi/check/checkfakes"
	"github.com/at-silva/ddapi/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckSignatures", func() {

	var (
		fakeNext             *handlerfakes.FakeHandler
		fakeSignatureChecker *checkfakes.FakeSignatureChecker
		recorder             *httptest.ResponseRecorder
		ehandler             http.Handler
	)

	BeforeEach(func() {
		fakeSignatureChecker = new(checkfakes.FakeSignatureChecker)
		fakeNext = new(handlerfakes.FakeHandler)
		recorder = httptest.NewRecorder()
		ehandler = CheckSignatures(fakeSignatureChecker, fakeNext)
	})

	It("should call the next handler when the validation succeeds", func() {
		req := request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          base64.StdEncoding.EncodeToString([]byte("valid-sql-signature")),
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: base64.StdEncoding.EncodeToString([]byte("valid-params-schema-signature")),
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusOK))

		s, ss := fakeSignatureChecker.CheckArgsForCall(0)
		Expect(string(s)).Should(Equal(req.SQL))
		Expect(string(ss)).Should(Equal("valid-sql-signature"))

		s, ss = fakeSignatureChecker.CheckArgsForCall(1)
		Expect(string(s)).Should(Equal(req.ParamsSchema))
		Expect(string(ss)).Should(Equal("valid-params-schema-signature"))
		Expect(fakeNext.ServeHTTPCallCount()).Should(Equal(1))

	})

	It("should return InternalServerError when it can't find a request in the context", func() {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not check signatures: invalid request" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return BadRequest when the sql signature decoding fails", func() {
		req := request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          "invalid-signature",
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: base64.StdEncoding.EncodeToString([]byte("valid-params-schema-signature")),
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not decode sql signature: illegal base64 data at input byte 7" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return BadRequest when the sql signature decoding fails", func() {
		req := request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          base64.StdEncoding.EncodeToString([]byte("valid-sql-signature")),
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: "invalid-signature",
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not decode params schema signature: illegal base64 data at input byte 7" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return Forbidden when the sql validation fails", func() {
		req := request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          base64.StdEncoding.EncodeToString([]byte("valid-sql-signature")),
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: base64.StdEncoding.EncodeToString([]byte("valid-params-schema-signature")),
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeSignatureChecker.CheckReturns(errors.New("invalid signature"))

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusForbidden))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not validate sql signature: invalid signature" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return Forbidden when the params schema validation fails", func() {
		req := request{
			SQL:                   "insert into product(name) values(:name)",
			SQLSignature:          base64.StdEncoding.EncodeToString([]byte("valid-sql-signature")),
			ParamsSchema:          `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`,
			ParamsSchemaSignature: base64.StdEncoding.EncodeToString([]byte("valid-params-schema-signature")),
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeSignatureChecker.CheckReturnsOnCall(1, errors.New("invalid signature"))

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusForbidden))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not validate params schema signature: invalid signature" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})
})
