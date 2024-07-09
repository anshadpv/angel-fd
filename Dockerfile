ARG GO_IMAGE_VER=1.21.3-alpine
FROM golang:1.21.3-alpine AS build

ARG service_name

WORKDIR /go/src/app
COPY . /go/src/app

ARG BUILD_STAMP
ARG GIT_COMMIT_ID
ARG GIT_PRIMARY_BRANCH
ARG GIT_URL
ARG HOST_NAME
ARG GIT_COMMIT_AUTHOR

RUN mkdir /go/bin/app
RUN apk add --no-cache build-base
RUN apk add --no-cache tzdata
RUN apk add --update --no-cache ca-certificates git
RUN go env -w GOPRIVATE="github.com/angel-one/*"
RUN export GOPRIVATE=github.com/angel-one/*
RUN GOOS=linux CGO_ENABLED=0 go build -a -tags musl -o  /go/bin/app/main -ldflags "-X github.com/sinhashubham95/go-actuator.BuildStamp=$BUILD_STAMP -X github.com/sinhashubham95/go-actuator.GitCommitID=$GIT_COMMIT_ID -X github.com/sinhashubham95/go-actuator.GitPrimaryBranch=$GIT_PRIMARY_BRANCH -X github.com/sinhashubham95/go-actuator.GitURL=$GIT_URL -X github.com/sinhashubham95/go-actuator.Username=shubham.sinha -X github.com/sinhashubham95/go-actuator/core.HostName=$HOST_NAME  -X github.com/sinhashubham95/go-actuator/core.GitCommitTime=$BUILD_STAMP -X github.com/sinhashubham95/go-actuator/core.GitCommitAuthor=$GIT_COMMIT_AUTHOR"
RUN mkdir /go/bin/app/log
COPY resources /go/bin/app/resources

##
# Archive #
##

FROM 732165046977.dkr.ecr.ap-south-1.amazonaws.com/sre-golang-base-image:3  AS artifact
ARG SERVICE_PORT
ARG service_name

#Copy location data
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ 	Asia/Calcutta

#Copy Aritfact
WORKDIR /app
COPY --from=build /go/bin/app /app

HEALTHCHECK NONE
EXPOSE ${SERVICE_PORT}

RUN echo -e "#!/bin/sh\n ./main --mode=release" > entrypoint.sh && \
    chmod +x entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
