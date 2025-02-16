FROM golang:1.23

WORKDIR /app

# copy dependency files
COPY go.mod .
COPY go.sum .

# download dependencies
RUN go mod tidy

# copy source code and .env
COPY . .
COPY .env .

# build the binary with the correct name
RUN go build -o robin-playlist ./cmd

# make the binary executable (optional, but ensure the name matches)
RUN chmod +x robin-playlist

# expose port
EXPOSE 9002
EXPOSE 7003

# run the binary
CMD ["./robin-playlist"]