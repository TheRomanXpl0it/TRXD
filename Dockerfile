# FROM node:24-alpine AS frontend

# WORKDIR /app
# COPY ./frontend/package*.json ./
# RUN npm install
# COPY ./frontend ./
# RUN npm run build


FROM golang:1.24-alpine AS backend

WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend/ .

RUN go build .


FROM alpine:latest

# RUN apk --no-cache add ca-certificates

WORKDIR /app

# COPY --from=frontend /app/build ./frontend
COPY --from=backend /app/frontend ./frontend

COPY --from=backend /app/trxd ./trxd
COPY --from=backend /app/static ./static
COPY --from=backend /app/sql ./sql

EXPOSE 1337

CMD ["/app/trxd"]
