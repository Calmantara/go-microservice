# Go Modules preparation
FROM alpine
COPY . .
ENV TZ=GMT
RUN apk --no-cache add ca-certificates
COPY ./main-app /main-app
COPY manifest /manifest
ENTRYPOINT ["/main-app"]