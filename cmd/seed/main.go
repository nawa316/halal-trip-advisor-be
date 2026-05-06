package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/idutil"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	app := bootstrap.App()
	db := app.DB
	defer app.CloseDBConnection()

	userRepo := repository.NewUserRepository(db)

	usersToSeed := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "adminpassword",
			Role:     "admin",
		},
		{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: "userpassword",
			Role:     "user",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("--- Starting Seeder ---")

	for _, u := range usersToSeed {
		// Check if user already exists
		_, err := userRepo.GetByEmail(ctx, u.Email)
		if err == nil {
			fmt.Printf("[Skipping] User with email %s already exists.\n", u.Email)
			continue
		}

		// Encrypt password
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[Error] Failed to hash password for %s: %v\n", u.Email, err)
			continue
		}

		user := domain.User{
			ID:       idutil.NewID(),
			Name:     u.Name,
			Email:    u.Email,
			Password: string(encryptedPassword),
			Role:     u.Role,
		}

		// Create user
		err = userRepo.Create(ctx, &user)
		if err != nil {
			log.Printf("[Error] Failed to create user %s: %v\n", u.Email, err)
			continue
		}

		fmt.Printf("[Success] Created %s user: %s (Password: %s)\n", u.Role, u.Email, u.Password)
	}

	fmt.Println("--- Seeding Completed ---")
}
