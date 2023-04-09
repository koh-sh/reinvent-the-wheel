FROM golang:1.20

WORKDIR /app

COPY chatgpt.go ./
COPY templates/ /app/templates
RUN CGO_ENABLED=0 GOOS=linux go build ./chatgpt.go

EXPOSE 8080
CMD [ "./chatgpt" ]
