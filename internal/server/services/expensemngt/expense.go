package expensemngt

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/stac47/myroomies/internal/server/services"
	"github.com/stac47/myroomies/internal/server/services/usermngt"
	"github.com/stac47/myroomies/pkg/models"
)

type ExpenseListOptions struct {
	Before time.Time
	After  time.Time
}

func (o ExpenseListOptions) SelectAll() bool {
	return o.Before.IsZero() && o.After.IsZero()
}

func GetExpensesList(ctx context.Context, options ExpenseListOptions) []models.Expense {
	if options.SelectAll() {
		return services.GetDataAccess().GetExpenseDataAccess().RetrieveExpenses(ctx)
	} else {
		// TODO: Not all
		return services.GetDataAccess().GetExpenseDataAccess().RetrieveExpenses(ctx)
	}
}

func CreateExpense(ctx context.Context, authenticatedUser models.User, newExpense models.Expense) *models.Expense {
	var payer models.User
	if authenticatedUser.IsAdmin && newExpense.PayerLogin != "" {
		foundUser := usermngt.SearchUser(ctx, usermngt.ByLoginCriteria(newExpense.PayerLogin))
		if foundUser != nil {
			payer = *foundUser
		}
	}
	if payer.Login == "" {
		payer = authenticatedUser
	}
	newExpense.PayerLogin = payer.Login
	expense, err := services.GetDataAccess().GetExpenseDataAccess().CreateExpense(ctx, newExpense)
	if err != nil {
		return nil
	} else {
		return &expense
	}
}

func validateExpense(expense *models.Expense) error {
	errorMessages := make([]string, 0)
	if expense.Amount == 0.0 {
		errorMessages = append(errorMessages, "Invalid or missing 'Amount' field")
	}
	if expense.Date.IsZero() {
		errorMessages = append(errorMessages, "Missing 'Date' field")
	}
	if expense.Recipient == "" {
		errorMessages = append(errorMessages, "Missing 'Recipient' field")
	}
	if expense.Description == "" {
		errorMessages = append(errorMessages, "Missing 'Description' field")
	}
	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, "/"))
	}
	return nil
}

func mergeExpense(old, new *models.Expense) {
	if new.Amount != 0.0 {
		old.Amount = new.Amount
	}
	if !new.Date.IsZero() {
		old.Date = new.Date
	}
	if new.Recipient != "" {
		old.Recipient = new.Recipient
	}
	if new.Description != "" {
		old.Description = new.Description
	}
}

func UpdateExpense(ctx context.Context, authenticatedUser models.User, update models.Expense, patch bool) (models.Expense, error) {
	var updatedExpense models.Expense
	if update.Id == "" {
		return updatedExpense, errors.New("No Expense ID provided")
	}
	expense := GetExpenseInfo(ctx, update.Id)
	if expense == nil {
		return updatedExpense, errors.New("Expense not found")
	}
	if authenticatedUser.Login != expense.PayerLogin && !authenticatedUser.IsAdmin {
		return updatedExpense, errors.New("Not the right do update this expense")
	}
	update.PayerLogin = expense.PayerLogin
	if !patch {
		if err := validateExpense(&update); err != nil {
			return updatedExpense, err
		}
		updatedExpense = update
	} else {
		mergeExpense(expense, &update)
		updatedExpense = *expense
	}
	err := services.GetDataAccess().GetExpenseDataAccess().UpdateExpense(ctx, updatedExpense)
	return updatedExpense, err
}

func DeleteExpense(ctx context.Context, authenticatedUser models.User, id string) error {
	expense := GetExpenseInfo(ctx, id)
	if expense == nil {
		return errors.New("Expense not found")
	}
	if authenticatedUser.Login != expense.PayerLogin && !authenticatedUser.IsAdmin {
		return errors.New("Not the right do delete this expense")
	}
	return services.GetDataAccess().GetExpenseDataAccess().DeleteExpense(ctx, id)
}

func GetExpenseInfo(ctx context.Context, id string) *models.Expense {
	return services.GetDataAccess().GetExpenseDataAccess().RetrieveExpenseFromId(ctx, id)
}
