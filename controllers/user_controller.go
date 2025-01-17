package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"user-service/models"
	"user-service/services"

	"github.com/labstack/echo/v4"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

// @Summary Get all users
// @Description Get a list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(service *services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := service.GetAllUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		}
		return c.JSON(http.StatusOK, users)
	}
}

// @Summary Create a new user
// @Description Add a new user to the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func CreateUser(service *services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if err := service.CreateUser(&user); err != nil {
			if err.Error() == "duplicate username" {
				return echo.NewHTTPError(http.StatusConflict, "username already exists")
			}
			// Check if the error is a validation failure
			if strings.HasPrefix(err.Error(), "validation failed:") {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
		}
		return c.JSON(http.StatusCreated, user)
	}
}

// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(service *services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user ID from request
		id := c.Param("id")

		// Convert id to integer
		userID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Prepare delete query
		if err := service.DeleteUser(userID); err != nil {
			if err.Error() == "user not found" {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

// @Summary Update a user
// @Description Update a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(service *services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}
		// Validate input
		if user.UserName == "" || user.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if err := service.UpdateUser(&user); err != nil {
			if err.Error() == "user not found" {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}

			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
		return c.JSON(http.StatusOK, user) // Return updated user
	}

}
