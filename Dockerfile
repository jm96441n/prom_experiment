FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

RUN ls /go/bin

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
