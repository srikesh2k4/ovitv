const express = require('express');
const next = require('next');
const http = require('http');
const WebSocket = require('ws');
const helmet = require('helmet');
const cors = require('cors');
const { RateLimiterMemory } = require('rate-limiter-flexible');

const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev, dir: './client' });
const handle = app.getRequestHandler();

const PORT = process.env.PORT || 3000;
const RATE_POINTS = parseInt(process.env.RATE_POINTS || '10');
const RATE_DURATION = parseInt(process.env.RATE_DURATION || '1'); // seconds

const rateLimiter = new RateLimiterMemory({
  points: RATE_POINTS,
  duration: RATE_DURATION,
});

app.prepare().then(() => {
  const serverApp = express();
  serverApp.use(helmet());
  serverApp.use(cors());
  serverApp.use(express.json());

  // simple report endpoint (store or forward to moderation pipeline)
  serverApp.post('/api/report', async (req, res) => {
    try {
      await rateLimiter.consume(req.ip);
    } catch (_) {
      return res.status(429).json({ error: 'Too many requests' });
    }
    const { reason, details } = req.body || {};
    console.log('REPORT:', { ip: req.ip, reason, details, at: new Date().toISOString() });
    // In production: persist to DB / send to moderation queue
    return res.json({ status: 'ok' });
  });

  // fallback to Next.js for all pages
  serverApp.all('*', (req, res) => handle(req, res));

  const server = http.createServer(serverApp);

  // WebSocket server for signaling
  const wss = new WebSocket.Server({ noServer: true });

  // Simple matchmaking: queue holds sockets waiting to be paired
  const waiting = [];
  const peers = new Map(); // ws -> peerWs

  function send(ws, obj) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(obj));
    }
  }

  wss.on('connection', (ws) => {
    ws.isAlive = true;
    ws.on('pong', () => ws.isAlive = true);

    // Add to waiting queue or pair
    if (waiting.length === 0) {
      waiting.push(ws);
      send(ws, { type: 'waiting' });
    } else {
      const peer = waiting.shift();
      if (peer && peer.readyState === WebSocket.OPEN) {
        peers.set(ws, peer);
        peers.set(peer, ws);
        send(ws, { type: 'paired' });
        send(peer, { type: 'paired' });
      } else {
        // peer not available, push this ws to queue
        waiting.push(ws);
        send(ws, { type: 'waiting' });
      }
    }

    ws.on('message', (message) => {
      try {
        const data = JSON.parse(message.toString());
        // Forward signaling messages to peer
        const peer = peers.get(ws);
        if (peer && peer.readyState === WebSocket.OPEN) {
          send(peer, data);
        } else {
          // If no peer, respond with error for certain message types
          if (data && data.type === 'offer') {
            send(ws, { type: 'no-peer' });
          }
        }

        // Handle 'next' request to switch partner
        if (data && data.type === 'next') {
          const currentPeer = peers.get(ws);
          if (currentPeer) {
            // unpair both and put them back to queue separately if open
            peers.delete(ws);
            peers.delete(currentPeer);
            try { send(currentPeer, { type: 'partner-left' }); } catch(e){}
            try { send(ws, { type: 'partner-left' }); } catch(e){}
            if (currentPeer.readyState === WebSocket.OPEN) waiting.push(currentPeer);
            if (ws.readyState === WebSocket.OPEN) waiting.push(ws);
            // attempt to pair again immediately
            while (waiting.length >= 2) {
              const a = waiting.shift();
              const b = waiting.shift();
              if (a.readyState === WebSocket.OPEN && b.readyState === WebSocket.OPEN) {
                peers.set(a, b);
                peers.set(b, a);
                send(a, { type: 'paired' });
                send(b, { type: 'paired' });
              } else {
                if (a.readyState === WebSocket.OPEN) waiting.push(a);
                if (b.readyState === WebSocket.OPEN) waiting.push(b);
              }
            }
          }
        }

      } catch (err) {
        console.warn('Invalid WS msg', err);
      }
    });

    ws.on('close', () => {
      // remove from waiting if present
      const idx = waiting.indexOf(ws);
      if (idx !== -1) waiting.splice(idx, 1);
      const peer = peers.get(ws);
      if (peer) {
        peers.delete(ws);
        peers.delete(peer);
        try { send(peer, { type: 'partner-left' }); } catch(e){}
        // put peer back to queue
        if (peer.readyState === WebSocket.OPEN) waiting.push(peer);
      }
    });
  });

  // heartbeat to close dead connections
  setInterval(() => {
    wss.clients.forEach((ws) => {
      if (ws.isAlive === false) return ws.terminate();
      ws.isAlive = false;
      ws.ping(()=>{});
    });
  }, 30000);

  // Upgrade HTTP server to handle WebSocket at path /signaling
  server.on('upgrade', (request, socket, head) => {
    if (request.url === '/signaling') {
      wss.handleUpgrade(request, socket, head, (ws) => {
        wss.emit('connection', ws, request);
      });
    } else {
      socket.destroy();
    }
  });

  server.listen(PORT, () => {
    console.log(`> Ready on http://localhost:${PORT} (env=${process.env.NODE_ENV})`);
  });

});
