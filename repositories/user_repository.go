package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"user-service/models"

	"github.com/Masterminds/squirrel"

	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	DB           *sql.DB
	QueryBuilder squirrel.StatementBuilderType
}

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrUserNotFound      = errors.New("user not found")
)

var validate = validator.New()

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB:           db,
		QueryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query, args, err := r.QueryBuilder.
		Select("id", "user_name", "email", "first_name", "last_name", "user_status", "department").
		From("users").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FirstName,
			&user.LastName, &user.Status, &user.Department); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Validate the user input
	if err := validate.Struct(user); err != nil {
		// If validation fails, return the first error message
		return fmt.Errorf("validation failed: %w", err)
	}

	query, args, err := r.QueryBuilder.
		Insert("users").
		Columns("user_name", "email", "first_name", "last_name", "user_status", "department").
		Values(user.UserName, user.Email, user.FirstName, user.LastName, user.Status, user.Department).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err) // Wrap the error
	}

	_, execErr := r.DB.Exec(query, args...)
	if execErr != nil {
		// Check if the error is a duplicate key error
		if execErr.Error() == "duplicate username" || isUniqueConstraintViolation(execErr) {
			// Wrap the error with ErrDuplicateUsername so it can be detected by errors.Is
			return fmt.Errorf("%w", ErrDuplicateUsername)
		}
		// Wrap the error and add context
		return fmt.Errorf("failed to execute query: %w", execErr)
	}
	return nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	// Validate the user input
	if err := validate.Struct(user); err != nil {
		// If validation fails, return the first error message
		return fmt.Errorf("validation failed: %w", err)
	}

	// Prepare the update query using squirrel
	query, args, err := squirrel.Update("users").
		Set("user_name", user.UserName).
		Set("email", user.Email).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("user_status", user.Status).
		Set("department", user.Department).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return err
	}

	// Execute the update query
	result, err := ur.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query, args, err := r.QueryBuilder.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	res, execErr := r.DB.Exec(query, args...)
	if execErr != nil {
		return execErr
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query, args, err := r.QueryBuilder.
		Select("id", "user_name", "email", "first_name", "last_name", "user_status", "department").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Check if we have any rows and scan them into a User struct
	if rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.FirstName,
			&user.LastName, &user.Status, &user.Department); err != nil {
			return nil, err
		}
		// Return the user if found
		return &user, nil
	}

	// Return an error if no rows were found
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}

// isUniqueConstraintViolation checks if the error is a unique constraint violation error
func isUniqueConstraintViolation(err error) bool {
	// Check for specific error related to unique constraint violation
	// For SQLite, check for a constraint violation
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true
		}
	}
	// Add checks for other DB types if needed (e.g., PostgreSQL, MySQL)
	return false
}
