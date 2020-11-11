# builder
FROM node:lts-alpine as builder
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
ARG server_url
RUN VUE_APP_SERVER_URL=${server_url} npm run build

# runner
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]