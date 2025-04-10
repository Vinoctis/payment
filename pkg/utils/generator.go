package utils

import(
	"math/rand"
	"time"
	"fmt"
)
func GenerateOrderID() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().UnixNano()
	randomNum := fmt.Sprintf("%04d", rand.Intn(10000))
	return time.Now().Format("20060102") + 
			fmt.Sprintf("%d", timestamp) +
	 		randomNum
}