FROM nginx:alpine

RUN apk update && apk add --no-cache inotify-tools

COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /etc/nginx/streams.d

COPY --chmod=755 nginx_reloader.sh /docker-entrypoint.d/nginx_reloader.sh