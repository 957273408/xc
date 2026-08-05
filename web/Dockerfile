FROM node:22-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV HTTP_PROXY=$HTTP_PROXY
ENV HTTPS_PROXY=$HTTPS_PROXY
ENV NO_PROXY=$NO_PROXY
COPY . /app
WORKDIR /app


FROM base AS prod-deps


FROM base AS build


FROM nginx:alpine
LABEL MAINTAINER="bypanghu@163.com"
COPY --from=base /app/.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf
COPY --from=build /app/dist /usr/share/nginx/html
RUN ls -al /usr/share/nginx/html
