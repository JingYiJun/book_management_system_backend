package tests

import (
	"book_management_system_backend/apis"
	. "book_management_system_backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testCreateABook(t *testing.T) {
	var bookCreateResponse apis.BookResponse
	superAdminTester.testPost(t, "/api/books", 201, Map{
		"title":  "testBook",
		"author": "testAuthor",
		"press":  "testPress",
		"isbn":   "90000000001",
		"price":  100,
	}, &bookCreateResponse)

	assert.Equal(t, "testBook", bookCreateResponse.Title)
}

func testGetABook(t *testing.T) {
	var bookGetResponse apis.BookListResponse
	superAdminTester.testGet(t, "/api/books", 200, Map{"id": 1}, &bookGetResponse)

	assert.Equal(t, 1, bookGetResponse.PageTotal)
	assert.Equal(t, "testBook", bookGetResponse.Books[0].Title)
}
