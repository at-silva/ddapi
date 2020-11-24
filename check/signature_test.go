package check_test

import (
	"encoding/hex"
	"time"

	"github.com/at-silva/ddapi/check"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SignatureChecker", func() {

	Describe("Sha256HMAC", func() {
		var (
			checker check.SignatureChecker
			q       []byte
			s       []byte
			err     error
		)

		BeforeEach(func() {
			key := []byte("my secret key")
			checker = check.Sha256HMAC(key)
			q = []byte("select name, id from users where id = :id")
			s, err = hex.DecodeString("e5ae3948a3345025bb29b70fd4df2a178aa3f528a50d3cfffe569d80968adfff")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be idempotent", func() {
			for i := 0; i < 10; i++ {
				Expect(checker.Check(q, s)).Should(Succeed())
				time.Sleep(100 * time.Millisecond)
			}
		})

		It("should fail if the query is empty", func() {
			Expect(checker.Check([]byte{}, s)).ShouldNot(Succeed())
		})

		It("should fail if the signature is empty", func() {
			Expect(checker.Check(q, []byte{})).ShouldNot(Succeed())
		})

		It("should fail if the signature is invalid", func() {
			s = []byte("select name, id from users where id = ?id")
			Expect(checker.Check(q, s)).ShouldNot(Succeed())
		})

		It("should fail if the query is invalid", func() {
			s = []byte("select name, id from users where id = ?id")
			s, err = hex.DecodeString("ae5ae3948a3345025bb29b70fd4df2a178aa3f528a50d3cfffe569d80968adff")
			Expect(err)
			Expect(checker.Check(q, s)).ShouldNot(Succeed())
		})

		Measure("it should run in less than half a millisecond (500 microseconds)", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				Expect(checker.Check(q, s)).Should(Succeed())
			})

			Expect(runtime.Microseconds()).Should(BeNumerically("<", 500))
			b.RecordValue("runtime (in ms)", float64(runtime.Microseconds()))
		}, 100)

	})

})
