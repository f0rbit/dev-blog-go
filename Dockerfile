# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port that your server will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
