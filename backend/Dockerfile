FROM golang:1.23-alpine

LABEL version="1.0" \
        description="Groupie-tracker Dockerized Application"\
        author="PR-MR" \
        maintainer="parisa.rahimi@gritlab.ax" \
        environment="development"

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

EXPOSE 8082

CMD ["go","run","."]