// +build unit_test

package terror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorStringWithCode(t *testing.T) {
	testCases := []errorStringWithCode{
		{
			msg:  "asd",
			code: 123,
		},
		{
			msg:  "a",
			code: 1,
		},
		{
			msg:  "asd ασδ !@#",
			code: 123987,
		},
	}

	for _, tc := range testCases {
		terr := New(tc.code, tc.msg)
		assert.Equal(t, tc.msg, terr.Error())
		assert.Equal(t, tc.code, terr.Code())
	}
}

func ExampleErrorStringWithCode() {
	f := func() error {
		err := New(404, "not found")
		return err
	}

	err := f()
	terr, isTerr := err.(Terror)

	fmt.Println(isTerr)
	fmt.Println(terr.Code(), terr.Error())

	// Output:
	// true
	// 404 not found
}
