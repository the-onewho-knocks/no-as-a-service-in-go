package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

// ---------- Config ----------
const (
	defaultPort = "8080"

	requestsPerMinute = 120
)

// ---------- Rate Limiter ----------
type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*clientLimiter)
	mu      sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), requestsPerMinute)
	clients[ip] = &clientLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

// Cleanup old IPs (memory hygiene)
func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// ---------- Utils ----------
func getClientIP(r *http.Request) string {
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// ---------- Main ----------
func main() {
	rand.Seed(time.Now().UnixNano())

	// Load reasons.json
	file, err := os.ReadFile("./no.json")
	if err != nil {
		log.Fatal("Failed to read reasons.json:", err)
	}

	var reasons []string
	if err := json.Unmarshal(file, &reasons); err != nil {
		log.Fatal("Invalid reasons.json:", err)
	}

	go cleanupClients()

	mux := http.NewServeMux()

	mux.HandleFunc("/no", func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too many requests, please try again later. (120 reqs/min/IP)",
			})
			return
		}

		reason := reasons[rand.Intn(len(reasons))]
		json.NewEncoder(w).Encode(map[string]string{
			"reason": reason,
		})
	})

	handler := cors.AllowAll().Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("No-as-a-Service is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
