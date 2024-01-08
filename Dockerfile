ARG NODE_BASE_IMAGE=node:21.5.0-alpine3.18
FROM ${NODE_BASE_IMAGE} as node-base

FROM golang:1.21.5-alpine3.18 as backend-builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /go/build
COPY app/backend /go/build
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /cms

FROM node-base as frontend-builder
WORKDIR /app
COPY app/frontend /app
RUN npm install
RUN npm run build

FROM node-base
COPY --from=caddy:2.7.5 /usr/bin/caddy /usr/local/bin/caddy
RUN apk add --no-cache git openssh-client supervisor
WORKDIR /app
COPY build/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY build/Caddyfile /etc/caddy/Caddyfile
COPY --from=backend-builder  /cms ./cms
COPY --from=frontend-builder /app/node_modules ./node_modules
COPY --from=frontend-builder /app/.next ./.next
COPY --from=frontend-builder /app/package.json ./package.json
COPY --from=frontend-builder /app/public ./public
HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
CMD [ "/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf" ]
