// +build unit_test

package cryptoutils

import "testing"

var testCases = []string{
	"a", "a ", "my text", "my text contains ελληνικά!",
	"this is a really long text : ad2os4dn3af1weo5fx39fn347fx935gxb3f736bf398dxv2389zd263vd238d6v2d",
	"+_)(*&^%$#@!@#$%^&*()&*(^7453$5#4^&*%(70^(&5645Y5GTG#EfEδ δωδφωδφωφερΒΡΡρβέφωέώε¨΄΄τφρφς΅΅δωφ",
}

func TestArgon2Hash(t *testing.T) {
	seenHashes := make(map[string]struct{})

	for _, tc := range testCases {
		h, err := Argon2Hash(tc)

		if err != nil {
			t.Error(err)
		}

		if _, exists := seenHashes[h]; exists {
			t.Errorf("collision detected! Hash(%s)=%s", tc, h)
		}

		if len(h) != 97 {
			t.Errorf("length %d is not valid: %s", len(h), h)
		}

		seenHashes[h] = struct{}{}
	}
}

func BenchmarkArgon2Hash(b *testing.B) {
	p := "thisIsMyRal3Saf5%Pa$w_ord!"

	for i := 0; i < b.N; i++ {
		_, err := Argon2Hash(p)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestArgon2Verify(t *testing.T) {
	for _, tc := range testCases {
		h, err := Argon2Hash(tc)

		if err != nil {
			t.Error(err)
		}

		if err = Argon2Verify(tc, h); err != nil {
			t.Error(err)
		}

		if err = Argon2Verify(tc, ""); err == nil {
			t.Errorf("Expected an error in verification of tc='%s' with nil string", tc)
		}

		tc2 := []byte(tc)
		tc2[0] = ' '
		if err = Argon2Verify(string(tc2), h); err == nil {
			t.Errorf("Expected an error in verification of tc='%s' with altered string", tc)
		}
	}
}
