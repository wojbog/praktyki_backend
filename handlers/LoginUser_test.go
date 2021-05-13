package handlers

import (
	"testing"
)

//TestCreateTokenReturnTrueIfFuncReturnToken test func CreateTokenR
func TestCreateTokenReturnTrueIfFuncReturnToken(t *testing.T) {

	if _, err := CreateToken("sdfsdfsd"); err != nil {
		t.Error("incorrect")
	}
}