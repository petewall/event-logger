FROM golang AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go /app/
ADD cmd /app/cmd
ADD internal /app/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /event-logger

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /event-logger /event-logger
ADD web /web

EXPOSE 8081

USER nonroot:nonroot

VOLUME [ "/data" ]

ENTRYPOINT ["/event-logger"]
