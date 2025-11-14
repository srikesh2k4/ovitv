import Link from 'next/link'

export default async function DashboardPage() {
  // In real app, fetch reports & bans from backend
  const reports = [
    { id: 1, reporter: '192.168.1.100', reported: '192.168.1.101', reason: 'Spam', time: '2025-04-05 10:30' },
  ]
  const banned = [
    { ip: '192.168.1.101', reason: 'Abuse', bannedAt: '2025-04-05 10:35' },
  ]

  return (
    <div className="min-h-screen bg-gray-100 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold">Admin Dashboard</h1>
          <Link href="/admin/logout" className="text-red-600 hover:underline font-medium">
            Logout
          </Link>
        </div>

        {/* Reports */}
        <div className="bg-white p-6 rounded-lg shadow mb-8">
          <h2 className="text-2xl font-semibold mb-4">Recent Reports</h2>
          <div className="overflow-x-auto">
            <table className="w-full table-auto">
              <thead className="bg-gray-200">
                <tr>
                  <th className="p-3 text-left">From</th>
                  <th className="p-3 text-left">To</th>
                  <th className="p-3 text-left">Reason</th>
                  <th className="p-3 text-left">Time</th>
                </tr>
              </thead>
              <tbody>
                {reports.map(r => (
                  <tr key={r.id} className="border-b">
                    <td className="p-3">{r.reporter}</td>
                    <td className="p-3">{r.reported}</td>
                    <td className="p-3">{r.reason}</td>
                    <td className="p-3">{r.time}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Ban IP */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-2xl font-semibold mb-4">Ban IP Address</h2>
          <form method="POST" action="/admin/ban" className="flex flex-col md:flex-row gap-3">
            <input
              name="ip"
              placeholder="IP Address"
              className="flex-1 p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              required
            />
            <input
              name="reason"
              placeholder="Reason"
              className="flex-1 p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              required
            />
            <button
              type="submit"
              className="px-6 py-3 bg-red-600 text-white rounded-lg font-bold hover:bg-red-700 transition"
            >
              Ban IP
            </button>
          </form>

          <h3 className="text-xl font-semibold mt-6 mb-3">Banned IPs</h3>
          <div className="overflow-x-auto">
            <table className="w-full table-auto">
              <thead className="bg-gray-200">
                <tr>
                  <th className="p-3 text-left">IP</th>
                  <th className="p-3 text-left">Reason</th>
                  <th className="p-3 text-left">Banned At</th>
                  <th className="p-3 text-left">Action</th>
                </tr>
              </thead>
              <tbody>
                {banned.map((b, i) => (
                  <tr key={i} className="border-b">
                    <td className="p-3">{b.ip}</td>
                    <td className="p-3">{b.reason}</td>
                    <td className="p-3">{b.bannedAt}</td>
                    <td className="p-3">
                      <a href={`/admin/unban?ip=${b.ip}`} className="text-red-600 hover:underline">
                        Unban
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}