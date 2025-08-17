
# Car Dealership App

A full-stack web application for managing a car dealership, featuring a React frontend and a Golang backend with MongoDB. The app allows users to view, add, update, reserve, and sell cars, with image upload and customer management.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Docker Usage](#docker-usage)
- [Contributing](#contributing)
- [Attribution](#attribution)

---

## Features

- **Home Page:** Overview of the dealership.
- **Car Management:** Add, update, delete, reserve, and sell cars.
- **Image Upload:** Store car images using GridFS in MongoDB.
- **Customer Management:** Reserve and sell cars to customers.
- **Status Tracking:** Cars can be available, reserved, or sold.
- **404 Page:** Custom not found page.

---

## Tech Stack

- **Frontend:** React, React Router DOM, Axios, HTML5, CSS3
- **Backend:** Golang, Gorilla Mux, MongoDB, GridFS
- **Containerization:** Docker, Docker Compose

---

## Project Structure

```
Car-Dealership/
   backend/      # Golang API server
   frontend/     # React client app
   docker-compose.yml
   Car-Dealership.json  # Postman API collection
   README.md
```

---

## Getting Started

### Prerequisites

- Node.js (>= 12)
- npm
- Go (>= 1.22)
- Docker (optional, for containerized setup)

### Local Development

**Frontend:**
```bash
cd frontend
npm install
npm start
```
Set `REACT_APP_API_URL` in `.env` to your backend URL (default: `http://localhost:8000`).

**Backend:**
```bash
cd backend
go run main.go
```
Set `MONGO_URI` in `.env` (default: `mongodb://localhost:27017/carDealershipDB`).

### Docker Compose

To run the entire stack with MongoDB:
```bash
docker-compose up --build
```
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend API: [http://localhost:8000](http://localhost:8000)
- MongoDB: [localhost:27017](localhost:27017)

---

## API Endpoints

See `Car-Dealership.json` for a full Postman collection.

**Main Endpoints:**
- `GET /cars/available` — List available cars
- `GET /cars/reserved` — List reserved cars
- `GET /cars/sold` — List sold cars
- `POST /cars` — Add a new car (multipart/form-data)
- `PUT /cars/{id}` — Update car details
- `DELETE /cars/{id}` — Delete a car
- `POST /cars/{id}/reserve` — Reserve a car
- `POST /cars/{id}/cancel-reservation` — Cancel reservation
- `POST /cars/{id}/sell` — Sell a car
- `GET /cars/image/{pictureID}` — Get car image

---

## Docker Usage

- **Build and run all services:** `docker-compose up --build`
- **Stop services:** `docker-compose down`
- **MongoDB data is persisted in the `mongo-data` volume.**

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss your ideas.

---

## Attribution

Favicon from [IconArchive](https://www.iconarchive.com/show/ionicons-icons-by-ionic/car-sport-sharp-icon.html), MIT license.