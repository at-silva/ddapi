package session_test

import (
	"time"

	"github.com/at-silva/ddapi/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {

	var (
		token  string
		reader session.Reader
		params map[string]interface{}
	)

	Describe("HS256JWT", func() {

		BeforeEach(func() {
			params = map[string]interface{}{}
			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.G0X_hg_UT05tEzRrZqW9EEfDMgeBqlxqBHbsmT9liYs"
			reader = session.Reader(session.HS256JWT([]byte("my_jwt_secret")))
		})

		It("should be idempotent", func() {
			for i := 0; i < 10; i++ {
				Expect(reader.Copy(token, params)).ShouldNot(HaveOccurred())
				Expect(params).Should(HaveKeyWithValue("iat", BeEquivalentTo(1516239022)))
				Expect(params).Should(HaveKeyWithValue("name", "John Doe"))
				Expect(params).Should(HaveKeyWithValue("sub", "1234567890"))
				time.Sleep(100 * time.Millisecond)
			}
		})

		It("should copy the claims of the given JWT into the given map", func() {
			Expect(reader.Copy(token, params)).ShouldNot(HaveOccurred())
			Expect(params).Should(HaveKeyWithValue("iat", BeEquivalentTo(1516239022)))
			Expect(params).Should(HaveKeyWithValue("name", "John Doe"))
			Expect(params).Should(HaveKeyWithValue("sub", "1234567890"))
		})

		It("should fail if an empty token gets passed in", func() {
			Expect(reader.Copy("", map[string]interface{}{})).ShouldNot(Succeed())
		})

		It("should fail if an nil map gets passed in", func() {
			Expect(reader.Copy(token, nil)).ShouldNot(Succeed())
		})

		It("should override preexisting keys in the given map", func() {
			Expect(reader.Copy(token, params)).ShouldNot(HaveOccurred())
			Expect(params).Should(HaveKeyWithValue("iat", BeEquivalentTo(1516239022)))
			Expect(params).Should(HaveKeyWithValue("name", "John Doe"))
			Expect(params).Should(HaveKeyWithValue("sub", "1234567890"))
			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphbmUgRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.6w99P5yVAaRf5yIubb8EsnX8rkOOd6zNuSQ6HetsFDo"
			Expect(reader.Copy(token, params)).ShouldNot(HaveOccurred())
			Expect(params).Should(HaveKeyWithValue("iat", BeEquivalentTo(1516239022)))
			Expect(params).Should(HaveKeyWithValue("name", "Jane Doe"))
			Expect(params).Should(HaveKeyWithValue("sub", "1234567890"))
		})

		It("should fail if the algorithm is not supported", func() {
			token = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphbmUgRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.6w99P5yVAaRf5yIubb8EsnX8rkOOd6zNuSQ6HetsFDo"
			Expect(reader.Copy(token, params)).ShouldNot(Succeed())
		})

		It("should fail if the token is invalid", func() {
			token = "invalid_token"
			Expect(reader.Copy(token, params)).ShouldNot(Succeed())
		})

		It("should fail if the signature is invalid", func() {
			reader = session.Reader(session.HS256JWT([]byte("invalid_secret")))
			Expect(reader.Copy(token, params)).ShouldNot(Succeed())
		})
	})

})
