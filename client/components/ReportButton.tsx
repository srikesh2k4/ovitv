'use client';

import { FC } from 'react';

interface Props {
  ip: string;
}

const ReportButton: FC<Props> = ({ ip }) => {
  const report = () => {
    const reason = prompt("Why report this user?");
    if (reason && ip) {
      fetch('/api/report', {
        method: 'POST',
        body: JSON.stringify({ reported_ip: ip, reason }),
        headers: { 'Content-Type': 'application/json' },
      });
      alert("Reported!");
    }
  };

  return (
    <button
      onClick={report}
      className="absolute top-4 right-4 bg-red-600 text-white px-4 py-2 rounded-full shadow-lg hover:scale-110 transition"
    >
      Report
    </button>
  );
};

export default ReportButton;