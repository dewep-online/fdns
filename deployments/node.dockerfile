FROM node:15.12.0-alpine

RUN apk update && \
    apk add --virtual build-dependencies build-base gcc && \
    apk add --no-cache bash git make python3 && \
    npm install -g npm@7.7.5 && \
    npm install -g @angular/cli@11.2.6

EXPOSE 4200
WORKDIR /app

HEALTHCHECK --interval=15s --timeout=15s --retries=60 \
  CMD netstat -an | grep 4200 > /dev/null; if [ 0 != $? ]; then exit 1; fi;