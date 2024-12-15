# Project Setup: Vue.js Frontend & Go Backend

## Prerequisites

- **Node.js & npm**: [Install Node.js](https://nodejs.org/)
- **Go**: [Install Go](https://golang.org/dl/)
- **MongoDB**: [Install MongoDB](https://www.mongodb.com/try/download/community)

 **Clone**:

**Backend**

git clone git@gitlab.com:tjasar/rurs-tvojkoticek-backend.git
cd backend
go mod init your-module-name
go mod tidy

go run main.go

**Frontend**

npm install -g @vue/cli
npm install
npm run serve

# Backend Explanation

- **Functions**: Contains all controllers and operations for database resources.
- **HTTP**: Holds all endpoint definitions.
- **Mongo**: Manages the database connection.
- **Schemas**: Defines the object parameters and data structures.

**main.go**: 
- Runs on port `8080`
- Configures CORS to allow requests only from `localhost` for the Vue.js frontend.

# Frontend Explanation

- **components**: Contains all components used in the website.
- **router.js**: Defines routes for navigation.
- **main.js**: Includes the router setup for rendering the application.






