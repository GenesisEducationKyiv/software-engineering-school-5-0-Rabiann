FROM node:19
FROM mcr.microsoft.com/playwright:v1.53.1
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json /app/
COPY tests/ /app/tests/
RUN apt-get update && apt-get -y install libnss3 libatk-bridge2.0-0 libdrm-dev libxkbcommon-dev libgbm-dev libasound-dev libatspi2.0-0 libxshmfence-dev
RUN npm install
RUN npx playwright test