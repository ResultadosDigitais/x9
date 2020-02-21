FROM golang:1.13.1

COPY .  $GOPATH/src/github.com/ResultadosDigitais/x9
WORKDIR $GOPATH/src/github.com/ResultadosDigitais/x9
RUN go run  main.go 
