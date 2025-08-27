FROM golang:1.24

WORKDIR /home/chall/
COPY ./ .
EXPOSE 1337

RUN go build

CMD [ "./echo-server" ]