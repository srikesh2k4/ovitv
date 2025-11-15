# ---------------------
# 1. Builder
# ---------------------
FROM node:18-alpine AS builder

WORKDIR /app

# Install deps
COPY package.json package-lock.json ./
RUN npm install --production=false

# Copy entire project
COPY . .

# Build Next.js
RUN npm run build


# ---------------------
# 2. Runner
# ---------------------
FROM node:18-alpine

WORKDIR /app

# Copy standalone Next.js build
COPY --from=builder /app/.next/standalone/ ./
COPY --from=builder /app/.next/static/ ./.next/static/

# ❌ DO NOT COPY PUBLIC (it does NOT exist)
# Next.js 14 standalone works fine without it

ENV NODE_ENV=production
ENV PORT=3001

EXPOSE 3001

CMD ["node", "server.js"]
