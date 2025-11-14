import { NextRequest } from 'next/server';

export async function POST(req: NextRequest) {
  const { reported_ip, reason } = await req.json();
  await fetch('http://backend:8080/api/report', {
    method: 'POST',
    body: JSON.stringify({ reported_ip, reason }),
    headers: { 'Content-Type': 'application/json' },
  });
  return Response.json({ success: true });
}