FROM golang:1.19-alpine

WORKDIR /app

COPY ./src/go.mod ./

RUN go mod download

COPY ./src ./

RUN go build -o /udppoc

EXPOSE 6001

CMD [ "/udppoc" ]