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

var _ = Describe("execHandler", func() {

	var (
		fakeDB     *dbfakes.FakeDB
		fakeResult *dbfakes.FakeResult

		recorder *httptest.ResponseRecorder
		ehandler http.Handler
	)

	BeforeEach(func() {
		fakeDB = new(dbfakes.FakeDB)
		fakeResult = new(dbfakes.FakeResult)

		recorder = httptest.NewRecorder()
		ehandler = execHandler{fakeDB}
	})

	It("should return Ok when a database call succeeds", func() {
		req := request{
			SQL:    "insert into product(name) values(:name)",
			Params: `{"name": "Product 1"}`,
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)

		Expect(err).ShouldNot(HaveOccurred())

		fakeResult.LastInsertIdReturns(1, nil)
		fakeResult.RowsAffectedReturns(1, nil)
		fakeDB.NamedExecContextReturns(fakeResult, nil)

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusOK))
		Expect(recorder.Body).Should(MatchJSON(`{
			"rowsAffected": 1,
			"lastInsertedId": 1,
			"error": null
		}`))
		_, sql, p := fakeDB.NamedExecContextArgsForCall(0)
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
			"rowsAffected": 0,
			"lastInsertedId": 0,
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
			"rowsAffected": 0,
			"lastInsertedId": 0,
			"error":"could not query the database: invalid params" 
		}`))
		Expect(fakeDB.NamedExecContextCallCount()).Should(BeZero())
	})

	It("should return InternalServerErrror when a database call fails", func() {
		req := request{
			SQL:    "insert into product(name) values(:name)",
			Params: `{"name": "Product 1"}`,
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)

		Expect(err).ShouldNot(HaveOccurred())

		fakeResult.LastInsertIdReturns(1, nil)
		fakeResult.RowsAffectedReturns(1, nil)
		fakeDB.NamedExecContextReturns(fakeResult, sql.ErrNoRows)

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"rowsAffected": 0,
			"lastInsertedId": 0,
			"error":"could not query the database: sql: no rows in result set" 
		}`))
		_, sql, p := fakeDB.NamedExecContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

	It("should return InternalServerErrror RowsAffected call fails", func() {
		req := request{
			SQL:    "insert into product(name) values(:name)",
			Params: `{"name": "Product 1"}`,
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)

		Expect(err).ShouldNot(HaveOccurred())

		fakeResult.LastInsertIdReturns(1, nil)
		fakeResult.RowsAffectedReturns(1, sql.ErrConnDone)
		fakeDB.NamedExecContextReturns(fakeResult, nil)

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"rowsAffected": 0,
			"lastInsertedId": 0,
			"error":"could not read the number of rows affected: sql: connection is already closed" 
		}`))
		_, sql, p := fakeDB.NamedExecContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

	It("should return InternalServerErrror LastInsertId call fails", func() {
		req := request{
			SQL:    "insert into product(name) values(:name)",
			Params: `{"name": "Product 1"}`,
		}
		ctx := context.WithValue(context.Background(), DecodedRequest, req)

		params := map[string]interface{}{"name": "Product1"}
		ctx = context.WithValue(ctx, DecodedParams, params)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/exec", nil)

		Expect(err).ShouldNot(HaveOccurred())

		fakeResult.LastInsertIdReturns(1, sql.ErrConnDone)
		fakeResult.RowsAffectedReturns(1, nil)
		fakeDB.NamedExecContextReturns(fakeResult, nil)

		ehandler.ServeHTTP(recorder, request)

		Expect(recorder.Code).Should(Equal(http.StatusInternalServerError))
		Expect(recorder.Body).Should(MatchJSON(`{
			"rowsAffected": 0,
			"lastInsertedId": 0,
			"error":"could not read the last inserted id: sql: connection is already closed" 
		}`))
		_, sql, p := fakeDB.NamedExecContextArgsForCall(0)
		Expect(sql).Should(Equal(req.SQL))
		Expect(p).Should(Equal(params))
	})

})
