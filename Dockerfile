# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

COPY ./ ./

WORKDIR admin

# RUN go mod download
RUN go build -o /admin

WORKDIR ../ballot

RUN go build -o /ballot

WORKDIR ../registration

RUN go build -o /registration

EXPOSE 10000
EXPOSE 10001
EXPOSE 10002

CMD ["sh", "-c", "/admin & /ballot & /registration"]