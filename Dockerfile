# deploying stage
FROM golang:1.18-alpine AS build

ENV GO111MODULE=on
# Updates the repository and installs git
RUN  apk update && apk upgrade
RUN  apk add --no-cache git

# Switches to /tmp/app as the working directory, similar to 'cd'
WORKDIR /tmp/app

## If you have a go.mod and go.sum file in your project, uncomment lines 13, 14, 15

# COPY go.mod .
# COPY go.sum .
# RUN go mod download

#Copy data project in image docker image
COPY . .

# Builds the current project to a binary file called api
# The location of the binary file is /tmp/app/out/urlshortener
RUN GOOS=linux go build -o ./out/urlshortener .

#RUN sudo chmod -R 777 /tmp/app
#########################################################

# The project has been successfully built and we will use a
# lightweight alpine image to run the server
FROM alpine:latest

# Adds CA Certificates to the image
#RUN sudo apk add ca-certificates

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /tmp/app/out/urlshortener /app/urlshortener
COPY .env /app/.env
# Switches working directory to /app
WORKDIR "/app"

# Exposes the 6067 port from the container
EXPOSE 6067

# Runs the binary once the container starts
CMD ["./urlshortener"]
