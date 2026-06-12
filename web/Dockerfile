# 如果需要用 cicd ，请设置环境变量：
# variables:
#     DOCKER_BUILDKIT: 1

FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM nginx:alpine
LABEL MAINTAINER="bypanghu@163.com"
COPY /app/.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf
COPY /app/dist /usr/share/nginx/html
RUN ls -al /usr/share/nginx/html
