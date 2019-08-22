package interest

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func main() {
	csvFileName := "interest.csv"
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("Error opening file '%s': %v", csvFileName, err)
		os.Exit(1)
	}

	csvFileReader := csv.NewReader(bufio.NewReader(csvFile))
	transactions, err := csvFileReader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records from file '%s': %v", csvFileName, err)
		os.Exit(2)
	}

	interest, err := calculateInterest(transactions, "1.5", time.Now())
	if err != nil {
		log.Fatalf("Error calculating interest: %v", err)
		os.Exit(3)
	}
	fmt.Printf("%d", interest)
	os.Exit(0)
}

func calculateInterest(transactions [][]string, rate string, until time.Time) (decimal.Decimal, error) {
	earned := decimal.NewFromFloat(0)
	const dateLayout = "2006-01-02"

	for i := 0; i < len(transactions); i++ {
		if i == 0 && len(transactions) > 1 {
			continue
		}

		if len(transactions[i]) < 4 {
			return decimal.Decimal{}, fmt.Errorf("Expected at least 4 fields but got %v", transactions[i])
		}
		var previousDate, currentDate time.Time
		var err error
		if i < len(transactions)-1 {
			previousDate, err = time.Parse(dateLayout, transactions[i-1][0])
			if err != nil {
				return decimal.Decimal{}, errors.Wrapf(err, "Error parsing '%s' in format '%s'", transactions[i-1][0], dateLayout)
			}

			currentDate, err = time.Parse(dateLayout, transactions[i][0])
			if err != nil {
				return decimal.Decimal{}, errors.Wrapf(err, "Error parsing '%s' in format '%s'", transactions[i][0], dateLayout)
			}
		} else {
			// For the last transaction, we calculate the interest from its
			// date until today
			previousDate, err = time.Parse(dateLayout, transactions[i][0])
			if err != nil {
				return decimal.Decimal{}, errors.Wrapf(err, "Error parsing '%s' in format '%s'", transactions[i][0], dateLayout)
			}

			currentDate = until
		}

		durationInDays := currentDate.Sub(previousDate).Hours() / 24
		currentBalance, err := decimal.NewFromString(transactions[i][3])
		if err != nil {
			return decimal.Decimal{}, errors.Wrapf(err, "parsing '%s' as decimal", transactions[i][3])
		}
		aprDecimal, err := decimal.NewFromString(rate)
		if err != nil {
			return decimal.Decimal{}, errors.Wrapf(err, "parsing '%s' as decimal", rate)
		}
		dailyRate := calculateDailyInterestRate(aprDecimal)
		earned = earned.Add(currentBalance.Mul(decimal.NewFromFloat(durationInDays).Mul(dailyRate)))

	}
	return earned.RoundBank(2), nil
}

func calculateDailyInterestRate(percentageRate decimal.Decimal) decimal.Decimal {
	return percentageRate.Div(decimal.NewFromFloat(100)).Div(decimal.NewFromFloat(365))
}
