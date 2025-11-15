# Simple Dockerfile to build and run the app on Railway (or other hosts)
FROM node:20-alpine

WORKDIR /app
COPY package.json package-lock.json* ./
COPY client/package.json ./client/
RUN npm install --production --ignore-scripts && cd client && npm install --production

COPY . .

# Build Next client
RUN npm run build

EXPOSE 3000
ENV NODE_ENV=production
CMD ["npm", "start"]
