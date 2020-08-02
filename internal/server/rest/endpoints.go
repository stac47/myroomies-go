package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/stac47/myroomies/internal/server/data"
	"github.com/stac47/myroomies/internal/server/services"
	"github.com/stac47/myroomies/internal/server/services/usermngt"
	"github.com/stac47/myroomies/pkg/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type AccessRight int

const (
	NoneRight AccessRight = 1 << iota
	AuthenticatedRight
	AdminRight
)

const (
	EnvRootLogin    = "MYROOMIES_ROOT_LOGIN"
	EnvRootPassword = "MYROOMIES_ROOT_PASSWORD"
)

func hasRights(accessRights AccessRight, f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accessRights != NoneRight {
			user, ok := GetAuthenticatedUser(r.Context())
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if accessRights&AdminRight != 0 && !user.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		f(w, r)
	})
}

func registerEndpoints() {
	router := mux.NewRouter()

	// Add your handlers hereafter
	// Technical handlers
	router.HandleFunc("/version", version).Methods("GET")
	router.HandleFunc("/echo/{key:\\d+}", echo).Methods("GET")
	router.HandleFunc("/reset",
		hasRights(AdminRight, resetServer)).Methods("POST")

	// User management handlers
	router.HandleFunc("/users",
		hasRights(AuthenticatedRight, retrieveUsers)).Methods("GET")
	router.HandleFunc("/users",
		hasRights(AdminRight, createUser)).Methods("POST")
	router.HandleFunc("/users/{login:[a-zA-Z0-9]+}",
		hasRights(AuthenticatedRight, getUserInfo)).Methods("GET")
	router.HandleFunc("/users/{login:[a-zA-Z0-9]+}",
		hasRights(AuthenticatedRight, updateUser)).Methods("PUT")
	router.HandleFunc("/users/{login:[a-zA-Z0-9]+}",
		hasRights(AdminRight, deleteUser)).Methods("DELETE")

	// Expense handlers
	router.HandleFunc("/expenses",
		hasRights(AuthenticatedRight, retrieveAllExpenses)).Methods("GET")
	router.HandleFunc("/expenses",
		hasRights(AuthenticatedRight, createExpense)).Methods("POST")
	router.HandleFunc("/expenses/{id:[a-z0-9]+}",
		hasRights(AuthenticatedRight, getExpenseInfo)).Methods("GET")
	router.HandleFunc("/expenses/{id:[a-z0-9]+}",
		hasRights(AuthenticatedRight, deleteExpense)).Methods("DELETE")

	var authentication authenticationMiddleware
	router.Use(authentication.Middleware)

	http.Handle("/", router)
}

func createAdminOnFirstStart() error {
	ctx := context.Background()
	if len(usermngt.GetUsersList(ctx)) == 0 {
		rootLogin := os.Getenv(EnvRootLogin)
		if rootLogin == "" {
			rootLogin = "root"
		}
		rootPassword := os.Getenv(EnvRootPassword)
		if rootPassword == "" {
			return errors.New(fmt.Sprintf("On first start, a root password must "+
				"be given through environment variable %s", EnvRootPassword))
		}

		log.Infof("As it is the first start of MyRoomies server, a 'root' account "+
			"is being created. Login: %s", rootLogin)
		err := usermngt.CreateUser(ctx, models.User{
			Firstname: "root",
			Lastname:  "root",
			IsAdmin:   true,
			Login:     rootLogin,
			Password:  models.PasswordType(rootPassword),
		})
		if err != nil {
			return err
		}
		log.Info("'root' account created")
	}
	return nil
}

type ServerConfig struct {
	// The URL of the the datastore like 'mongodb://localhost:27017'. If empty,
	// or not match the MongoDB URL, the datastore will be the memory (data not
	// persisted).
	Storage string

	// The address to bind the listening socket to. Example: 'localhost:8080'
	// or ':8080'.
	BindTo string

	// The location of the certificate or the certificates chain.
	CertificatePath string

	// The location of the private key file corresponding to the given
	// certificate.
	KeyPath string
}

func registerSignalHandlers() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-c
		log.Infof("Received signal (%s). The server will shutdown.", s)
		services.Shutdown()
		os.Exit(0)
	}()
}

func validateConfiguration(config ServerConfig) error {
	if config.CertificatePath != "" && config.KeyPath == "" {
		return errors.New("A path to a certificate was given but no path to a private key.")
	}
	if config.CertificatePath == "" && config.KeyPath != "" {
		return errors.New("A path to a private key was given but no path to a certificate.")
	}

	return nil
}

func Start(config ServerConfig) (err error) {
	if err = validateConfiguration(config); err != nil {
		return
	}
	registerSignalHandlers()
	if strings.Contains(config.Storage, "mongodb://") {
		services.Configure(data.MongoDataAccessParams{Server: config.Storage})
	} else {
		// Non-persistent storage (memory) generally for tests
		log.Warn("No persistent data store selected: the data will be lost on " +
			"server shutdown.")
		services.Configure(nil)
	}
	registerEndpoints()
	if err = createAdminOnFirstStart(); err != nil {
		return
	}
	if config.CertificatePath == "" && config.KeyPath == "" {
		log.Warn("No certificate and no private key provided. The server " +
			"will run in plain HTTP: this is unsecured.")
		log.Printf("Listening on [%s] (HTTPS not activated).", config.BindTo)
		err = http.ListenAndServe(config.BindTo, nil)
	} else {
		log.Printf("Listening on [%s] (HTTPS activated).", config.BindTo)
		err = http.ListenAndServeTLS(config.BindTo,
			config.CertificatePath,
			config.KeyPath,
			nil)
	}
	return
}
