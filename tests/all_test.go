package tests

import "testing"

func TestAllTests(t *testing.T) {
	t.Run("testLogin", testLogin)
	t.Run("testRegister", testRegister)
	t.Run("testUserMe", testUserMe)
	t.Run("testUserList", testUserList)
}
