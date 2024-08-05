# syntax=docker/dockerfile:1

# Build stage
FROM node:18-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code and build the React app
COPY . .
RUN npm run build

# Test stage
FROM node:18-alpine AS tester
WORKDIR /app

# Copy only the necessary files for testing
COPY package.json package-lock.json ./
RUN npm install
COPY . .

# Run tests
CMD ["npm", "test"]

# Serve stage
FROM node:18-alpine AS runner
WORKDIR /app

# Install `serve` globally
RUN npm install -g serve

# Copy the build artifacts from the builder stage
COPY --from=builder /app/build /app/build

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the application
CMD ["serve", "-s", "build"]