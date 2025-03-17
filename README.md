# Go Web Service with Google Authentication

This project implements a Go web service with user authentication via Google OAuth 2.0. It uses the following key technologies:

* **Go:** The programming language for the backend service.
* **goth:** A multi-provider authentication package for Go.
* **gorilla/mux:** A powerful URL router and dispatcher for Go.
* **gorilla/sessions:** Secure cookie-based session management.
* **Prisma:** An ORM for database interactions.
* **PostgreSQL:** The database used for storing user data.

## Features

* User registration and login with Google.
* Secure session management.
* Database integration with Prisma.

## Getting Started

1. **Clone the repository:**

    ```bash
    git clone https://github.com/Akshay2642005/go-oauth2-service.git
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up environment variables:**

    Create a `.env` file in the root directory and set the following environment variables:

    ```
    PUBLIC_HOST=http://localhost
    PORT=8000
    DATABASE_URL="your_database_url"
    COOKIES_AUTH_SECRET="some-very-secret-key"
    GOOGLE_CLIENT_ID="your_google_client_id"
    GOOGLE_CLIENT_SECRET="your_google_client_secret"
    ```

4. **Run database migrations:**

    ```bash
    prisma migrate dev
    ```

5. **Run the server:**

    ```bash
    go run main.go
    ```

## API Endpoints

* `/auth/google`: Initiates the Google authentication flow.
* `/auth/google/callback`: Handles the callback from Google after authentication.
* `/logout`: Logs the user out.

## Contributing

Contributions are welcome! Please feel free to submit pull requests.

## License

This project is licensed under the MIT License.
