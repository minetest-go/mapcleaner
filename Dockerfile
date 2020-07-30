# Stage 1 testing
FROM node:14.6.0-alpine as builder

COPY . /data

# build
RUN cd /data &&\
  npm ci &&\
  npm test &&\
  npm run jshint_backend

# Stage 2 package
FROM node:14.6.0-alpine

COPY . /data
RUN apk update && apk add curl

RUN cd /data && npm ci --only=production
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1

CMD ["npm", "start"]
