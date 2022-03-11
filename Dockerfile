FROM golang

EXPOSE 8000

COPY ./ ./
RUN go build -v -o /usr/local/bin/app main.go

CMD ["app"]
