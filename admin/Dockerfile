# ---- Build ----
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package files
COPY package*.json ./
COPY package-lock.json ./

# Install deps (including dev)
RUN npm install

# Copy source
COPY . .

# Build
RUN npm run build

# ---- Runtime ----
FROM node:18-alpine
WORKDIR /app

# Copy built app
COPY --from=builder /app/.next ./.next

# ❌ REMOVE this (public folder doesn't exist)
# COPY --from=builder /app/public ./public

# Copy package files
COPY --from=builder /app/package*.json ./

# Install only production deps
RUN npm ci --only=production

EXPOSE 3001
CMD ["npm", "start"]
