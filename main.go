package main

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var receiptIdVsPoints = make(map[string]int)

type item struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}
type receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

/*
One point for every alphanumeric character in the retailer name.
*/
func countAlphanumericCharacters(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

/*
50 points if the total is a round dollar amount with no cents.
25 points if the total is a multiple of 0.25.
*/
func checkIfTotalPointsAreRounded(totalStr string) int {
	totalPoints, err := strconv.ParseFloat(totalStr, 64)
	var points = 0
	if err == nil {
		if math.Mod(totalPoints, 1) == 0 {
			points = 50
		}
		if math.Mod(totalPoints, 0.25) == 0 {
			points += 25
		}
	}
	return points
}

/*
5 points for every two items on the receipt.
*/
func calculateTotalNumberOfItemsInReceipt(lengthOfItems int) int {
	return (lengthOfItems / 2) * 5
}

/*
If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and
round up to the nearest integer.
The result is the number of points earned.
*/
func calculatePointsBasedOnItemDescription(items []item) int {
	var priceValue = 0
	for _, itemVal := range items {
		descptTrimmed := strings.TrimSpace(itemVal.Description)
		if len(descptTrimmed)%3 == 0 {
			itemPrice, err := strconv.ParseFloat(itemVal.Price, 64)
			if err == nil {
				priceValue += int(math.Ceil(itemPrice * 0.2))
			}
		}
	}
	return priceValue
}

/*
6 points if the day in the purchase date is odd.
*/
func checkIfDateIsOdd(dateStr string) int {
	layout := "2006-01-02"
	priceVal := 0
	date, err := time.Parse(layout, dateStr)
	if err == nil && date.Day()%2 != 0 {
		priceVal = 6
	}
	return priceVal
}

/*
10 points if the time of purchase is after 2:00pm and before 4:00pm.
*/
func checkIfTimeIsBetweenTwoAndFour(timeStr string) int {
	layout := "15:04"
	priceVal := 0
	parsedTime, err := time.Parse(layout, timeStr)
	if err == nil && ((parsedTime.Hour() == 15) || (parsedTime.Hour() == 14 && parsedTime.Minute() > 0)) {
		priceVal = 10
	}
	return priceVal
}

func calculatePoints(newReceipt receipt) (totalCalculatedPoints int) {

	totalCalculatedPoints = 0
	totalCalculatedPoints += countAlphanumericCharacters(newReceipt.Retailer)
	totalCalculatedPoints += checkIfTotalPointsAreRounded(newReceipt.Total)
	totalCalculatedPoints += calculateTotalNumberOfItemsInReceipt(len(newReceipt.Items))
	totalCalculatedPoints += calculatePointsBasedOnItemDescription(newReceipt.Items)
	totalCalculatedPoints += checkIfDateIsOdd(newReceipt.PurchaseDate)
	totalCalculatedPoints += checkIfTimeIsBetweenTwoAndFour(newReceipt.PurchaseTime)

	return totalCalculatedPoints
}
func processReceipts(context *gin.Context) {
	var newReceipt receipt
	err := context.BindJSON(&newReceipt)
	if err != nil {
		return
	}
	var generatedId = generateRandomString(8) + "-" + generateRandomString(4) + "-" + generateRandomString(4) + "-" + generateRandomString(4) + "-" + generateRandomString(12)

	var points = calculatePoints(newReceipt)
	receiptIdVsPoints[generatedId] = points
	context.IndentedJSON(http.StatusCreated, gin.H{"id": generatedId})
}

func getPoints(context *gin.Context) {
	id := context.Param("id")
	var points = receiptIdVsPoints[id]

	points, keyexists := receiptIdVsPoints[id]
	if keyexists {
		context.IndentedJSON(http.StatusOK, gin.H{"points": points})
	} else {
		context.IndentedJSON(http.StatusNotFound, gin.H{"id": "id not found"})
	}

}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipts)
	router.GET("/receipts/:id/points", getPoints)
	router.Run("localhost:8080")

}
