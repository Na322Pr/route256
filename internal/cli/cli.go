package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
)

var (
	orderUC *usecase.OrderUseCase
	err     error
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
	Use:   "receive-courier",
	Short: "Receive order from courier",
	Long: `Usage: receive-courier orderID clientID storeUntil cost weight [package1] [package2]
Example: receive-courier 1 1 2024-10-01 15:20:00 1200 7 bag tape`,
	Run: func(cmd *cobra.Command, args []string) {
		minArgsCount := 6
		maxArgsCount := 8

		if len(args) < minArgsCount || len(args) > maxArgsCount {
			fmt.Println("Incorrect args count")
			return
		}

		var orderID, clientID, cost, weight int
		var storeUntil time.Time

		if orderID, err = strconv.Atoi(args[0]); err != nil {
			fmt.Println("orderID is incorrect")
		}

		if clientID, err = strconv.Atoi(args[1]); err != nil {
			fmt.Println("clientID is incorrect")
		}

		if storeUntil, err = time.Parse("2006-01-02 15:04:05", args[2]+" "+args[3]); err != nil {
			fmt.Println("storeUntil is incorrect")
		}

		cost, err = strconv.Atoi(args[4])
		if err != nil {
			fmt.Println("cost is incorrect")
		}

		weight, err := strconv.Atoi(args[5])
		if err != nil {
			fmt.Println("weight is incorrect")
		}

		packages := []string{"unknown", "unknown"}

		if len(args) >= 7 {
			packages[0] = args[6]
		}

		if len(args) == 8 {
			packages[1] = args[7]
		}

		req := dto.AddOrder{
			ID:         orderID,
			ClientID:   clientID,
			StoreUntil: storeUntil,
			Cost:       cost,
			Weight:     weight,
			Packages:   packages,
		}

		err = orderUC.ReceiveOrderFromCourier(context.Background(), req)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Order added successfully")
	},
}

var returnOrderToCourierCmd = &cobra.Command{
	Use:   "return-courier",
	Short: "Return order to courier",
	Long: `Usage: return-courier orderID
Example: return-courier 1`,
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

		err = orderUC.ReturnOrderToCourier(context.Background(), orderID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Order returned to courier successfully")
	},
}

var giveOutOrderToClientCmd = &cobra.Command{
	Use:   "give-out-client",
	Short: "Give out order to client",
	Long: `Usage: give-out-client [orderIDs...]
Example: give-out-client 1 2 3 4`,
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

		err := orderUC.GiveOrderToClient(context.Background(), orderIDs)
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
			return
		}

		// if len(args) == 2 {
		// 	lastCount, err := strconv.Atoi(args[1])
		// 	if err != nil {
		// 		fmt.Println("lastCount is incorrect")
		// 		return
		// 	}
		// }

		orders, err := orderUC.OrderList(context.Background(), clientID)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Order IDs list:")
		for i, order := range orders.Orders {
			fmt.Printf("%d:\t%d\n", i+1, order.ID)
		}

	},
}

var refundFromCustomerCmd = &cobra.Command{
	Use:   "refund-client",
	Short: "Refund order from client",
	Long: `Usage: refund-client clientID orderID
Example: refund-client 10 12`,
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

		err = orderUC.GetRefundFromСlient(context.Background(), clientID, orderID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Product has been successfully returned")
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
				return
			}
		}

		if len(args) == 2 {
			offset, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("offset is incorrect")
				return
			}
		}

		refunds, err := orderUC.RefundList(context.Background(), limit, offset)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Refund IDs list:")
		for i, order := range refunds.Orders {
			fmt.Printf("%d:\t%d\n", i+1, order.ID)
		}
	},
}

func Run(orderUseCase *usecase.OrderUseCase) {
	orderUC = orderUseCase
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
