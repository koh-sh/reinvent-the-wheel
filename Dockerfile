FROM alpine

WORKDIR /app

COPY chatgpt ./
COPY templates/ /app/templates
EXPOSE 8080
CMD [ "./chatgpt" ]
