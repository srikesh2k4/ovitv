'use client'

import { useEffect, useRef, useState } from 'react'

interface VideoChatProps {
  ws: WebSocket
  onNext: () => void
}

export default function VideoChat({ ws, onNext }: VideoChatProps) {
  const localVideoRef = useRef<HTMLVideoElement>(null)
  const remoteVideoRef = useRef<HTMLVideoElement>(null)
  const [connected, setConnected] = useState(false)
  const [partnerIP, setPartnerIP] = useState('')
  const [stream, setStream] = useState<MediaStream | null>(null)

  useEffect(() => {
    let pc: RTCPeerConnection | null = null

    const start = async () => {
      try {
        const mediaStream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        })
        setStream(mediaStream)
        if (localVideoRef.current) {
          localVideoRef.current.srcObject = mediaStream
        }

        pc = new RTCPeerConnection({
          iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
        })

        mediaStream.getTracks().forEach((track) => {
          pc!.addTrack(track, mediaStream)
        })

        pc.ontrack = (event) => {
          if (remoteVideoRef.current) {
            remoteVideoRef.current.srcObject = event.streams[0]
          }
        }

        pc.onicecandidate = (event) => {
          if (event.candidate) {
            ws.send(
              JSON.stringify({
                type: 'candidate',
                candidate: event.candidate.toJSON(),
              })
            )
          }
        }

        ws.onmessage = async (e) => {
          const msg = JSON.parse(e.data)

          if (msg.type === 'ip') {
            setPartnerIP(msg.ip)
            setConnected(true)
          }

          if (msg.type === 'offer') {
            await pc!.setRemoteDescription(new RTCSessionDescription(msg))
            const answer = await pc!.createAnswer()
            await pc!.setLocalDescription(answer)
            ws.send(JSON.stringify({ type: 'answer', sdp: answer.sdp }))
          }

          if (msg.type === 'answer') {
            await pc!.setRemoteDescription(new RTCSessionDescription(msg))
          }

          if (msg.type === 'candidate') {
            await pc!.addIceCandidate(new RTCIceCandidate(msg.candidate))
          }
        }

        // Create offer
        const offer = await pc.createOffer()
        await pc.setLocalDescription(offer)
        ws.send(JSON.stringify({ type: 'offer', sdp: offer.sdp }))
      } catch (err) {
        console.error('Media error:', err)
      }
    }

    start()

    return () => {
      pc?.close()
      stream?.getTracks().forEach((t) => t.stop())
    }
  }, [ws])

  return (
    <div className="flex flex-col items-center p-4">
      <div className="flex flex-col md:flex-row gap-4 mb-6">
        <div className="relative">
          <video
            ref={localVideoRef}
            autoPlay
            muted
            playsInline
            className="w-80 h-60 bg-black rounded-xl shadow-lg object-cover"
          />
          <div className="absolute bottom-2 left-2 bg-black/60 px-2 py-1 rounded text-xs">
            You
          </div>
        </div>

        <div className="relative">
          <video
            ref={remoteVideoRef}
            autoPlay
            playsInline
            className="w-80 h-60 bg-gray-900 rounded-xl shadow-lg object-cover"
          />
          {connected ? (
            <div className="absolute bottom-2 left-2 bg-green-600 px-2 py-1 rounded text-xs">
              Live • {partnerIP}
            </div>
          ) : (
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="animate-spin rounded-full h-12 w-12 border-4 border-purple-500 border-t-transparent"></div>
            </div>
          )}
        </div>
      </div>

      <div className="flex gap-4">
        <button
          onClick={onNext}
          className="px-8 py-3 bg-red-600 text-white rounded-xl font-bold hover:bg-red-700 transition shadow-lg"
        >
          Next
        </button>
        <button className="px-8 py-3 bg-gray-700 text-white rounded-xl font-bold hover:bg-gray-600 transition shadow-lg">
          Report
        </button>
      </div>
    </div>
  )
}