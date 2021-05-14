package handlers

import (
	"os"
	"testing"
)

//TestCreateTokenReturnTrueIfFuncReturnToken
//return error if return error
func TestCreateTokenReturnErrorIfFuncReturnError(t *testing.T) {

	if _, err := CreateToken("sdfsdfsd"); err != nil {
		t.Error("incorrect")
	}
}

//TestCreateTokenReturnErrorIfFuncReturnToken
//return error if return Token
func TestCreateTokenReturnErrorIfFuncReturnToken(t *testing.T) {
	os.Setenv("ACCESS_SECRET", "")
	if _, err := CreateToken("fgsfg"); err == nil {
		t.Error("incorrect")
	}
}
