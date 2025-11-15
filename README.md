# Mini OmeTV-like — Next.js + WebSocket Signaling (Deployable on Railway)

This repository contains a minimal, deployable example of a random video chat app (OmeTV-like) using:
- Next.js frontend (client/)
- Express + Next custom server + WebSocket signaling (server.js)
- Simple matchmaking queue and signaling (offer/answer/ice)
- Report endpoint (`POST /api/report`)
- Security basics: helmet, rate-limiter (in server), CORS

## How it works
- The frontend obtains camera/mic and connects to the signaling WebSocket at `/signaling`.
- The server pairs two WebSocket clients and forwards SDP offers/answers and ICE candidates.
- The frontend uses RTCPeerConnection to establish direct peer-to-peer WebRTC media.

## Development
1. Install dependencies:
```bash
npm install
cd client
npm install
```

2. Start dev:
```bash
npm run dev
```
Open `http://localhost:3000`.

## Production build
```bash
npm run build
npm start
```
This builds the Next app then runs `server.js` which serves the built frontend and handles signaling.

## Deploy to Railway
1. Push this repository to GitHub.
2. Create a new Railway project and connect your GitHub repo.
3. Railway will detect Node; set `Start Command` to `npm start`.
4. Ensure Railway exposes the port via `PORT` env (Railway does this automatically).
5. Optionally set `NODE_ENV=production` in Railway environment variables.

For a Docker deployment, you can add the provided `Dockerfile` and set Railway to build with Docker.

## Security & production notes
- This is a minimal example. For production:
  - Use HTTPS (Railway provides TLS).
  - Add user authentication (avoid fully anonymous systems without moderation).
  - Store reports in a persistent DB and integrate a moderation pipeline.
  - Limit and sanitize any metadata exchanged.
  - Implement abuse prevention, profanity filters, and optionally content moderation via ML.
  - Consider TURN server (recommended) for NAT traversal / mobile networks.
    - Configure `RTCPeerConnection` to use a TURN server (coturn) when needed.

## Files of interest
- `server.js` — custom server + signaling
- `client/pages/index.js` — WebRTC UI & signaling logic
- `client/styles/globals.css` — basic styling
- `Dockerfile` — Docker image for deployment

