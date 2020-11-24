package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/at-silva/ddapi/session/sessionfakes"

	"github.com/at-silva/ddapi/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {

	var (
		fakeNext          *handlerfakes.FakeHandler
		fakeSessionReader *sessionfakes.FakeReader
		recorder          *httptest.ResponseRecorder
		ehandler          http.Handler
	)

	BeforeEach(func() {
		fakeSessionReader = new(sessionfakes.FakeReader)
		fakeNext = new(handlerfakes.FakeHandler)
		recorder = httptest.NewRecorder()
		ehandler = ReadSession(fakeSessionReader, fakeNext)
	})

	It("should call the next handler when the session is read successfully", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		request.Header.Set("Authorization", "Bearer valid-jwt")

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusOK))
		t, pm := fakeSessionReader.CopyArgsForCall(0)
		Expect(t).Should(Equal("valid-jwt"))
		Expect(pm).Should(Equal(params))
		Expect(fakeNext.ServeHTTPCallCount()).Should(Equal(1))

	})

	It("should return InternalServerError when it can't find the params in the context", func() {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not copy session params: invalid params" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return BadRequest when it can't find the Authorization header", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not copy session params: invalid Authorization header" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return BadRequest when the Authorization header is invalid", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		request.Header.Set("Authorization", "invalid-header")

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not copy session params: invalid Authorization header" 
		}`))
		Expect(fakeNext.ServeHTTPCallCount()).Should(BeZero())
	})

	It("should return InternalServerError when the session reader fails", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		request.Header.Set("Authorization", "Bearer valid-jwt")
		fakeSessionReader.CopyReturns(errors.New("invalid jwt"))

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not copy session params: invalid jwt" 
		}`))

	})

})
