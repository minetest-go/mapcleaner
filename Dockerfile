FROM node:14.7.0-alpine

COPY . /data

RUN cd /data && npm ci --only=production

WORKDIR /data

EXPOSE 8080

CMD ["npm", "start"]
