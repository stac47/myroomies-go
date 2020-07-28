package rest

import (
	"net/http"

	"github.com/stac47/myroomies/internal/server/services/expensemngt"
	"github.com/stac47/myroomies/internal/server/services/usermngt"

	log "github.com/sirupsen/logrus"
)

func resetServer(w http.ResponseWriter, r *http.Request) {
	log.Info("Resetting the server")
	authenticatedUser, _ := GetAuthenticatedUser(r.Context())

	// Removing all users except the logged user
	users := usermngt.GetUsersList(r.Context())
	for _, user := range users {
		if user != authenticatedUser {
			usermngt.DeleteUser(r.Context(), user.Login)
		}
	}

	// Removing expenses
	expenses := expensemngt.GetExpensesList(r.Context(), expensemngt.ExpenseListOptions{})
	for _, expense := range expenses {
		expensemngt.DeleteExpense(r.Context(), authenticatedUser, expense.Id)
	}
}
