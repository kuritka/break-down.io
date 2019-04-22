FROM golang:1.12

# specifying default value if user is not defined
ARG USER=appuser

RUN mkdir /app

WORKDIR /app

#copy from build context (mostly .) to just defined workdir
COPY . .

# RUN go mod init github.com/pennylab.io/calendarBackend
RUN go mod init --force

RUN go mod download

RUN go mod vendor

RUN go build -o main .

RUN adduser -S -D -H -h /app $USER

USER $USER


# what gets executed after container is created from the image
CMD ["./main"]
#ENTRYPOINT ["aws"]