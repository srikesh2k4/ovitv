'use client'

import { useState } from 'react'
import FilterModal from '@/components/FilterModal'
import VideoChat from '@/components/VideoChat'

interface Filters {
  gender: string
  minAge: number
  maxAge: number
  interest: string
}

export default function Home() {
  const [filters, setFilters] = useState<Filters | null>(null)
  const [ws, setWs] = useState<WebSocket | null>(null)

  const handleApplyFilters = (f: Filters) => {
    setFilters(f)
    const socket = new WebSocket('wss://client.localhost/ws')
    
    socket.onopen = () => {
      socket.send(
        JSON.stringify({
          type: 'join',
          gender: f.gender,
          age: 25,
          interest: f.interest,
          minAge: f.minAge,
          maxAge: f.maxAge,
        })
      )
    }

    setWs(socket)
  }

  const handleNext = () => {
    ws?.close()
    setFilters(null)
  }

  return (
    <main className="min-h-screen flex items-center justify-center p-4">
      {!filters ? (
        <FilterModal onApply={handleApplyFilters} />
      ) : (
        <VideoChat ws={ws!} onNext={handleNext} />
      )}
    </main>
  )
}