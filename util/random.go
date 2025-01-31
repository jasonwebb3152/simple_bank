package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	/** Generates a random integer from min to max. */
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	/** Generates a random string of length n. */
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	/** Generates a random string owner name of length 6. */
	return RandomString(6)
}

func RandomMoney() int64 {
	/** Generates a random amount of money. */
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	/** Generates a random currency code */
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
