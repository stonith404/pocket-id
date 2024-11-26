# Stage 1: Build Frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY ./frontend/package*.json ./
RUN npm ci
COPY ./frontend ./
RUN npm run build
RUN npm prune --production

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend
WORKDIR /app/backend
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

RUN apk add --no-cache gcc musl-dev

COPY ./backend ./
WORKDIR /app/backend/cmd
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/backend/pocket-id-backend .

# Stage 3: Production Image
FROM node:20-alpine
RUN apk add --no-cache caddy sqlite uuidgen

USER node
COPY --chown=node ./reverse-proxy /etc/caddy/

WORKDIR /app
COPY --from=frontend --chown=node /app/frontend/build ./frontend/build
COPY --from=frontend --chown=node /app/frontend/node_modules ./frontend/node_modules
COPY --from=frontend --chown=node /app/frontend/package.json ./frontend/package.json

COPY --from=backend --chown=node /app/backend/pocket-id-backend ./backend/pocket-id-backend
COPY --from=backend --chown=node /app/backend/migrations ./backend/migrations
COPY --from=backend --chown=node /app/backend/email-templates ./backend/email-templates
COPY --from=backend --chown=node /app/backend/images ./backend/images

COPY --chown=node ./scripts ./scripts
RUN chmod +x ./scripts/*.sh
RUN mkdir -p /app/backend/data

EXPOSE 80
ENV APP_ENV=production
ENTRYPOINT ["sh", "./scripts/docker/entrypoint.sh"]