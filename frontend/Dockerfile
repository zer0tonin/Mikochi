FROM node:19-alpine

RUN mkdir /app
COPY package.json /app/

WORKDIR /app
RUN npm install

EXPOSE 8080

CMD ["npm", "run", "dev"]
