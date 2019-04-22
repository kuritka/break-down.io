FROM golang:1.8

# specifying default value if user is not defined
ARG USER=gouser

WORKDIR /go/src/app

#copy from build context (mostly .) to just defined workdir
COPY . .

# RUN go mod init github.com/pennylab.io/calendarBackend
RUN go mod init

RUN go mod download

RUN go mod vendor

RUN adduser -D $USER

RUN go build

#WORKDIR /home/$USER
#
#USER $USER


# what gets executed after container is created from the image
CMD ["main"]
#ENTRYPOINT ["aws"]