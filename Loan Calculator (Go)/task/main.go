package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
)

var errInput = errors.New("Incorrect parameters")

const (
	annuity = "annuity"
	diff    = "diff"
)

func main() {
	var principal float64
	flag.Float64Var(&principal, "principal", 0, "loan principal")

	var payment float64
	flag.Float64Var(&payment, "payment", 0, "monthly payment")

	var months int
	flag.IntVar(&months, "periods", 0, "number of months")

	var interest float64
	flag.Float64Var(&interest, "interest", 0, "annual interest rate")

	var typ string
	flag.StringVar(&typ, "type", "", "type of interest")

	flag.Parse()

	if err := validateInterest(interest); err != nil {
		fmt.Println(err)
		return
	}

	if err := validateOption(typ); err != nil {
		fmt.Println(err)
		return
	}

	loan := Loan{interest: interest}

	switch typ {
	case diff:
		payments, overpayment, err := loan.calculateDiff(principal, months)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, payment := range payments {
			fmt.Printf("Month %d: payment is %.0f\n", i+1, payment)
		}

		if overpayment > 0 {
			fmt.Println("Overpayment = ", overpayment)
		}

		return
	case annuity:
		if principal > 0 && months > 0 {
			payment := loan.calculatePayment(principal, months)
			fmt.Println(fmt.Sprintf("Your annuity payment = %.0f!", payment))
			overpayment := payment*float64(months) - principal
			if overpayment > 0 {
				fmt.Println("Overpayment = ", overpayment)
			}
			return
		}

		if months > 0 && payment > 0 {
			principal := loan.calculatePrincipal(payment, months)
			fmt.Println(fmt.Sprintf("Your loan principal = %.0f!", principal))
			overpayment := payment*float64(months) - principal
			if overpayment > 0 {
				fmt.Println("Overpayment = ", overpayment)
			}
			return
		}

		if principal > 0 && payment > 0 {
			months := loan.calculatePeriod(principal, payment)
			years, restMonths := convertToYears(months)
			if years > 0 && restMonths > 0 {
				fmt.Println(fmt.Sprintf("It will take %d years and %d months to repay this loan!", years, months))
			} else if years > 0 {
				fmt.Println(fmt.Sprintf("It will take %d years to repay this loan!", years))
			} else {
				fmt.Println(fmt.Sprintf("It will take %d months to repay this loan!", restMonths))
			}

			overpayment := payment*float64(months) - principal
			if overpayment > 0 {
				fmt.Println("Overpayment = ", overpayment)
			}
			return
		}

		fmt.Println(errInput)
	}
}

type Loan struct {
	interest float64
}

func (l *Loan) calculateNominalInterest() float64 {
	return l.interest / 1200
}

func (l *Loan) calculateDiff(principal float64, months int) ([]float64, float64, error) {
	if principal <= 0 || months <= 0 {
		return nil, 0, errInput
	}

	nominalInterest := l.calculateNominalInterest()
	payments := make([]float64, months)
	var sum float64
	for m := 0; m < months; m++ {
		dm := (principal / float64(months)) + nominalInterest*(principal-(principal*float64(m))/float64(months))
		payments[m] = math.Ceil(dm)
		sum += payments[m]
	}

	overpayment := sum - principal
	if overpayment > 0 {
		return payments, overpayment, nil
	}

	return payments, 0, nil
}

func (l *Loan) calculatePayment(principal float64, months int) float64 {
	nominalInterest := l.calculateNominalInterest()
	pow := math.Pow(1+nominalInterest, float64(months))
	return math.Ceil(principal * (nominalInterest * pow) / (pow - 1))
}

func (l *Loan) calculatePeriod(principal, payment float64) int {
	nominalInterest := l.calculateNominalInterest()
	log := math.Log(payment/(payment-nominalInterest*principal)) / math.Log(1+nominalInterest)
	return int(math.Ceil(log))
}

func (l *Loan) calculatePrincipal(payment float64, months int) float64 {
	nominalInterest := l.calculateNominalInterest()
	pow := math.Pow(1+nominalInterest, float64(months))
	return payment / ((nominalInterest * pow) / (pow - 1))
}

func convertToYears(months int) (int, int) {
	return int(math.Floor(float64(months) / 12)), int(months) % 12
}

func validateInterest(interest float64) error {
	if interest <= 0 {
		return errInput
	}

	return nil
}

func validateOption(option string) error {
	if option != diff && option != annuity {
		return errInput
	}

	return nil
}
