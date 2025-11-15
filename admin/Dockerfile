# ---------------------
# 1. Builder
# ---------------------
FROM node:18-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY package.json package-lock.json ./
RUN npm install --production=false

# Copy entire project
COPY . .

# Build Next.js (creates .next/standalone)
RUN npm run build


# ---------------------
# 2. Runner
# ---------------------
FROM node:18-alpine

WORKDIR /app

# Copy standalone output (contains server.js + internal node_modules)
COPY --from=builder /app/.next/standalone/ ./

# Copy static files
COPY --from=builder /app/.next/static/ ./.next/static/

# Copy public folder ONLY if exists
RUN mkdir -p /app/public
COPY --from=builder /app/public/ ./public/ || true

ENV NODE_ENV=production
ENV PORT=3001

EXPOSE 3001

# Start Next.js standalone server
CMD ["node", "server.js"]
