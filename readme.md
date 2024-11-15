# Coffee-shop: Coffee Shop Management System

Hot-Coffee is a RESTful API application built in Go to manage a coffee shop's backend operations. It enables handling orders, managing menu items, tracking inventory, and running reports, all designed with a three-layered software architecture for maintainability and scalability.

## Features

- **Order Management**: Create, retrieve, update, delete, and close orders.
- **Menu Management**: Add, retrieve, update, and delete menu items.
- **Inventory Management**: Track ingredient stock levels, update quantities, and check availability for orders.
- **Reports**: Generate total sales and popular items reports.
- **Logging**: Integrated logging using `log/slog` for significant events and errors.
- **Error Handling**: Graceful error responses with appropriate HTTP status codes.

## Project Structure

- **cmd/**: Application entry point (`main.go`)
- **internal/**: Organized by layers
  - **handler/**: HTTP request handlers
  - **service/**: Business logic layer
  - **dal/**: Data Access Layer (repositories)
- **models/**: Data models for orders, menu items, and inventory
- **data/**: JSON files for persisting data (`orders.json`, `menu_items.json`, `inventory.json`)

## API Endpoints

### Orders

- `POST /orders` - Create a new order
- `GET /orders` - Retrieve all orders
- `GET /orders/{id}` - Retrieve order by ID
- `PUT /orders/{id}` - Update an order
- `DELETE /orders/{id}` - Delete an order
- `POST /orders/{id}/close` - Close an order

### Menu Items

- `POST /menu` - Add a new menu item
- `GET /menu` - Retrieve all menu items
- `GET /menu/{id}` - Retrieve a menu item by ID
- `PUT /menu/{id}` - Update a menu item
- `DELETE /menu/{id}` - Delete a menu item

### Inventory

- `POST /inventory` - Add a new inventory item
- `GET /inventory` - Retrieve all inventory items
- `GET /inventory/{id}` - Retrieve an inventory item by ID
- `PUT /inventory/{id}` - Update an inventory item
- `DELETE /inventory/{id}` - Delete an inventory item

### Reports

- `GET /reports/total-sales` - Retrieve total sales
- `GET /reports/popular-items` - Retrieve popular menu items

## Usage

```bash
./hot-coffee --port <N> --dir <data_directory>
./hot-coffee --help
