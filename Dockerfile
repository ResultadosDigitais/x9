FROM golang:1.13.1

COPY .  $GOPATH/src/github.com/ResultadosDigitais/x9
WORKDIR $GOPATH/src/github.com/ResultadosDigitais/x9
RUN ls
RUN go build -o worker main.go 