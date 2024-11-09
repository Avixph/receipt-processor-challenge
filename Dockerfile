ARG BINARY_PATH="tmp/bin/api"
ARG MAIN_PATH="cmd/api"
ARG MIGRATION_PATH="./migrations"

#FROM golang:1.23.2-alpine AS BUILD
#WORKDIR /server
#RUN cd server  \
#    && COPY go.mod go.sum ./
#RUN go mod download  \
#    && COPY ..
#RUN go build -o /${BINARY_PATH} ./${MAIN_PATH}

#FROM golang:1.23.2-alpine AS BUILD
#WORKDIR /server
#COPY server /server
#RUN go build -o /${BINARY_PATH} ./${MAIN_PATH}

FROM golang:1.23.2-alpine AS BUILD
WORKDIR /server
COPY server /server/
RUN go mod download
RUN go build -ldflags "-w -s" -o /${BINARY_PATH} ./${MAIN_PATH}

FROM alpine:3.20.3 AS RUNNER
COPY --from=BUILD /server/${BINARY_PATH} ./
EXPOSE 8080
CMD ["./${BINARY_PATH}"]