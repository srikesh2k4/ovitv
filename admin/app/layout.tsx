import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'OmeTV Pro Admin',
  description: 'Admin panel for OmeTV Pro',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}