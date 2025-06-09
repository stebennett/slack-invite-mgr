import React, { useEffect, useState } from 'react';

interface Invite {
  name: string;
  role: string;
  email: string;
  company: string;
  yearsExperience: string;
  reasons: string;
  source: string;
}

export const InvitesTable: React.FC = () => {
  const [invites, setInvites] = useState<Invite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchInvites = async () => {
      try {
        const response = await fetch('/api/invites');
        if (!response.ok) {
          throw new Error('Failed to fetch invites');
        }
        const data = await response.json();
        setInvites(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchInvites();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white border border-gray-300">
        <thead>
          <tr className="bg-gray-100">
            <th className="px-6 py-3 border-b text-left">Name</th>
            <th className="px-6 py-3 border-b text-left">Role</th>
            <th className="px-6 py-3 border-b text-left">Email</th>
            <th className="px-6 py-3 border-b text-left">Company</th>
            <th className="px-6 py-3 border-b text-left">Years Experience</th>
            <th className="px-6 py-3 border-b text-left">Reasons</th>
            <th className="px-6 py-3 border-b text-left">Source</th>
          </tr>
        </thead>
        <tbody>
          {invites.map((invite, index) => (
            <tr key={index} className="hover:bg-gray-50">
              <td className="px-6 py-4 border-b">{invite.name}</td>
              <td className="px-6 py-4 border-b">{invite.role}</td>
              <td className="px-6 py-4 border-b">{invite.email}</td>
              <td className="px-6 py-4 border-b">{invite.company}</td>
              <td className="px-6 py-4 border-b">{invite.yearsExperience}</td>
              <td className="px-6 py-4 border-b">{invite.reasons}</td>
              <td className="px-6 py-4 border-b">{invite.source}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}; 