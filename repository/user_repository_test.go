package repository_test

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	mockUser := &domain.User{
		ID:       "user-1",
		Name:     "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockDB.ExpectExec(`INSERT INTO users`).WithArgs(mockUser.ID, mockUser.Name, mockUser.Email, mockUser.Password).WillReturnResult(sqlmock.NewResult(1, 1))

		ur := repository.NewUserRepository(db)

		err = ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockDB.ExpectExec(`INSERT INTO users`).WithArgs(mockUser.ID, mockUser.Name, mockUser.Email, mockUser.Password).WillReturnError(errors.New("Unexpected"))

		ur := repository.NewUserRepository(db)

		err = ur.Create(context.Background(), mockUser)

		assert.Error(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

}
