package expensemngt

import (
	"context"
	"errors"
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
	if authenticatedUser.IsAdmin &&
		newExpense.PayerLogin != "" {

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
