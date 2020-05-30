FROM golang:1.12.9
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]
LABEL "Author"="Qarapaiym"

EXPOSE 8080