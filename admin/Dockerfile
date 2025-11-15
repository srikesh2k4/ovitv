# ---------------------
# 1. Builder
# ---------------------
FROM node:18-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install --production=false

# Copy full project
COPY . .

# Build Next.js standalone output
RUN npm run build


# ---------------------
# 2. Runner
# ---------------------
FROM node:18-alpine

WORKDIR /app

# Copy Next.js standalone build
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static/

# Environment port
ENV NODE_ENV=production
ENV PORT=3001

# Expose the same port Next.js uses
EXPOSE 3001

# Run standalone server
CMD ["node", "server.js"]
