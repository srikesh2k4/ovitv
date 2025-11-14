'use client'

import { useState } from 'react'

interface FilterModalProps {
  onApply: (filters: Filters) => void
}

interface Filters {
  gender: string
  minAge: number
  maxAge: number
  interest: string
}

export default function FilterModal({ onApply }: FilterModalProps) {
  const [gender, setGender] = useState('any')
  const [minAge, setMinAge] = useState(18)
  const [maxAge, setMaxAge] = useState(99)
  const [interest, setInterest] = useState('any')

  const handleSubmit = () => {
    onApply({ gender, minAge, maxAge, interest })
  }

  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-900 p-8 rounded-2xl w-full max-w-md shadow-2xl">
        <h2 className="text-3xl font-bold mb-8 text-center bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
          Find Your Match
        </h2>

        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium mb-2">I am</label>
            <div className="grid grid-cols-2 gap-3">
              <button
                onClick={() => setGender('male')}
                className={`p-3 rounded-lg font-medium transition ${
                  gender === 'male'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
                }`}
              >
                Male
              </button>
              <button
                onClick={() => setGender('female')}
                className={`p-3 rounded-lg font-medium transition ${
                  gender === 'female'
                    ? 'bg-pink-600 text-white'
                    : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
                }`}
              >
                Female
              </button>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Looking for</label>
            <select
              value={interest}
              onChange={(e) => setInterest(e.target.value)}
              className="w-full p-3 bg-gray-800 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              <option value="any">Anyone</option>
              <option value="male">Men</option>
              <option value="female">Women</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Age Range</label>
            <div className="flex items-center gap-3">
              <input
                type="number"
                value={minAge}
                onChange={(e) => setMinAge(Math.max(18, Math.min(98, +e.target.value)))}
                className="w-full p-3 bg-gray-800 rounded-lg text-white text-center focus:outline-none focus:ring-2 focus:ring-purple-500"
                min="18"
                max="98"
              />
              <span className="text-gray-400">—</span>
              <input
                type="number"
                value={maxAge}
                onChange={(e) => setMaxAge(Math.min(99, Math.max(19, +e.target.value)))}
                className="w-full p-3 bg-gray-800 rounded-lg text-white text-center focus:outline-none focus:ring-2 focus:ring-purple-500"
                min="19"
                max="99"
              />
            </div>
          </div>

          <button
            onClick={handleSubmit}
            className="w-full py-4 bg-gradient-to-r from-purple-600 to-pink-600 rounded-xl font-bold text-lg hover:from-purple-700 hover:to-pink-700 transition shadow-lg"
          >
            Start Chatting
          </button>
        </div>
      </div>
    </div>
  )
}