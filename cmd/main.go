package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"hot-coffee/internal/dal"
)

// Define command-line flags for directory path and port number
var (
	dir  = flag.String("dir", "data", "Path to the directory")
	port = flag.String("port", "8080", "Port number")
)

func main() {
	flag.Parse()
	// Start the server with the specified port, handling any errors
	if *port == "0" {
		slog.Error("Port number must be between 1024 and 65535")
		os.Exit(1)
	}
	err := StartServer(*port)
	if err != nil {
		log.Fatal(err)
	}
}

// CreatedJSONfile ensures JSON files are created for orders, menu, and inventory data
func CreatedJSONfile() {
	directory := fmt.Sprintf("../%s", *dir)
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	orders := fmt.Sprintf("%s/%s", directory, dal.OrdersFile)
	_, err = os.Create(orders)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	menuPath := fmt.Sprintf("%s/%s", directory, dal.MenuItemFile)
	if check, err := dal.FileExistsInDirectory(menuPath); !check && err == nil {
		menu, err := os.Open(dal.ReserveMenu)
		if err != nil {

			slog.Error(err.Error())
			return
		}

		newMenu, err := os.Create(menuPath)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		_, err = io.Copy(newMenu, menu)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
	if err != nil {
		fmt.Println(8)
		return
	}
	inventoryPath := fmt.Sprintf("%s/%s", directory, dal.InventoryitemFile)
	if check, err := dal.FileExistsInDirectory(inventoryPath); !check && err == nil {
		newInventory, err := os.Create(inventoryPath)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		inventory, err := os.Open(dal.ReserveInventory)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		_, err = io.Copy(newInventory, inventory)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
	if err != nil {
		fmt.Println(9)
		return
	}
	order, err := os.Open(dal.ReserveOrder())
	if err != nil {

		slog.Error(err.Error())
		return
	}

	newOrder, err := os.Create(dal.Orders())
	if err != nil {
		slog.Error(err.Error())
		return
	}
	_, err = io.Copy(newOrder, order)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

// init sets up the usage information displayed when the help flag is used
func init() {
	flag.Usage = func() {
		fmt.Println(
			`Coffee Shop Management System

Usage:
	hot-coffee [--port <N>] [--dir <S>] 
	hot-coffee --help
			
Options:
	--help       Show this screen.
	--port N     Port number.
	--dir S      Path to the data directory.`)
	}
}
