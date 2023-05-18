package tests

import "testing"

func TestAllTests(t *testing.T) {
	t.Run("testLogin", testLogin)
	t.Run("testRegister", testRegister)
	t.Run("testGetUserMe", testGetUserMe)
	t.Run("testModifyUserMe", testModifyUserMe)
	t.Run("testListUsers", testListUsers)
}
