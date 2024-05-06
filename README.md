# CRUD Users with Login System (Local Auth) using Gin

This project is a comprehensive user management system built with Gin, a web framework for Go. It includes features such as login, registration, logout, and CRUD operations for user data. It implements authorization for certain routes using JSON Web Tokens (JWT), providing secure access to protected endpoints.

When a user successfully logs in, a JWT is generated and set in the cookies, with an expiration time of 1 day. This ensures that authenticated users can access protected routes without having to repeatedly log in.

User data management involves hashing passwords for security purposes, ensuring that sensitive information is securely stored in the database. The project uses SQLite as the database engine and leverages the ORM framework GORM for database operations, simplifying database interactions and improving code readability.

To ensure code quality and reliability, the project includes comprehensive testing using the Ginkgo and Gomega testing frameworks. This helps identify and fix bugs early in the development process, ensuring a robust and stable user management system.

# Project Structure
## project-starter
Use the following command to create the basic project structure:
```sh
mkdir -p cmd/server internal/handlers internal/models internal/repository internal/server db/migrations config
touch cmd/server/main.go \
      internal/handlers/user_handler.go \
      internal/models/user.go \
      internal/repository/user_repository.go \
      internal/server/server.go \
      go.mod
```

> This command creates the basic structure for the project, including directories for server code, handlers, models, repository, server setup, database migrations, and configuration files.

## Define Models
Define the models needed for the application in the internal/models directory. This typically includes a User model to represent user data.

## Define Repository
Define the repository layer responsible for interacting with the database in the internal/repository directory. Implement functions to perform CRUD operations on user data.
## Define Handlers
Define the HTTP request handlers for the application in the internal/handlers directory. These handlers will handle incoming HTTP requests, validate input, and call the appropriate repository functions.
## Define Server
 the server setup and configuration in the internal/server directory. This includes setting up the HTTP server, defining routes, middleware, and any other server-related configurations.
## Define Configuration
Define configuration files for the application in the config directory. This may include environment-specific configuration files (e.g., development.yaml, production.yaml) for managing different settings in different environments.
## Define Migrations
Define database migration files in the db/migrations directory. These files contain SQL statements to create, modify, or delete database schema objects.
## Define Tests
Write tests in the internal directory, alongside models, handlers, and repository code, to ensure that your application functions correctly. Use testing frameworks like Go's built-in testing package or external packages like Ginkgo and Gomega for behavior-driven development.

# Endpoints Documentation

#### No Authentication Required:

- **Welcome**
  - **Description:** Get a welcome message.
  - **Method:** GET
  - **Endpoint:** `/`
  - **Handler:** `userHandler.Welcome`

- **User Login**
  - **Description:** Authenticate user and generate JWT token.
  - **Method:** POST
  - **Endpoint:** `/users/login`
  - **Handler:** `userHandler.UserLogin`

- **User Registration**
  - **Description:** Register a new user.
  - **Method:** POST
  - **Endpoint:** `/users/register`
  - **Handler:** `userHandler.UserRegister`

- **User Logout**
  - **Description:** Logout the currently authenticated user.
  - **Method:** GET
  - **Endpoint:** `/users/logout`
  - **Handler:** `userHandler.UserLogout`

#### Authentication Required:

- **Get All Users**
  - **Description:** Get a list of all users.
  - **Method:** GET
  - **Endpoint:** `/users`
  - **Handler:** `userHandler.GetAllUsers`
  - **Authentication:** Required (JWT token in cookies)

- **Get User by ID**
  - **Description:** Get details of a specific user by their ID.
  - **Method:** GET
  - **Endpoint:** `/users/{id}`
  - **Handler:** `userHandler.GetUserByID`
  - **Authentication:** Required (JWT token in cookies)

- **Create User**
  - **Description:** Create a new user.
  - **Method:** POST
  - **Endpoint:** `/users`
  - **Handler:** `userHandler.CreateUser`
  - **Authentication:** Required (JWT token in cookies)

- **Update User**
  - **Description:** Update details of a specific user.
  - **Method:** PUT
  - **Endpoint:** `/users/{id}`
  - **Handler:** `userHandler.UpdateUser`
  - **Authentication:** Required (JWT token in cookies)

- **Delete User**
  - **Description:** Delete a specific user by their ID.
  - **Method:** DELETE
  - **Endpoint:** `/users/{id}`
  - **Handler:** `userHandler.DeleteUser`
  - **Authentication:** Required (JWT token in cookies)

These endpoints provide functionalities for user authentication, registration, data management, and access control. The authentication middleware ensures that certain routes are only accessible to authenticated users with valid JWT tokens.