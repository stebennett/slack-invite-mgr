import React, { useEffect, useState } from 'react';

interface Invite {
  name: string;
  role: string;
  email: string;
  company: string;
  yearsExperience: string;
  reasons: string;
  source: string;
  status?: 'pending' | 'approved' | 'denied';
}

export const InvitesTable: React.FC = () => {
  const [invites, setInvites] = useState<Invite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'pending' | 'approved' | 'denied'>('pending');

  useEffect(() => {
    const fetchInvites = async () => {
      try {
        const response = await fetch('/api/invites');
        if (!response.ok) {
          throw new Error('Failed to fetch invites');
        }
        const data = await response.json();
        // Initialize all invites as pending
        setInvites(data.map((invite: Invite) => ({ ...invite, status: 'pending' })));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchInvites();
  }, []);

  const handleApprove = (email: string) => {
    setInvites(invites.map(invite => 
      invite.email === email ? { ...invite, status: 'approved' } : invite
    ));
  };

  const handleDeny = (email: string) => {
    setInvites(invites.map(invite => 
      invite.email === email ? { ...invite, status: 'denied' } : invite
    ));
  };

  const filteredInvites = invites.filter(invite => invite.status === activeTab);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <div className="mb-4">
        <div className="flex space-x-4">
          <button
            onClick={() => setActiveTab('pending')}
            className={`px-4 py-2 rounded ${
              activeTab === 'pending'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 text-gray-700'
            }`}
          >
            Pending
          </button>
          <button
            onClick={() => setActiveTab('approved')}
            className={`px-4 py-2 rounded ${
              activeTab === 'approved'
                ? 'bg-green-500 text-white'
                : 'bg-gray-200 text-gray-700'
            }`}
          >
            Approved
          </button>
          <button
            onClick={() => setActiveTab('denied')}
            className={`px-4 py-2 rounded ${
              activeTab === 'denied'
                ? 'bg-red-500 text-white'
                : 'bg-gray-200 text-gray-700'
            }`}
          >
            Denied
          </button>
        </div>
      </div>

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
              {activeTab === 'pending' && (
                <th className="px-6 py-3 border-b text-left">Actions</th>
              )}
            </tr>
          </thead>
          <tbody>
            {filteredInvites.map((invite, index) => (
              <tr key={index} className="hover:bg-gray-50">
                <td className="px-6 py-4 border-b">{invite.name}</td>
                <td className="px-6 py-4 border-b">{invite.role}</td>
                <td className="px-6 py-4 border-b">{invite.email}</td>
                <td className="px-6 py-4 border-b">{invite.company}</td>
                <td className="px-6 py-4 border-b">{invite.yearsExperience}</td>
                <td className="px-6 py-4 border-b">{invite.reasons}</td>
                <td className="px-6 py-4 border-b">{invite.source}</td>
                {activeTab === 'pending' && (
                  <td className="px-6 py-4 border-b">
                    <div className="flex space-x-2">
                      <button
                        onClick={() => handleApprove(invite.email)}
                        className="px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600"
                      >
                        Approve
                      </button>
                      <button
                        onClick={() => handleDeny(invite.email)}
                        className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                      >
                        Deny
                      </button>
                    </div>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}; 