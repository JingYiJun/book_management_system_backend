package tests

import (
	"book_management_system_backend/apis"
	. "book_management_system_backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testLogin(t *testing.T) {
	var response = Map{}
	defaultTester.testPost(t, "/api/login", 200, Map{
		"username": "admin",
		"password": "adminadmin",
	}, &response)
	assert.Equal(t, "登录成功", response["message"])
	assert.NotEmpty(t, response["access"])
	superAdminTester.Token = response["access"].(string)

	defaultTester.testPost(t, "/api/login", 400, Map{
		"username": "user",
		"password": "12345678",
	}, &response)
	assert.Equal(t, "用户名或密码错误", response["message"])
}

func testRegister(t *testing.T) {
	var response = Map{}
	superAdminTester.testPost(t, "/api/register", 200, Map{
		"username": "user",
		"password": "12345678",
	}, &response)
	assert.Equal(t, "注册成功", response["message"])
	assert.NotEmpty(t, response["access"])
	adminTester.Token = response["access"].(string)

	superAdminTester.testPost(t, "/api/register", 400, Map{
		"username": "user",
		"password": "123456789",
	}, &response)
	assert.Equal(t, "用户名已存在", response["message"])
}

func testGetUserMe(t *testing.T) {
	var user = User{}
	superAdminTester.testGet(t, "/api/users/me", 200, nil, &user)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, true, user.IsAdmin)

	adminTester.testGet(t, "/api/users/me", 200, nil, &user)
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, false, user.IsAdmin)

	defaultTester.testGet(t, "/api/users/me", 401, nil, nil)
}

func testModifyUserMe(t *testing.T) {
	var userResponse apis.UserResponse
	var dbUser User

	// test partial update self
	superAdminTester.testPatch(t, "/api/users/me", 200, Map{
		"real_name": "小明",
		"staff_id":  "s1001",
		"something": "something",
	}, &userResponse)
	assert.Equal(t, "admin", userResponse.Username)
	assert.Equal(t, true, userResponse.IsAdmin)
	assert.Equal(t, "小明", userResponse.RealName)
	assert.Equal(t, "s1001", userResponse.StaffID)
	DB.First(&dbUser, userResponse.ID)
	assert.Equal(t, "小明", dbUser.RealName)
	assert.Equal(t, "s1001", dbUser.StaffID)

	defaultTester.testPatch(t, "/api/users/me", 401, Map{
		"username": "user",
		"password": "user",
	}, nil)
}

func testListUsers(t *testing.T) {
	var usersResponse apis.UserListResponse
	superAdminTester.testGet(t, "/api/users", 200, Map{
		"page_num":  1,
		"page_size": 10,
	}, &usersResponse)
	assert.Equal(t, 2, len(usersResponse.Users))

	adminTester.testGet(t, "/api/users", 403, Map{
		"page_num":  1,
		"page_size": 10,
	}, nil)

	defaultTester.testGet(t, "/api/users", 401, Map{
		"page_num":  1,
		"page_size": 10,
	}, nil)
}
