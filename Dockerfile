#docker run -p 3000:8080  auth:1.0
#docker build --build-arg USER=gouser -t auth:1.6 .

#we use multistage docker image. First stage is called build and is used
#for building app in container where go and git exists..
FROM golang:1.12 as build

# specifying default value if user is not defined
ARG USER=appuser

# helps with debugging because shows dokcer instruction
RUN set -x

RUN mkdir /build

WORKDIR /build

#copy from build context to workdir
# COPY . .

#staicaly linking go libraries
ENV CGO_ENABLED=0
#Do all in one layer otherwise downloaded data stay in previous layers
#and if we delete them in upcomming layer it will not have effect
RUN git clone http://github.com/kuritka/break-down.io . && \
    go mod init github.com/kuritka/break-down.io && \
    go list -e $(go list -f . -m all) && \
    go build -o main .
    #mv './templates' './static' 'main' 'config.json' /app && \
    ##remove all files directories and hidden dirs like .git, .idea...
    #rm -rf .[^.] .??* ./*
#960MB

#-----------------------------------------------------
#FROM scratch as release - -5MB and is difficult to troubleshoot (missing bash)
FROM alpine:latest as release

#installing certificates otherwise app doesnt connect to github
#this needs to be extra solved for scratch image (+2MB)
RUN     set -x \
    &&  mkdir /app \
    &&  apk update \
    &&  apk upgrade \
    &&  apk add --no-cache \
            ca-certificates \
    && update-ca-certificates 2>/dev/null || true

WORKDIR /app

#multistage containers - copying from build stage /build to /app
COPY --from=build /build/static /app/static
COPY --from=build /build/templates /app/templates
COPY --from=build /build/main /app/main
COPY --from=build /build/config.json /app/config.json

ENTRYPOINT ["./main"]
#27MB

#delete all <none> images
#sudo docker rmi $(sudo docker images | grep "^<none>" | awk '{ print $3 }')
#docker container prune
#docker image prune
#docker network prune
#docker volume prune
