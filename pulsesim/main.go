package main

import (
	"flag"
	"github.com/go-faker/faker/v4"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	NumUsers int
	BaseUrl  string
}

func intEnvOrDefault(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}

func stringEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	cfg := config{}
	// Parse command-line flags
	flag.IntVar(&cfg.NumUsers, "n", intEnvOrDefault("NUM_USERS", 100), "Number of users to simulate")
	flag.StringVar(&cfg.BaseUrl, "api-url", stringEnvOrDefault("API_URL", "http://localhost:8080"), "Backend URL to ping")
	flag.Parse()

	log.SetFlags(log.Ltime)
	// Register and login all users
	log.Printf("Registering and logging in %d users...\n", cfg.NumUsers)
	users := make([]user, cfg.NumUsers)

	for i := 0; i < cfg.NumUsers; i++ {
		name := faker.Name()
		loginResponse, err := registerAndLogin(name, cfg.BaseUrl)
		if err != nil {
			panic(err)
		}
		users[i] = user{name, loginResponse.SessionID}
		log.Printf("%s: registered and logged in\n", name)
	}

	// Initialize random seed
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch a goroutine for each user
	for _, activeUser := range users {
		wg.Add(1)
		go func(user user) {
			defer wg.Done()

			for {
				// Generate a random interval between 5 and 15 seconds
				// This ensures a fair mix of online (< 10s) and offline (>= 10s) statuses
				interval := 5 + random.Intn(11) // 5 to 15 seconds

				// Ping the server
				err := Ping(user.SessionId, cfg.BaseUrl)
				if err != nil {
					log.Printf("%s ping failed: %v\n", user.Name, err)
				} else {
					log.Printf("%s: pinged successfully (next ping in %ds)\n",
						user.Name, interval)
				}

				// Sleep for the random interval
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}(activeUser)
	}
	wg.Wait()
}

func registerAndLogin(name string, baseUrl string) (LoginResponse, error) {

	registerResponse, err := Register(name, baseUrl)
	if err != nil {
		return LoginResponse{}, err
	}

	loginResponse, err := Login(registerResponse.Id, registerResponse.Token, baseUrl)
	if err != nil {
		return LoginResponse{}, err
	}

	return loginResponse, nil
}

type user struct {
	Name      string
	SessionId string
}
