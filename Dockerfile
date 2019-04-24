#docker run -p 3000:8080  auth:1.0
#docker build --build-arg USER=gouser -t auth:1.0 .
FROM golang:1.12

# specifying default value if user is not defined
ARG USER=appuser

RUN mkdir /app

WORKDIR /app

#copy from build context (mostly .) to just defined workdir
# COPY . .

RUN git clone http://github.com/kuritka/break-down.io

WORKDIR /app/break-down.io

RUN go mod init github.com/kuritka/break-down.io

RUN go list -e $(go list -f . -m all)

RUN go build -o main .

CMD ["./main"]

# this is not working at all llook at th ehack...
#RUN go mod init github.com/kuritka/break-down.io
#RUN go mod download
#RUN go mod vendorc
#
#RUN go build -o main .
# HACK, waiting for next go version
# Populate the module cache based on the go.{mod,sum} files.
# https://stackoverflow.com/questions/51126349/how-to-download-all-dependencies-with-vgo-and-a-given-go-mod
#COPY go.mod .
#COPY go.sum .
#RUN go list -e $(go list -f . -m all)

#RUN adduser -D $USER
#USER $USER

# what gets executed after container is created from the image

#ENTRYPOINT ["aws"]