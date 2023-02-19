FROM nginx:alpine

RUN apk update && apk add --no-cache inotify-tools

COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /etc/nginx/streams.d

COPY --chmod=755 nginx_reloader.sh /nginx_reloader.sh
COPY --chmod=755 99-start-reloader.sh /docker-entrypoint.d/99-start-reloader.sh