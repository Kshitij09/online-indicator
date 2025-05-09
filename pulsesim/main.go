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

func main() {
	cfg := config{}
	// Parse command-line flags
	flag.IntVar(&cfg.NumUsers, "n", 100, "Number of users to simulate")
	flag.StringVar(&cfg.BaseUrl, "api-url", "http://localhost:8080", "Backend URL to ping")
	flag.Parse()

	envNumUsers, set := os.LookupEnv("NUM_USERS")
	if set {
		if envNumUsers, err := strconv.Atoi(envNumUsers); err != nil {
			cfg.NumUsers = envNumUsers
		}
	}
	envBaseUrl, set := os.LookupEnv("API_URL")
	if set {
		cfg.BaseUrl = envBaseUrl
	}

	log.SetFlags(log.Ltime)
	// Register and login all users
	log.Printf("Registering and logging in %d users...\n", cfg.NumUsers)
	users := make([]user, cfg.NumUsers)

	for i := 0; i < cfg.NumUsers; i++ {
		name := faker.Name()
		id, loginResponse, err := registerAndLogin(name, cfg.BaseUrl)
		if err != nil {
			panic(err)
		}
		users[i] = user{id, name, loginResponse.SessionToken}
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
				err := Ping(user.Id, user.SessionId, cfg.BaseUrl)
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

func registerAndLogin(name string, baseUrl string) (string, LoginResponse, error) {

	registerResponse, err := Register(name, baseUrl)
	if err != nil {
		return "", LoginResponse{}, err
	}

	loginResponse, err := Login(registerResponse.Id, registerResponse.ApiKey, baseUrl)
	if err != nil {
		return "", LoginResponse{}, err
	}

	return registerResponse.Id, loginResponse, nil
}

type user struct {
	Id        string
	Name      string
	SessionId string
}
