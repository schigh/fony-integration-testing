FROM golang:alpine AS builder

WORKDIR /
ADD ./ /
RUN apk add --no-cache git bash
RUN mkdir /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app .

FROM gcr.io/distroless/base
COPY --from=builder  /build/app /
EXPOSE 80
CMD ["/app"]
