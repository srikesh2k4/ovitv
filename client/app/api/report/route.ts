import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest) {
  try {
    const { reported_ip, reason } = await req.json();

    const backendURL = process.env.NEXT_PUBLIC_BACKEND_URL || "https://ovitv-backend.up.railway.app";

    const res = await fetch(`${backendURL}/api/report`, {
      method: 'POST',
      body: JSON.stringify({ reported_ip, reason }),
      headers: { 'Content-Type': 'application/json' },
    });

    if (!res.ok) {
      return NextResponse.json({ success: false, error: 'Backend error' }, { status: 500 });
    }

    return NextResponse.json({ success: true });
  } catch (error) {
    return NextResponse.json({ success: false, error: 'Server error' }, { status: 500 });
  }
}
