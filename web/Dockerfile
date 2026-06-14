FROM node:22-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV HTTP_PROXY=$HTTP_PROXY
ENV HTTPS_PROXY=$HTTPS_PROXY
ENV NO_PROXY=$NO_PROXY
RUN corepack enable && pnpm config set registry https://registry.npmmirror.com && pnpm config set fetch-timeout 600000
COPY . /app
WORKDIR /app


FROM base AS prod-deps
RUN pnpm install --prod


FROM base AS build
COPY --from=prod-deps /app/node_modules /app/node_modules
RUN pnpm install && pnpm run build


FROM nginx:alpine
LABEL MAINTAINER="bypanghu@163.com"
COPY --from=base /app/.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf
COPY --from=build /app/dist /usr/share/nginx/html
RUN ls -al /usr/share/nginx/html
