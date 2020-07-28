package data

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/stac47/myroomies/pkg/models"
)

var (
	// TODO: protect with a mutex
	memoryExpenseDataAccess *MemoryExpenseDataAccess
)

type MemoryExpenseDataAccess struct {
	expenses map[string]models.Expense
	counter  int
}

func GetMemoryExpenseDataAccess() *MemoryExpenseDataAccess {
	if memoryExpenseDataAccess == nil {
		memoryExpenseDataAccess = &MemoryExpenseDataAccess{
			expenses: make(map[string]models.Expense),
			counter:  0,
		}
	}
	return memoryExpenseDataAccess
}

func (dao *MemoryExpenseDataAccess) RetrieveExpenses(ctx context.Context) []models.Expense {
	result := make([]models.Expense, 0)
	for _, value := range dao.expenses {
		result = append(result, value)
	}
	return result
}

func (dao *MemoryExpenseDataAccess) CreateExpense(ctx context.Context, newExpense models.Expense) (models.Expense, error) {
	newExpense.Id = strconv.Itoa(dao.counter)
	dao.counter++
	dao.expenses[newExpense.Id] = newExpense
	return newExpense, nil
}

func (dao *MemoryExpenseDataAccess) RetrieveExpenseFromId(ctx context.Context, id string) *models.Expense {
	expense, ok := dao.expenses[id]
	if ok {
		return &expense
	}
	return nil
}

func (dao *MemoryExpenseDataAccess) UpdateExpense(ctx context.Context, updatedExpense models.Expense) (err error) {
	id := updatedExpense.Id
	_, ok := dao.expenses[id]
	if ok {
		dao.expenses[id] = updatedExpense
	} else {
		err = errors.New(fmt.Sprintf("Expense [id=%s] does not exist", id))
	}
	return
}

func (dao *MemoryExpenseDataAccess) DeleteExpense(ctx context.Context, id string) (err error) {
	_, ok := dao.expenses[id]
	if ok {
		delete(dao.expenses, id)
	} else {
		err = errors.New(fmt.Sprintf("Expense [id=%s] does not exist", id))
	}
	return
}
