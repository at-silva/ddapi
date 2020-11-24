package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/at-silva/ddapi/check/checkfakes"
	"github.com/at-silva/ddapi/handler/handlerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckParams", func() {

	var (
		fakeNext          *handlerfakes.FakeHandler
		fakeParamsChecker *checkfakes.FakeParamsChecker
		recorder          *httptest.ResponseRecorder
		ehandler          http.Handler
	)

	BeforeEach(func() {
		fakeParamsChecker = new(checkfakes.FakeParamsChecker)
		fakeNext = new(handlerfakes.FakeHandler)
		recorder = httptest.NewRecorder()
		ehandler = CheckParams(fakeParamsChecker, fakeNext)
	})

	It("should call the next handler when the validation succeeds", func() {
		req := request{ParamsSchema: `{"type":"object", "required": ["name"], "properties": {"name": {"type": "string"}}}`}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusOK))
		Expect(fakeNext.ServeHTTPCallCount()).Should(Equal(1))

		p, ps := fakeParamsChecker.CheckArgsForCall(0)
		Expect(p).Should(Equal(params))
		Expect(ps).Should(Equal(req.ParamsSchema))
	})

	It("should return InternalServerErrror when it can't find a request in the context", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not check params: invalid request" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return InternalServerError when it can't find the params in the context", func() {
		req := request{Params: `{"name": "Product 1"}`}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not check params: invalid params" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return Forbidden when the validation fails", func() {
		req := request{Params: `{"name": "Product 1"}`}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeParamsChecker.CheckReturns(errors.New("id is required"))

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusForbidden))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"invalid params: id is required" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})
})
