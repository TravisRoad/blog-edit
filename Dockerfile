ARG NODE_BASE_IMAGE=node:21.5.0-alpine3.18
FROM ${NODE_BASE_IMAGE} as node-base

FROM golang:1.21.5-alpine3.18 as backend-builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /go/build
COPY app/backend /go/build
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /cms

FROM node-base as frontend-deps
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add --no-cache libc6-compat
WORKDIR /app
COPY app/frontend/package.json app/frontend/yarn.lock* app/frontend/package-lock.json* app/frontend/pnpm-lock.yaml* ./
RUN \
  if [ -f yarn.lock ]; then yarn --frozen-lockfile; \
  elif [ -f package-lock.json ]; then npm ci; \
  elif [ -f pnpm-lock.yaml ]; then yarn global add pnpm && pnpm i --frozen-lockfile; \
  else echo "Lockfile not found." && exit 1; \
  fi

FROM node-base as frontend-builder
WORKDIR /app
COPY app/frontend /app
COPY --from=frontend-deps /app/node_modules ./node_modules
run yarn build

FROM node-base as runner
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add --no-cache git openssh-client supervisor

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1
RUN addgroup --system --gid 1001 cms
RUN adduser --system --uid 1001 cms
WORKDIR /home/cms

COPY --from=frontend-builder /app/public ./public

RUN mkdir .next && chown cms:cms .next

COPY --from=backend-builder  --chown=cms:cms /cms /usr/local/bin/cms
COPY --from=frontend-builder --chown=cms:cms /app/.next/standalone ./
COPY --from=frontend-builder --chown=cms:cms /app/.next/static ./.next/static

COPY --from=caddy:2.7.5 /usr/bin/caddy /usr/local/bin/caddy
COPY build/Caddyfile /etc/caddy/Caddyfile
COPY build/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

USER cms

EXPOSE 8000
ENV HOSTNAME "0.0.0.0"

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
CMD [ "/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf" ]
