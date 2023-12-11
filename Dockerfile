# Step 1: Build Go Backend
FROM golang:1.21-alpine AS build-dm
WORKDIR /build

# Copy source files
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod .
COPY go.sum .
# Get go packages
RUN go mod download
# Build
RUN go build -o ./bin/kikuri ./cmd/kikuri/main.go
# Step 2: Build React Frontend
#FROM node:18-alpine AS build-fe
#WORKDIR /build

# Copy web files
#COPY web .

# Install dependencies
#RUN npm install

# Build
#RUN npm run build --base=/ --outDir=dist

# Step 3: Build runtime
# Set up runtime
FROM alpine:3 AS runtime
WORKDIR /app

# Copy Kikuri
COPY --from=build-dm /build/bin/kikuri /app/kikuri
#COPY --from=build-fe /build/dist web/dist

# Prepare config directory
RUN mkdir -p /etc/kikuri

# not sure if this is needed
EXPOSE 8080
ENTRYPOINT ["/app/kikuri"]
CMD ["-c", "/etc/kikuri/config.toml"]