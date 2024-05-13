# Eniqilo Store API

The eniqilo store product management API is a comprehensive solution designed to facilitate efficient management of products and streamline the checkout process for cashiers. This API provides a range of functionalities tailored to meet the needs of store staff, enabling them to register, log in, manage products, search SKUs, register customers, and handle checkout operations seamlessly.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development.

## MakeFile

build the application
```bash
make build
```

run the application
```bash
make run
```

live reload the application
```bash
make watch
```

clean up binary from the last build
```bash
make clean
```

## API

### Authentication

#### Register Staff
- **Method:** `POST`
- **Endpoint:** `/v1/staff/register`
- **Description:** Registers a new staff member.
- **Request Body:**
  - `phoneNumber` (string, required): The phone number of the staff member.
  - `name` (string, required): The name of the staff member.
  - `password` (string, required): The password of the staff member.
- **Response:** Returns staff details upon successful registration.

#### Staff Login
- **Method:** `POST`
- **Endpoint:** `/v1/staff/login`
- **Description:** Logs in an existing staff member.
- **Request Body:**
  - `phoneNumber` (string, required): The username of the staff member.
  - `password` (string, required): The password of the staff member.
- **Response:** Returns authentication token upon successful login.

### Product Management

#### Add Product
- **Method:** `POST`
- **Endpoint:** `/v1/product`
- **Description:** Adds a new product to the inventory.
- **Request Body:**
  - `name` (string, required): The name of the product.
  - `sku` (string, required): The sku of the product.
  - `category` (string, required): The category of the product.
  - `imageUrl` (string, required): The image of the product.
  - `notes` (string, required): The notes of the product.
  - `price` (integer, required): The price of the product.
  - `stock` (integer, required): The stock of the product.
  - `location` (string, required): The category of the product.
  - `isAvailable` (boolean, required): The category of the product.
- **Response:** Returns details of the added product.

#### Get Products
- **Method:** `GET`
- **Endpoint:** `/v1/product`
- **Description:** Retrieves all products from the inventory.
- **Response:** Returns a list of products.

#### Update Product
- **Method:** `PUT`
- **Endpoint:** `/v1/product/{id}`
- **Description:** Updates the details of a product in the inventory.
- **Request Body:** Same as Add Product.
- **Response:** Returns updated details of the product.

#### Delete Product
- **Method:** `DELETE`
- **Endpoint:** `/v1/product/{id}`
- **Description:** Deletes a product from the inventory.
- **Response:** Returns a success message upon successful deletion.

### Search SKU

#### Search Product by SKU
- **Method:** `GET`
- **Endpoint:** `/v1/product/customer`
- **Description:** Searches for a product in the inventory based on SKU (Stock Keeping Unit).
- **Response:** Returns details of the matching product.

### Checkout

#### Register Customer
- **Method:** `POST`
- **Endpoint:** `/v1/customer/register`
- **Description:** Registers a new customer.
- **Request Body:**
  - `name` (string, required): The name of the customer.
  - `phoneNumber` (string): The phone number of the customer.
- **Response:** Returns customer details upon successful registration.

#### Get Customers
- **Method:** `GET`
- **Endpoint:** `/v1/customer`
- **Description:** Retrieves all registered customers.
- **Response:** Returns a list of customers.

#### Product Checkout
- **Method:** `POST`
- **Endpoint:** `/v1/product/checkout`
- **Description:** Processes a product checkout for a customer.
- **Request Body:**
  - `customerId` (string, required): The ID of the customer making the purchase.
  - `productDetails` (array of products, required): 
	  - `productId` (string, required)
	  - `quantity` (integer, required)
  - `paid` (integer, required): The quantity of the product being purchased.
  - `change` (integer, required): The quantity of the product being purchased.
- **Response:** Returns details of the checkout transaction.

#### Get Checkout History
- **Method:** `GET`
- **Endpoint:** `/v1/product/checkout/history`
- **Description:** Retrieves the checkout history of products.
- **Response:** Returns a list of checkout transactions.
