package utils

import (
	"fmt"
	"testing"
)

func TestMakePassword(t *testing.T) {
	fmt.Println(MakePassword("adminadmin"))
}
