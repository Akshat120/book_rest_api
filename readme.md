# BooksAPI Setup Guide

## Prerequisites

- [Homebrew](https://brew.sh/) (for macOS)
- [dbmate](https://github.com/amacneil/dbmate) (for database migrations)
- [MySQL](https://dev.mysql.com/downloads/mysql/)

## 1. Install Dependencies

```sh
brew install dbmate mysql
```

## 2. Configure Database Connection

Set your database URL as an environment variable (replace placeholders):

```sh
export DATABASE_URL="mysql://<username>:<password>@127.0.0.1:3306/<database_name>"
```

## 3. Database Migration Commands

- **Create the database (if not already created):**
  ```sh
  dbmate create
  ```
- **Create a new migration file:**
  ```sh
  dbmate new <migration_file_name>
  ```
  > Migration files are created in `./db/migration/`

- **Apply all migrations:**
  ```sh
  dbmate up
  ```

- **Seed the database:**
  ```sh
  mysql -u root -p <database_name> < ./db/seeds/books.sql
  ```

## 4. Application Configuration

Create a `config.json` file in the project root with the following structure:

```json
{
  "app_name": "BooksApi",
  "server": {
    "port": 8080,
    "read_timeout": 15,
    "write_timeout": 20,
    "idle_timeout": 60
  },
  "basic_auth": {
    "username": "username",
    "password": "password"
  },
  "database": {
    "type": "mysql",
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "password",
    "db_name": "my_book"
  }
}
```

> **Note:**  
> - Adjust the values as needed for your environment.
> - All fields are required unless your code sets defaults.

---

## 5. Running the Application

Build and run your Go application as usual:

```sh
go run main.go
```

---

## Additional Notes

- Use `dbmate help` for more migration commands.
- Ensure your MySQL server is running before running migrations or seeding.
- For any issues, check your `config.json` and environment variables.

---






