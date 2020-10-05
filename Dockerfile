# Stage 1 testing
FROM node:14.7.0-alpine as builder

COPY . /data

# build
RUN cd /data &&\
  npm ci &&\
  npm test &&\
  npm run jshint_backend

# Stage 2 package
FROM node:14.7.0-alpine

COPY . /data

RUN cd /data && npm ci --only=production
COPY --from=builder /data/public /data/public

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
