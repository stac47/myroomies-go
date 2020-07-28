package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/stac47/myroomies/internal/server/services/expensemngt"
	"github.com/stac47/myroomies/pkg/models"
)

func retrieveAllExpenses(w http.ResponseWriter, r *http.Request) {
	listOptions := expensemngt.ExpenseListOptions{}
	expenses := expensemngt.GetExpensesList(r.Context(), listOptions)
	encodeJson(w, expenses)
}

func getExpenseInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if expense := expensemngt.GetExpenseInfo(r.Context(), id); expense != nil {
		encodeJson(w, expense)
	} else {
		msg := fmt.Sprintf("The user %s does not exist", id)
		http.Error(w, models.NewRestError(msg).ToJSON(),
			http.StatusNotFound)
	}
}

func createExpense(w http.ResponseWriter, r *http.Request) {
	var newExpense models.Expense
	if err := decodeJson(w, r.Body, &newExpense); err != nil {
		return
	}
	authenticatedUser, _ := GetAuthenticatedUser(r.Context())
	createdExpense := expensemngt.CreateExpense(r.Context(), authenticatedUser, newExpense)
	if createdExpense == nil {
		http.Error(w, models.NewRestError("Cannot create expense").ToJSON(),
			http.StatusNotFound)
		return
	}
	w.Header().Add("Location", fmt.Sprintf("/expenses/%s", createdExpense.Id))
	w.WriteHeader(http.StatusCreated)
}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	authenticatedUser, _ := GetAuthenticatedUser(r.Context())
	if err := expensemngt.DeleteExpense(r.Context(), authenticatedUser, id); err != nil {
		http.Error(w, models.NewRestError(err.Error()).ToJSON(),
			http.StatusNotFound)
	}
}
