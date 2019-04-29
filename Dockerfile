#docker run -p 3000:8080  auth:1.0
#docker build --build-arg USER=gouser --target release-stage -t kuritka/auth:1.6 .
#or running those two commands will correctly tag particular stages
#sudo docker build --target build-stage -t auth-builder-image:1.0 . \
#sudo docker build --build-arg USER=gouser --target release-stage -t kuritka/auth:1.6 .

#we use multistage docker image. First stage is called build and is used
#for building app in container where go and git exists..
FROM golang:1.12 as build-stage

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
# FROM alpine:latest as release-stage +5MB and is easy to troubleshoot (bash is present)
FROM scratch as release-stage

WORKDIR /app

#multistage containers - copying from build stage /build to /app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /build/static /app/static
COPY --from=build-stage /build/templates /app/templates
COPY --from=build-stage /build/main /app/main
COPY --from=build-stage /build/config.json /app/config.json

ENTRYPOINT ["./main"]
#19.6MB

#delete all <none> images
#sudo docker rmi $(sudo docker images | grep "^<none>" | awk '{ print $3 }')
#docker container prune
#docker image prune
#docker network prune
#docker volume prune
