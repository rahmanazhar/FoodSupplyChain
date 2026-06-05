// Command token mints a development JWT for testing the protected endpoints.
//
//	go run ./cmd/token -secret "your-secret-key-here" -role admin
//
// The secret must match the running service's auth.jwt_secret (or JWT_SECRET).
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
)

func main() {
	secret := flag.String("secret", os.Getenv("JWT_SECRET"), "JWT signing secret (defaults to $JWT_SECRET)")
	subject := flag.String("sub", "test-user", "token subject (user id)")
	role := flag.String("role", auth.RoleAdmin, "role: admin | manager | operator | viewer")
	tenant := flag.String("tenant", "tenant-1", "tenant id")
	ttl := flag.Duration("ttl", time.Hour, "token lifetime")
	flag.Parse()

	if *secret == "" {
		fmt.Fprintln(os.Stderr, "error: provide -secret or set JWT_SECRET")
		os.Exit(1)
	}

	token, err := auth.NewManager(*secret, *ttl).GenerateToken(*subject, *role, *tenant)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(token)
}
