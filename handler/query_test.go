package handler

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/at-silva/ddapi/db/dbfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("queryHandler", func() {

	var (
		fakeDB   *dbfakes.FakeDB
		fakeRows *dbfakes.FakeRows

		recorder *httptest.ResponseRecorder
		ehandler http.Handler
	)

	BeforeEach(func() {
		fakeDB = new(dbfakes.FakeDB)
		fakeRows = new(dbfakes.FakeRows)

		recorder = httptest.NewRecorder()
		ehandler = queryHandler{fakeDB}
	})

	It("should return Ok when a database call succeeds", func() {
		req := request{SQL: "select * from product where name = :name"}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeRows.NextReturnsOnCall(0, true)
		fakeRows.MapScanStub = func(m map[string]interface{}) error {
			m["name"] = "Product 1"
			return nil
		}
		fakeDB.NamedQueryContextReturns(fakeRows, nil)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Body).Should(MatchJSON(`{
			"data": [{"name":"Product 1"}],
			"error": null
		}`))
		Expect(recorder.Code).Should(Equal(http.StatusOK))
		_, sql, p := fakeDB.NamedQueryContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

	It("should cast []bytes to string in the result", func() {
		req := request{SQL: "select * from product where name = :name"}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeRows.NextReturnsOnCall(0, true)
		fakeRows.MapScanStub = func(m map[string]interface{}) error {
			m["name"] = []byte("Product 1")
			return nil
		}
		fakeDB.NamedQueryContextReturns(fakeRows, nil)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Body).Should(MatchJSON(`{
			"data": [{"name":"Product 1"}],
			"error": null
		}`))
		Expect(recorder.Code).Should(Equal(http.StatusOK))
		_, sql, p := fakeDB.NamedQueryContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

	It("should return InternalServerErrror when it can't find a request in the context", func() {
		params := map[string]interface{}{"name": "Product1"}
		ctx := context.WithValue(context.Background(), DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not query the database: invalid request" 
		}`))
		Expect(fakeDB.NamedExecContextCallCount()).Should(BeZero())
	})

	It("should return InternalServerErrror when it can't find the params in the context", func() {
		req := request{
			SQL:    "insert into product(name) values(:name)",
			Params: `{"name": "Product 1"}`,
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)
		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not query the database: invalid params" 
		}`))
		Expect(fakeDB.NamedExecContextCallCount()).Should(BeZero())
	})

	It("should return InternalServerErrror when a database call fails", func() {
		req := request{SQL: "select * from product where name = :name"}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		fakeDB.NamedQueryContextReturns(nil, sql.ErrNoRows)
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)

		Expect(err).ShouldNot(HaveOccurred())

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"error":"could not query the database: sql: no rows in result set" 
		}`))
		_, sql, p := fakeDB.NamedQueryContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

})
