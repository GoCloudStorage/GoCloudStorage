package model

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	pwd := "est1111"
	p1, _ := HashPassword(pwd)
	pwd2 := "est1111"
	p2, _ := HashPassword(pwd2)
	fmt.Println(pwd, pwd2, p1, p2)
	if p1 != p2 {
		t.Fatal("failed to impl hashpassword()")
	}

}
