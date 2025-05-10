package main

import (
	"app/src/config"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Simple utility to verify Midtrans signatures
// Usage: go run src/cmd/verify_signature.go <order_id> <status_code> <gross_amount> <signature_key>

func main() {
	// Check if enough arguments are provided
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run src/cmd/verify_signature.go <order_id> <status_code> <gross_amount> <signature_key>")
		os.Exit(1)
	}

	// Get arguments
	orderID := os.Args[1]
	statusCode := os.Args[2]
	grossAmount := os.Args[3]
	signatureKey := os.Args[4]

	// Get server key (no need to load config, init() in config package handles it)
	serverKey := config.MidtransServerKey

	fmt.Println("======= MIDTRANS SIGNATURE VERIFICATION =======")
	fmt.Println("Order ID:", orderID)
	fmt.Println("Status Code:", statusCode)
	fmt.Println("Gross Amount:", grossAmount)
	fmt.Printf("Server Key (masked): %s***%s\n", serverKey[:4], serverKey[len(serverKey)-4:])

	// Try different formats
	fmt.Println("\n======= TESTING DIFFERENT FORMATS =======")

	// Format 1: Direct concatenation with string values
	data1 := orderID + statusCode + grossAmount
	sig1 := generateSignature(data1, serverKey)
	isValid1 := sig1 == signatureKey
	fmt.Printf("Format 1 (direct string concat): %s + %s + %s\n", orderID, statusCode, grossAmount)
	fmt.Printf("Signature: %s\n", sig1)
	fmt.Printf("Valid: %v\n\n", isValid1)

	// Format 2: Clean amounts from decimal point
	grossAmountClean := strings.Replace(grossAmount, ".", "", -1)
	data2 := orderID + statusCode + grossAmountClean
	sig2 := generateSignature(data2, serverKey)
	isValid2 := sig2 == signatureKey
	fmt.Printf("Format 2 (no decimal point): %s + %s + %s\n", orderID, statusCode, grossAmountClean)
	fmt.Printf("Signature: %s\n", sig2)
	fmt.Printf("Valid: %v\n\n", isValid2)

	// Format 3: Just the numbers
	data3 := orderID + statusCode + "10000"
	sig3 := generateSignature(data3, serverKey)
	isValid3 := sig3 == signatureKey
	fmt.Printf("Format 3 (raw number): %s + %s + %s\n", orderID, statusCode, "10000")
	fmt.Printf("Signature: %s\n", sig3)
	fmt.Printf("Valid: %v\n\n", isValid3)

	// Format 4: Decimal with precision 00
	data4 := orderID + statusCode + "10000.00"
	sig4 := generateSignature(data4, serverKey)
	isValid4 := sig4 == signatureKey
	fmt.Printf("Format 4 (decimal with precision): %s + %s + %s\n", orderID, statusCode, "10000.00")
	fmt.Printf("Signature: %s\n", sig4)
	fmt.Printf("Valid: %v\n\n", isValid4)

	// Exit with success if any format is valid
	if isValid1 || isValid2 || isValid3 || isValid4 {
		fmt.Println("✅ Valid signature found!")
		os.Exit(0)
	}

	fmt.Println("❌ No valid signature format found!")
	os.Exit(1)
}

// Helper function to generate signature
func generateSignature(data, serverKey string) string {
	signatureData := data + serverKey
	h := sha512.New()
	h.Write([]byte(signatureData))
	return hex.EncodeToString(h.Sum(nil))
}
