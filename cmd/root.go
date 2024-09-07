package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/marchenkosasha2/homework/storage"
)

var (
	store *storage.Store
	err   error
)

var rootCmd = &cobra.Command{
	Use:   "homework",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var receiveOrderFromCourierCmd = &cobra.Command{
	Use:   "receive-order",
	Short: "Receive order from courier",
	Long: `Usage: receive-order orderID clientID storeUntil
Example: receive-order 1 1 2024-09-10 15:20:00`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println("Incorrect args count. Expected 3 arguments: orderID clientID storeUntil")
			return
		}

		orderID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("orderID is incorrect")
		}

		clientID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("clientID is incorrect")
		}

		storeUntil, err := time.Parse("2006-01-02 15:04:05", args[2]+" "+args[3])
		if err != nil {
			fmt.Println("storeUntil is incorrect")
		}

		err = store.GetOrderFromСourier(orderID, clientID, storeUntil)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Order added successfully")
	},
}

var returnOrderToCourierCmd = &cobra.Command{
	Use:   "return-order",
	Short: "Return order to courier",
	Long: `Usage: return-order orderID
Example: return-order 1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Incorrect args count. Expected 1 argument: orderID")
			return
		}

		orderID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("orderID is incorrect")
			return
		}

		err = store.GiveOrderToCourier(orderID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Order returned to courier successfully")
	},
}

var giveOutOrderToClientCmd = &cobra.Command{
	Use:   "give-out-order",
	Short: "Give out order to client",
	Long: `Usage: give-out-order [orderIDs...]
Example: give-out-order 1 2 3 4`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No arguments. Expected arguments: [orderIDs...]")
			return
		}

		var orderIDs []int

		for i := 0; i < len(args); i++ {
			orderID, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Println("One of orderIDs is incorrect")
				return
			}

			orderIDs = append(orderIDs, orderID)
		}

		err = store.GiveOrderToClient(orderIDs)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Orders have been successfully issued")
	},
}

var getOrderListCmd = &cobra.Command{
	Use:   "order-list",
	Short: "Get order list",
	Long: `Usage: order-list clientID [lastCount]
Example 1, return all orders:    order-list 10,
Example 2, return n last orders: order-list 10 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Println("Incorrect args count. Expected 1-2 arguments: clientID [lastCount]")
			return
		}

		clientID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("clientID is incorrect")
		}

		if len(args) == 2 {
			lastCount, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("lastCount is incorrect")
			}

			err = store.OrderList(clientID, lastCount)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		err = store.OrderList(clientID)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var refundFromCustomerCmd = &cobra.Command{
	Use:   "refund",
	Short: "Refund order",
	Long: `Usage: refund clientID orderID
Example: refund 10 12`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Incorrect args count. Expected 2 arguments: clientID orderID")
			return
		}

		clientID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("clientID is incorrect")
			return
		}

		orderID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("orderID is incorrect")
			return
		}

		err = store.GetRefundFromСlient(clientID, orderID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Refund successfully")
	},
}

var getRefundListCmd = &cobra.Command{
	Use:   "refund-list",
	Short: "Get refund list",
	Long: `Usage: refund-list [limit] [offset]
Example 1, return all refunds: 			 order-list 10,
Example 2, return n refunds: 			 order-list 10 10,
Example 3, return n refunds with offset: order-list 10 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			fmt.Println("Incorrect args count. Expected 2 or less arguments: [limit] [offset]")
			return
		}

		limit, offset := 0, 0

		if len(args) >= 1 {
			limit, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("limit is incorrect")
			}
		}

		if len(args) == 2 {
			offset, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("offset is incorrect")
			}
		}

		err = store.RefundList(limit, offset)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	store, err = storage.NewStore("storage/data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.AddCommand(receiveOrderFromCourierCmd)
	rootCmd.AddCommand(returnOrderToCourierCmd)
	rootCmd.AddCommand(giveOutOrderToClientCmd)
	rootCmd.AddCommand(getOrderListCmd)
	rootCmd.AddCommand(refundFromCustomerCmd)
	rootCmd.AddCommand(getRefundListCmd)

	HandleUserInput()
}

func HandleUserInput() {
	fmt.Println("Running App... Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("Exiting...")
				os.Exit(1)
			}

			commandArgs := strings.Fields(input)
			rootCmd.SetArgs(commandArgs)

			if err := rootCmd.Execute(); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}
