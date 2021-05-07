package person

import (
	"testing"
	"github.com/go-playground/validator"
	
)

func TestVlidatePasswordFalse(t *testing.T) {
	validate:= validator.New()
	validate.RegisterValidation("password", validatePassword)
	arr:=[]string{"Wojtekir","asdd","asdasdaq","12345678","Wojtekqw+_';\"["}
	for i:=0;i<5;i++ {
		err:=validate.Var(arr[i],"password")
		if err==nil {
			t.Error(err)
		}
	} 
	
}

func TestVlidatePasswordTrue(t *testing.T) {
	validate:= validator.New()
	validate.RegisterValidation("password", validatePassword)
	arr:=[]string{"Wojtek9r","asDddfgdfg887","asdasdaqE4","12345678Q","Wojtekqw4"}
	for i:=0;i<5;i++ {
		err:=validate.Var(arr[i],"password")
		if err!=nil {
			t.Error(err)
		}
	} 
	
}

func TestVlidatePostFalse(t *testing.T) {
	validate:= validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	arr:=[]string{"aw-asd","123435","234-87","34-876h","32y768"}
	for i:=0;i<5;i++ {
		err:=validate.Var(arr[i],"postCode")
		if err==nil {
			t.Error(err)
		}
	} 
	
}

func TestVlidatePostTrue(t *testing.T) {
	validate:= validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	arr:=[]string{"12-678","87-456","09-000","00-000","99-999"}
	for i:=0;i<5;i++ {
		err:=validate.Var(arr[i],"postCode")
		if err!=nil {
			t.Error(err)
		}
	} 
	
}

func TestValidationStructTrue (t *testing.T) {
	per:=[]Person{{"Andasd","asdasdP","asd.pooijk-asda@asdasd.pl","asdasd","12a","ssdfsdfQWEWEQ","45-456","Wertyui9"},
	{"ASDSSD","asdasd","asd.pooijk-asda@asdasd.asdasd.asd-asd.pl","ASADDSA","12243","WEQQWEasdasd","00-000","W+-==9tt"}}
	validate:= validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	validate.RegisterValidation("password", validatePassword)
	for i:=0;i<len(per);i++ {
		err:=validate.Struct(per[i])
		if err!=nil {
			t.Error(err)
		}
	}

}

func TestValidationStructFalse (t *testing.T) {
	per:=[]Person{{"Andasd","asdasdP","asd.pooijk-asda@asdasd.pl@","asdasd","12a","ssdfsdfQWEWEQ","45-456","Wertyui9"},
	{"ASDSSD3","asdasd","asd.pooijk-asda@asdasd.asdasd.asd-asd.pl","ASADDSA","12243","WEQQWEasdasd","00-000","W+-==9tt"}}
	validate:= validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	validate.RegisterValidation("password", validatePassword)
	for i:=0;i<len(per);i++ {
		err:=validate.Struct(per[i])
		if err==nil {
			t.Error(err)
		}
	}

}