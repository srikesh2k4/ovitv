import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'OmeTV Pro',
  description: 'Anonymous video chat with filters',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-gradient-to-br from-purple-900 to-pink-900">
        {children}
      </body>
    </html>
  )
}