# ---------------------
# 1. Builder
# ---------------------
FROM node:18-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy rest of the app
COPY . .

# Build Next.js
RUN npm run build


# ---------------------
# 2. Runner
# ---------------------
FROM node:18-alpine

WORKDIR /app

# Copy standalone build (Next.js 14 generates server.js + .next)
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Safe: Only copy public IF exists (won't break if missing)
RUN mkdir -p ./public
COPY --from=builder /app/public ./public 2>/dev/null || true

ENV NODE_ENV=production
ENV PORT=3001

EXPOSE 3001

# Run Next.js standalone server
CMD ["node", "server.js"]
