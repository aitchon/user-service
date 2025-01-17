package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"user-service/controllers"
	"user-service/models"
	"user-service/repositories"
	"user-service/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")
}

var _ = Describe("UserController and UserService with Squirrel Repository", func() {
	var (
		db          *sql.DB
		mock        sqlmock.Sqlmock
		userRepo    *repositories.UserRepository
		userService *services.UserService
		e           *echo.Echo
		rec         *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).To(BeNil())
		userRepo = repositories.NewUserRepository(db)
		userService = services.NewUserService(userRepo)
		e = echo.New()
		rec = httptest.NewRecorder()
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("UserService", func() {
		Context("CreateUser", func() {
			It("should call the repository's CreateUser method successfully", func() {
				// Arrange
				user := &models.User{UserName: "john_doe", Email: "john@example.com", FirstName: "John", LastName: "Doe", Status: "A", Department: "IT"}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status, user.Department).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Act
				err := userService.CreateUser(user)

				// Assert
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})

			It("should return an error if the repository returns an error", func() {
				// Arrange
				user := &models.User{UserName: "john_doe", Email: "john@example.com", FirstName: "Jane", LastName: "Doe", Status: "I", Department: "IT"}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status, user.Department).
					WillReturnError(errors.New("duplicate username"))

				// Act
				err := userService.CreateUser(user)

				// Assert
				Expect(err).To(MatchError("duplicate username"))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})
		})
	})

	Describe("UserController", func() {
		Context("CreateUser", func() {
			It("should return 201 when user is created successfully", func() {
				// Arrange
				user := models.User{
					UserName:   "john_doe",
					Email:      "john@example.com",
					FirstName:  "John",
					LastName:   "Doe",
					Status:     "A",
					Department: "IT",
				}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status, user.Department).
					WillReturnResult(sqlmock.NewResult(1, 1))

				handler := controllers.CreateUser(userService)
				body, _ := json.Marshal(user)
				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)

				// Act
				err := handler(c)

				// Assert
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusCreated))
				var createdUser models.User
				json.Unmarshal(rec.Body.Bytes(), &createdUser)
				Expect(createdUser.UserName).To(Equal(user.UserName))
				Expect(createdUser.Email).To(Equal(user.Email))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})

			It("should return 409 when username already exists", func() {
				// Arrange
				user := models.User{
					UserName:   "existing_user",
					Email:      "existing@example.com",
					FirstName:  "Existing",
					LastName:   "User",
					Status:     "A",
					Department: "IT",
				}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status, user.Department).
					WillReturnError(errors.New("duplicate username"))

				handler := controllers.CreateUser(userService)
				body, _ := json.Marshal(user)
				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)

				// Act
				err := handler(c)

				// Assert
				Expect(err).NotTo(BeNil())
				httpErr, ok := err.(*echo.HTTPError)
				Expect(ok).To(BeTrue())
				Expect(httpErr.Code).To(Equal(http.StatusConflict))
				Expect(httpErr.Message).To(Equal("username already exists"))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})

			It("should return 400 when input is invalid", func() {
				// Arrange
				invalidBody := `{"username": 123}` // Invalid JSON structure
				handler := controllers.CreateUser(userService)
				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(invalidBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)

				// Act
				err := handler(c)

				// Assert
				httpErr, ok := err.(*echo.HTTPError)
				Expect(ok).To(BeTrue())
				Expect(httpErr.Code).To(Equal(http.StatusBadRequest))
				Expect(httpErr.Message).To(Equal("Invalid input"))
			})
		})
		Describe("UpdateUser", func() {
			It("should return 200 when the user is successfully updated", func() {
				// Arrange
				user := models.User{
					ID:         1,
					UserName:   "updated_user",
					Email:      "updated_email@example.com",
					FirstName:  "Updated",
					LastName:   "User",
					Status:     "A",
					Department: "IT",
				}
				handler := controllers.UpdateUser(userService)
				body, _ := json.Marshal(user)
				req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues("1")

				mock.ExpectExec(`UPDATE users SET user_name = \?, email = \?, first_name = \?, last_name = \?, user_status = \?, department = \? WHERE id = \?`).
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status,
						user.Department, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Act
				err := handler(c)

				// Assert
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(mock.ExpectationsWereMet()).To(BeNil())
			})

			It("should return 404 when the user is not found", func() {
				// Arrange
				user := models.User{
					ID:         999,
					UserName:   "updated_user",
					Email:      "updated_email@example.com",
					FirstName:  "Updated",
					LastName:   "User",
					Status:     "A",
					Department: "IT",
				}
				handler := controllers.UpdateUser(userService)
				body, _ := json.Marshal(user)
				req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues("999")

				mock.ExpectExec(`UPDATE users SET user_name = \?, email = \?, first_name = \?, last_name = \?, user_status = \?, department = \? WHERE id = \?`).
					WithArgs(user.UserName, user.Email, user.FirstName, user.LastName, user.Status,
						user.Department, 999).
					WillReturnError(errors.New("user not found"))

				// Act
				err := handler(c)

				// Assert
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusNotFound))
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal("User not found"))
			})
		})

		Describe("DeleteUser", func() {
			It("should return 204 when the user is successfully deleted", func() {
				// Arrange
				handler := controllers.DeleteUser(userService)
				req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues("1")

				// mock service behavior
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Act
				err := handler(c)

				// Assert
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusNoContent))
			})

			It("should return 404 when the user is not found", func() {
				// Arrange
				handler := controllers.DeleteUser(userService)
				req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues("999")

				// Mock service behavior
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs(999).
					WillReturnError(errors.New("user not found"))

				// Act
				err := handler(c)

				// Assert
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusNotFound))
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal("User not found"))
			})
		})
	})

})
