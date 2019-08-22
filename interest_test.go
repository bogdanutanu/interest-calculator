package interest

import (
	"testing"
	"time"
)

func TestCalculateInterest(t *testing.T) {
	transactions := [][]string{
		[]string{"2019-01-01", "Credit", "1000", "1000"},
		[]string{"2019-06-01", "Credit", "1000", "2000"},
	}
	until, err := time.Parse("2006-01-02", "2020-01-01")

	interest, err := calculateInterest(transactions, "1", until)
	if err != nil {
		t.Fatalf("Error calculating interest: %v", err)
	}

	if interest.String() != "11.73" {
		t.Fatalf("Expected interest to be 11.73 but was '%s'", interest.String())
	}
}

func TestCalculateInterestSingleDeposit(t *testing.T) {
	transactions := [][]string{
		[]string{"2019-01-01", "Credit", "1000", "1000"},
	}
	until, err := time.Parse("2006-01-02", "2020-01-01")

	interest, err := calculateInterest(transactions, "1", until)
	if err != nil {
		t.Fatalf("Error calculating interest: %v", err)
	}

	if interest.String() != "10" {
		t.Fatalf("Expected interest to be 10 but was '%s'", interest.String())
	}
}
