FROM golang:1.16.4

ENV APP_NAME nomi_mux
ENV PORT 8080

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go mod download
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${PORT}