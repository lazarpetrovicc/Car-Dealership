
# Car Dealership App

A full-stack web application for managing a car dealership. The project combines a React frontend, a Go backend, and MongoDB to support inventory browsing, car creation and updates, image upload, reservation workflows, and sales tracking.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Testing](#testing)
- [API Endpoints](#api-endpoints)
- [Docker Usage](#docker-usage)
- [Contributing](#contributing)
- [Attribution](#attribution)

---

## Features

- **Inventory management:** Create, update, and delete car listings.
- **Reservation and sales flow:** Reserve cars and mark them as sold.
- **Image upload:** Store and serve car images with MongoDB GridFS.
- **Status tracking:** Keep cars in available, reserved, or sold states.
- **Modern UI:** Navigate the app with React Router and a responsive frontend experience.
- **Helpful UX:** Includes a custom 404 page for invalid routes.

---

## Tech Stack

- **Frontend:** React, React Router DOM, Axios, HTML5, CSS3
- **Backend:** Go, Gorilla Mux, MongoDB, GridFS
- **Testing:** Jest, React Testing Library, Go test
- **Containerization:** Docker, Docker Compose

---

## Project Structure

```
Car-Dealership/
   backend/      # Go API server
   frontend/     # React client app
   docs/         # OpenAPI + Postman docs
   docker-compose.yml
   README.md
```

---

## Getting Started

### Prerequisites

- Node.js (>= 12)
- npm
- Go (>= 1.22)
- Docker Desktop or Docker Engine (recommended for the database and full-stack setup)

### Quick Start with Docker

The easiest way to run the full application is with Docker Compose. This starts MongoDB, the backend API, and the frontend together.

```bash
docker compose up --build
```

Once the containers are running:

- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend API: [http://localhost:8000](http://localhost:8000)
- MongoDB: [localhost:27017](localhost:27017)

To stop everything:

```bash
docker compose down
```

### Local Development

If you prefer to run the services directly on your machine, start by making sure MongoDB is available.

#### Backend

```bash
cd backend
copy nul .env
go run main.go
```

The backend reads `MONGO_URI` from the environment or from a local `.env` file in the backend folder.

Example:

```env
MONGO_URI=mongodb://localhost:27017/carDealershipDB
```

#### Frontend

```bash
cd frontend
npm install
npm start
```

The frontend uses the backend URL from `REACT_APP_API_URL`. You can define it in a `.env` file inside the frontend folder:

```env
REACT_APP_API_URL=http://localhost:8000
```

> If the frontend cannot connect to the backend, verify that the backend is running on port 8000 and that the API URL in the frontend environment matches it.

---

## Environment Variables

### Backend

- `MONGO_URI` — MongoDB connection string.
  - Example: `mongodb://localhost:27017/carDealershipDB`

### Frontend

- `REACT_APP_API_URL` — Base URL for the API server.
  - Default: `http://localhost:8000`

---

## Testing

### Local tests

> Backend tests require a running MongoDB instance. If you do not already have MongoDB installed locally, start the database first with Docker:
>
> ```bash
> docker compose up -d mongo
> ```

#### Backend

```bash
cd backend
go test ./tests/...
```

#### Frontend

```bash
cd frontend
npm test
```

### Docker tests

The repository includes Docker Compose test services for backend and frontend.

```bash
docker compose up --build --abort-on-container-exit test-backend test-frontend
```

Or run individual test services:

```bash
docker compose run --rm test-backend
docker compose run --rm test-frontend
```

---

## API Endpoints

The backend exposes car-related routes through the router in the backend. API documentation is available in the docs folder:

- OpenAPI spec: [docs/openapi.yaml](docs/openapi.yaml)
- Postman collection: [docs/Car-Dealership.json](docs/Car-Dealership.json)

### Car listing and management

- `GET /cars/{status}` — List cars by status, where `status` is one of `available`, `reserved`, or `sold`
- `POST /cars` — Create a new car (multipart/form-data)
- `PUT /cars/{id}` — Update an existing car
- `DELETE /cars/{id}` — Remove a car from the database

### Reservation and sales actions

- `POST /cars/{id}/reserve` — Reserve a specific car
- `POST /cars/{id}/cancel-reservation` — Cancel an existing reservation
- `POST /cars/{id}/sell` — Mark a car as sold

### Images

- `GET /cars/image/{id}` — Retrieve the image associated with a car

---

## Docker Usage

- **Build and run all services:** `docker compose up --build`
- **Stop services:** `docker compose down`
- **Persisted data:** MongoDB storage is kept in the `mongo-data` volume
- **Test services:** `docker compose up --build --abort-on-container-exit test-backend test-frontend`

---
## API Documentation

You can find the API reference files in [docs/openapi.yaml](docs/openapi.yaml) and [docs/Car-Dealership.json](docs/Car-Dealership.json).
## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss your ideas.

---

## Attribution

Favicon from [IconArchive](https://www.iconarchive.com/show/ionicons-icons-by-ionic/car-sport-sharp-icon.html), MIT license.