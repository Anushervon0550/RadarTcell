# Multi-stage Dockerfile for @radartcell/web-public.
# Build context: repo root. See docker-compose.yml.

# -------------------- deps --------------------
FROM node:20-alpine AS deps
WORKDIR /repo
RUN corepack enable

# Copy workspace manifests only to leverage layer caching.
COPY frontend/package.json frontend/pnpm-workspace.yaml frontend/.npmrc ./
COPY frontend/packages/api/package.json ./packages/api/
COPY frontend/packages/ui/package.json ./packages/ui/
COPY frontend/apps/web-public/package.json ./apps/web-public/
COPY frontend/apps/web-admin/package.json ./apps/web-admin/

RUN pnpm install --frozen-lockfile=false

# -------------------- build --------------------
FROM node:20-alpine AS build
WORKDIR /repo
RUN corepack enable

COPY --from=deps /repo/node_modules ./node_modules
COPY --from=deps /repo/apps/web-public/node_modules ./apps/web-public/node_modules
COPY --from=deps /repo/packages/api/node_modules ./packages/api/node_modules
COPY --from=deps /repo/packages/ui/node_modules ./packages/ui/node_modules

COPY frontend/package.json frontend/pnpm-workspace.yaml frontend/.npmrc ./
COPY frontend/packages ./packages
COPY frontend/apps/web-public ./apps/web-public

# Public build-time env. Same-origin by default: leave VITE_API_BASE empty.
ARG VITE_API_BASE=
ARG VITE_ADMIN_URL=/admin
ARG VITE_DEFAULT_LOCALE=ru
ENV VITE_API_BASE=$VITE_API_BASE \
    VITE_ADMIN_URL=$VITE_ADMIN_URL \
    VITE_DEFAULT_LOCALE=$VITE_DEFAULT_LOCALE

RUN pnpm --filter @radartcell/web-public build

# -------------------- runtime --------------------
FROM nginx:1.27-alpine AS runtime
COPY --from=build /repo/apps/web-public/dist /usr/share/nginx/html
COPY deploy/nginx/web-public.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
HEALTHCHECK --interval=15s --timeout=3s --retries=5 \
    CMD wget -qO- http://127.0.0.1/ >/dev/null || exit 1
