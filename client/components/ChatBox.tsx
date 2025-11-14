'use client';

import { FC, useState } from 'react';
import { motion } from 'framer-motion';

interface Props {
  messages: string[];
  onSend: (text: string) => void;
}

const ChatBox: FC<Props> = ({ messages, onSend }) => {
  const [input, setInput] = useState('');

  const send = () => {
    if (input.trim()) {
      onSend(input);
      setInput('');
    }
  };

  return (
    <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-4 h-96 flex flex-col">
      <div className="flex-1 overflow-y-auto space-y-2 mb-4">
        {messages.map((m, i) => (
          <motion.div
            key={i}
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            className={`p-3 rounded-lg ${m.startsWith('You:') ? 'bg-blue-600 ml-auto max-w-xs' : 'bg-gray-700 max-w-xs'}`}
          >
            {m}
          </motion.div>
        ))}
      </div>
      <div className="flex gap-2">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && send()}
          placeholder="Type a message..."
          className="flex-1 p-3 bg-white/20 rounded-lg text-white placeholder-gray-400"
        />
        <button
          onClick={send}
          className="px-6 bg-gradient-to-r from-green-500 to-emerald-600 rounded-lg font-bold"
        >
          Send
        </button>
      </div>
    </div>
  );
};

export default ChatBox;