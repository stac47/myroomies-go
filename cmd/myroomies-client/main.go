package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/table"
	//"github.com/jedib0t/go-pretty/text"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stac47/myroomies/pkg/libmyroomies"
	"github.com/stac47/myroomies/pkg/models"
)

type GlobalSettings struct {
	serverUrl string
	user      string
	password  string
	debug     bool
}

type commonParams struct {
	date        string
	recipient   string
	description string
	amount      string
	id          int
}

var params commonParams

var settings GlobalSettings

func init() {
	cobra.OnInitialize(initConfig)

	// List command
	expenseCmd.AddCommand(listCmd)

	// Create expense command
	createCmd.Flags().StringVar(&params.recipient,
		"recipient",
		"",
		"Recipient of the expense")
	createCmd.MarkFlagRequired("recipient")
	createCmd.Flags().StringVar(&params.date,
		"date",
		"",
		"Date of the expense")
	createCmd.MarkFlagRequired("date")
	createCmd.Flags().StringVar(&params.description,
		"description",
		"",
		"Description of the expense")
	createCmd.MarkFlagRequired("description")
	createCmd.Flags().StringVar(&params.amount,
		"amount",
		"",
		"Amount of the expense")
	createCmd.MarkFlagRequired("amount")
	expenseCmd.AddCommand(createCmd)

	// Update expense command
	updateCmd.Flags().StringVar(&params.recipient,
		"recipient",
		"",
		"Recipient of the expense")
	updateCmd.Flags().StringVar(&params.date,
		"date",
		"",
		"Date of the expense")
	updateCmd.Flags().StringVar(&params.description,
		"description",
		"",
		"Description of the expense")
	updateCmd.Flags().StringVar(&params.amount,
		"amount",
		"",
		"Amount of the expense")
	expenseCmd.AddCommand(updateCmd)

	// Delete expense command
	expenseCmd.AddCommand(deleteCmd)

	rootCmd.PersistentFlags().String("host",
		"http://fsgtcyclisme06.fr:8080",
		"Server on which MyRoomies is running")
	rootCmd.PersistentFlags().String("user",
		"",
		"Current username")
	rootCmd.PersistentFlags().String("password",
		"",
		"Current user password")
	rootCmd.PersistentFlags().Bool("debug",
		false,
		"Run the client in debug mode")
	viper.BindPFlag("myroomies.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("myroomies.login", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("myroomies.password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("myroomies.debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(expenseCmd)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Cannot find the $HOME directory: %s\n", err)
		os.Exit(1)
	}
	viper.SetConfigName("myroomiesrc")
	viper.SetConfigType("toml")
	viper.AddConfigPath(home)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		logrus.Errorf("Error reading config. Reason: %s", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "myroomies",
	Short: "A client to access the myroomies server",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var expenseCmd = &cobra.Command{
	Use:   "expense",
	Short: "Manage the houseshare expenses",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List the house share expenses",
	RunE:  listCmdRun,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Record a new expense",
	RunE:  createCmdRun,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing expense",
	Args:  cobra.ExactArgs(1),
	RunE:  updateCmdRun,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing expense",
	Args:  cobra.MinimumNArgs(1),
	RunE:  deleteCmdRun,
}

func deleteCmdRun(cmd *cobra.Command, args []string) error {
	ids := make([]string, 0)
	for _, id := range args {
		ids = append(ids, id)
	}
	cli, err := getClient()
	if err != nil {
		return err
	}
	for _, id := range ids {
		err := cli.ExpenseDelete(id)
		if err != nil {
			fmt.Printf("WARNING: Could not delete expense [id=%s]: %s\n", id, err)
		}
	}
	return nil
}

func updateCmdRun(cmd *cobra.Command, args []string) error {
	id := args[0]
	expense := &models.Expense{}
	if params.date != "" {
		t, err := time.Parse("2006-01-02", params.date)
		if err != nil {
			return err
		}
		expense.Date = t
	}
	if params.recipient != "" {
		expense.Recipient = params.recipient
	}
	if params.description != "" {
		expense.Description = params.description
	}
	if params.amount != "" {
		amount, err := strconv.ParseFloat(params.amount, 64)
		if err != nil {
			return err
		}
		expense.Amount = amount
	}
	cli, err := getClient()
	if err != nil {
		return err
	}
	if err := cli.ExpenseUpdate(id, expense); err != nil {
		return err
	}
	fmt.Printf("Updated expense: %s\n", expense.Id)
	return nil
}

func createCmdRun(cmd *cobra.Command, args []string) error {
	t, err := time.Parse("2006-01-02", params.date)
	if err != nil {
		return err
	}
	amount, err := strconv.ParseFloat(params.amount, 64)
	if err != nil {
		return err
	}
	expense := models.Expense{
		Recipient:   params.recipient,
		Date:        t,
		Description: params.description,
		Amount:      amount,
	}
	cli, err := getClient()
	if err != nil {
		return err
	}
	id, err := cli.ExpenseCreate(expense)
	if err != nil {
		return err
	}
	fmt.Printf("Created expense: %s\n", id)

	return nil
}

func listCmdRun(cmd *cobra.Command, args []string) error {
	cli, err := getClient()
	if err != nil {
		return err
	}

	expenses, err := cli.ExpenseList()
	if err != nil {
		return err
	}

	roomiesMap := make(map[string]float64)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Date", "Recipient", "Description", "Roomie", "Amount"})
	spentTotal := 0.0
	for _, expense := range expenses {
		roomiesMap[expense.PayerLogin] += expense.Amount
		spentTotal += expense.Amount

		t.AppendRow(table.Row{
			expense.Id,
			expense.Date.Format("2006-01-02"),
			expense.Recipient,
			expense.Description,
			expense.PayerLogin,
			strconv.FormatFloat(expense.Amount, 'f', 2, 64),
		})
	}
	for k, v := range roomiesMap {
		balance := v - spentTotal/float64(len(roomiesMap))
		t.AppendFooter(table.Row{"", "", "", "Total", k,
			fmt.Sprintf("%.2f (%+.2f)", v, balance)})
	}

	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

func getClient() (libmyroomies.Client, error) {
	return libmyroomies.NewClient(
		viper.GetString("myroomies.host"),
		viper.GetString("myroomies.login"),
		viper.GetString("myroomies.password"),
		viper.GetBool("myroomies.debug"))
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	Execute()
}
