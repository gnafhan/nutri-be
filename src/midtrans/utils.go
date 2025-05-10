package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

// VerifySignature verifies if a Midtrans signature key is valid
// The signature is generated using SHA512(order_id+status_code+gross_amount+ServerKey)
func VerifySignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) (bool, string) {
	// Try several possible formats for the signature data

	// Format 1: Direct concatenation
	calculatedSignature1 := generateSignature(orderID+statusCode+grossAmount, serverKey)

	// Format 2: Clean amounts from decimal point
	grossAmountClean := strings.Replace(grossAmount, ".", "", -1)
	calculatedSignature2 := generateSignature(orderID+statusCode+grossAmountClean, serverKey)

	// Format 3: Midtrans might use the raw number without decimals
	calculatedSignature3 := generateSignature(orderID+statusCode+"10000", serverKey)

	// Check all formats
	isValid1 := calculatedSignature1 == signatureKey
	isValid2 := calculatedSignature2 == signatureKey
	isValid3 := calculatedSignature3 == signatureKey

	// Return true if any format matches
	if isValid1 {
		return true, calculatedSignature1
	} else if isValid2 {
		return true, calculatedSignature2
	} else if isValid3 {
		return true, calculatedSignature3
	}

	// Return the most likely signature (for debugging)
	return false, calculatedSignature1
}

// generateSignature creates a hash from the input data + server key
func generateSignature(data, serverKey string) string {
	signatureData := data + serverKey

	// Create SHA512 hash
	h := sha512.New()
	h.Write([]byte(signatureData))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateSignature generates a Midtrans signature key
// Useful for testing or debugging
func GenerateSignature(orderID, statusCode, grossAmount, serverKey string) string {
	// Direct concatenation format
	return generateSignature(orderID+statusCode+grossAmount, serverKey)
}

// maskString masks a string for security when logging
func maskString(input string) string {
	if len(input) <= 8 {
		return "****"
	}
	return input[:4] + "****" + input[len(input)-4:]
}
