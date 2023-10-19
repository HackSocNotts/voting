FROM golang:alpine as builder
ARG target

WORKDIR /build

ADD ${target}/go.mod ./${target}/go.mod
ADD ${target}/go.sum ./${target}/go.sum

ADD common/go.mod ./common/go.mod
ADD common/go.sum ./common/go.sum

RUN cd /build/${target} && go mod download

ADD ${target}/. ${target}/.
ADD common/. common/.

RUN go build -C ${target} -o /build/a.out

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/a.out ./a.out

CMD [ "./a.out" ]