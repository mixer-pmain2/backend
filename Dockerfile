FROM node:16-alpine as build-frontend

WORKDIR /app

COPY ./frontend/public ./public
COPY ./frontend/src ./src
COPY ./frontend/.babelrc.json ./
COPY ./frontend/package.json ./
COPY ./frontend/webpack.config.js ./
COPY ./frontend/yarn.lock ./
RUN mkdir build

RUN yarn install
RUN yarn build


FROM golang:1.17-alpine as build-backend

WORKDIR /go/src/pmain2

COPY cmd/app/main.go ./
COPY internal ./internal
COPY pkg ./pkg
COPY go.mod ./
COPY go.sum ./

RUN go mod download

ENV GO111MODULE=on
RUN go build -o /pmain2 main.go


FROM alpine

WORKDIR /app

COPY --from=build-backend /pmain2 ./
COPY --from=build-frontend /app/build/static ./static
COPY --from=build-frontend /app/build/index.html ./
COPY frontend.routes ./
COPY .env ./


EXPOSE 80

CMD ["/app/pmain2"]

