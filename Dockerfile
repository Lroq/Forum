FROM golang:1.21.0

WORKDIR /app

# copi les fichiers de l'appli dans conteneur
COPY . .

# le port 8080 pour HTTP et 8443 pour HTTPS
EXPOSE 8080
EXPOSE 8443

# construit l'appli
RUN go build -o main .

CMD ["/app/main"]