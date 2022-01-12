# Stage 1 testing
FROM node:17.3.1-alpine as builder

COPY . /data

# build
RUN cd /data &&\
  npm ci &&\
  npm test &&\
  npm run jshint

# Stage 2 package
FROM node:17.3.1-alpine

COPY . /data

RUN cd /data && npm ci --only=production

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
