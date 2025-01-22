# Thales Backend

This guide explains how to set up the `Thales` backend, including configuring PostgreSQL with Docker and setting up Goose for database migrations.

---

## Prerequisites

Ensure the following tools are installed on your system:

- Docker (Recommended for PostgreSQL setup)
- Go (Golang) version 1.23.2
- Git (for version control)

---

## Installation Steps

### Step 1: Set Up PostgreSQL Using Docker
Run the following command to set up and start a PostgreSQL database container:
1. Run the following command to download the PostgreSQL Docker image: 
```bash
docker pull postgres
```
2. Start a PostgreSQL container with the required environment variables:
```bash
docker run --name thales_db -e POSTGRES_USER=thales_user -e POSTGRES_PASSWORD=thales_pw -e POSTGRES_DB=thales -p 5432:5432 -d postgres
```
### Step 2: Set Up Goose ( migrations ) Environment Variables
Run the following command to set up and goose migration environment:
- ( Linux/macOS )
```bash
export GOOSE_DRIVER="postgres"
export GOOSE_DBSTRING="postgresql://thales_user:thales_pw@localhost:5432/thales"
export GOOSE_MIGRATION_DIR="./migrations"
```
- Window ( Powershell )
``` bash
$env:GOOSE_DRIVER = "postgres"
$env:GOOSE_DBSTRING = "postgresql://thales_user:thales_pw@localhost:5432/thales"
$env:GOOSE_MIGRATION_DIR = "./migrations"
```
### Step 3: Set Up the `.env` File

Create a `.env` file in the root directory of your project and add the following environment variables:

```bash
DB_USER="thales_user"
DB_PASSWORD="thales_pw"
DB_NAME="thales"
DB_CONNECTION="postgresql://thales_user:thales_pw@localhost:5432/thales?sslmode=disable"
```

### Step 4: Install Dependencies and Run the Application

1. **Install Dependencies:**  
   Run the following command to install the dependencies listed in your `go.mod` file:
   ```bash
   go get
   ```
2. **Clean Up Dependencies:**  
   Use `go mod tidy` to remove any unused dependencies and ensure everything is clean:
   ```bash
   go mod tidy
   ```
3. **Run the Application:**  
   Finally, run the application with the following command:
   ```bash
   go run main.go
   ```

### Step 5: Run Goose Migration
1. **Run Goose Migration**

    If the environment variables from Step 2 have been set up, you can run the migration by executing:
    ```bash
    goose up
    ```
    If the environment variables are not set, you can run the migration with the full connection string:
    ``` bash
    goose -dir ./migrations postgres "postgresql://thales_user:thales_pw@localhost:5432/thales?sslmode=disable" up
    ```
### Troubleshooting
- **PostgreSQL Connection Issues:** Ensure that your Docker container is running and the port is correctly mapped to `5432`.
- **Migration Errors:** If you encounter issues with Goose migrations, ensure that your migrations are correctly placed in the `./migrations` directory.

## API Reference

### API Testing with Postman
To quickly get started with testing the Thales API, you can import the provided Postman collection.
#### Steps to Import the Postman Collection
1. Download the `thales.postman_collection.json` file from the `/postman` directory in this repository.
2. Open Postman.
3. In the top left corner of Postman, click on the **Import** button.
4. In the Import dialog, select the **File** tab.
5. Click on **Choose Files** and select the `thales.postman_collection.json` file you downloaded.
6. Click **Open** to import the collection into Postman.

Once imported, you will see all available API requests organized by category, and you can start testing them by sending requests directly from Postman.

---

### API Documentations
#### Get all Products

```http
  GET /api/products
```
**Response**
- Success (200)
```bash
{
    "paginate": {
        "page": 1,
        "page_size": 2,
        "total_count": 2
    },
    "products": [
        {
            "id": 1,
            "name": "Product A",
            "type": "Type A",
            "picture_url": "uploads/product.png",
            "price": 72.4,
            "description": "Product A description",
            "created_at": "2025-01-22T16:03:51.69736Z",
            "updated_at": "2025-01-22T16:03:51.69736Z"
        }
    ]
}
```

---

#### Create Product

```http
  POST /api/products
```
**Request body**

The request body should be a raw JSON object with the following fields:

| Parameter      | Type     | Description                       |
| :--------      | :------- | :-------------------------------- |
| `name`         | `string` | **Required**. The name of the product |
| `type`         | `string` | **Required**. The type of the product |
| `price`        | `float`  | **Required**. The price of the product |
| `description`  | `string` | The description of the product |
| `image`        | `file`   | **Required**. The image file of the product. This should be a valid image file (e.g., `.jpg`, `.png`, etc.). |

**Response**
- Success (200)
```bash
{
    "product": {
        "id": 2,
        "name": "Product B",
        "type": "Type B",
        "picture_url": "uploads/product-b.png",
        "price": 100,
        "description": "Product B description",
        "created_at": "2025-01-22T16:31:33.791265Z",
        "updated_at": "2025-01-22T16:31:33.791265Z"
    }
}
```

---

#### Update Product by Id

```http
  PUT /api/products/{id}
```
**Request body**

The request body should be a raw JSON object with the following fields:

| Parameter      | Type     | Description                       |
| :--------      | :------- | :-------------------------------- |
| `name`         | `string` | **Required**. The name of the product |
| `type`         | `string` | **Required**. The type of the product |
| `price`        | `float`  | **Required**. The price of the product |
| `description`  | `string` | The description of the product |
| `image`        | `file`   | **Required**. The image file of the product. This should be a valid image file (e.g., `.jpg`, `.png`, etc.). |

**Response**
- Success (200)
```bash
{
    "product": {
        "id": 1,
        "name": "Product C",
        "type": "Type C",
        "picture_url": "uploads/product-c.png",
        "price": 80,
        "description": "Product C description",
        "created_at": "2025-01-22T16:36:32.627234Z",
        "updated_at": "2025-01-22T16:38:10.238455Z"
    }
}
```

---

#### Get Product by Id

```http
  GET /api/products/{id}
```
**Response**
- Success (200)
```bash
{
    "product": {
        "id": 1,
        "name": "Product B",
        "type": "Type B",
        "picture_url": "uploads/product-B.png",
        "price": 100,
        "description": "Product B description",
        "created_at": "2025-01-22T16:03:51.69736Z",
        "updated_at": "2025-01-22T16:36:27.08286Z"
    }
}
```

---

#### Delete Product by Id

```http
  DELETE /api/products/{id}
```
**Query Parameters**
| Parameter    | Type     | Description                       |
| :--------    | :------- | :-------------------------------- |
| `applicant`  | `string` | **Required.** The unique ID of the applicant.|

**Response**
- Success (200)
```bash
{
    "message": "Product deleted successfully"
}
```

---
