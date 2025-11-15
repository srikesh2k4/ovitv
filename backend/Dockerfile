# --- Builder ---
FROM node:18-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install --production=false

COPY . .
RUN npm run build

# --- Runner ---
FROM node:18-alpine

WORKDIR /app

# Copy Next.js standalone output
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Only copy public if it exists
# (safe even if empty)
COPY --from=builder /app/public ./public

ENV NODE_ENV=production
ENV PORT=3001

EXPOSE 3001

CMD ["node", "server.js"]
