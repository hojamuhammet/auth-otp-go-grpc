package utils

import (
	"math/rand"
	"strconv"
	"sync"
)

// Define a global map to store recently generated OTPs
var recentOTPs = make(map[string]struct{})

// Mutex for concurrent access to the recentOTPs map
var recentOTPsMutex sync.Mutex

// Function to generate a random OTP
func GenerateRandomOTP() string {
    // Define the range for OTPs
    min := 100000
    max := 999999

    // Generate a new OTP until it is unique
    for {
        otp := strconv.Itoa(rand.Intn(max-min+1) + min)

        // Check if the OTP is not in the recentOTPs map (i.e., it's unique)
        recentOTPsMutex.Lock()
        _, exists := recentOTPs[otp]
        if !exists {
            // If it's unique, add it to the recentOTPs map and return it
            recentOTPs[otp] = struct{}{}
            recentOTPsMutex.Unlock()
            return otp
        }
        recentOTPsMutex.Unlock()
    }
}
