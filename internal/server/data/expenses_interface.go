package data

import (
	"context"

	"github.com/stac47/myroomies/pkg/models"
)

type ExpenseDataAccess interface {
	CreateExpense(ctx context.Context, newExpense models.Expense) (models.Expense, error)
	RetrieveExpenseFromId(ctx context.Context, id string) *models.Expense
	UpdateExpense(ctx context.Context, updatedExpense models.Expense) error
	DeleteExpense(ctx context.Context, id string) error

	RetrieveExpenses(ctx context.Context) []models.Expense
}
