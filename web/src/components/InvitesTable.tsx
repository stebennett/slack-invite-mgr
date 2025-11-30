import React, { useEffect, useState } from 'react';

interface Invite {
  name: string;
  role: string;
  email: string;
  company: string;
  yearsExperience: string;
  reasons: string;
  source: string;
  status?: 'pending' | 'approved' | 'denied' | 'sent';
}

type WorkflowStep = 'screening' | 'send-invites' | 'slack-preparation' | 'mark-denied' | 'complete';

export const InvitesTable: React.FC = () => {
  const [invites, setInvites] = useState<Invite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'pending' | 'approved' | 'denied'>('pending');
  const [currentStep, setCurrentStep] = useState<WorkflowStep>('screening');
  const [copySuccess, setCopySuccess] = useState(false);
  const [updateStatus, setUpdateStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');

  // Get the base URL for API calls, respecting PUBLIC_URL for subpath deployments
  // Uses runtime config (window.APP_CONFIG) which is generated at container startup
  const getApiUrl = (endpoint: string) => {
    const publicUrl = (window as any).APP_CONFIG?.PUBLIC_URL || '';
    return `${publicUrl}/api${endpoint}`;
  };

  const fetchInvites = async () => {
    try {
      const response = await fetch(getApiUrl('/invites'));
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

  useEffect(() => {
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

  const handleUndo = (email: string) => {
    setInvites(invites.map(invite => 
      invite.email === email ? { ...invite, status: 'pending' } : invite
    ));
  };

  const handleCopyEmails = async () => {
    const emailList = approvedInvites.map(invite => invite.email).join(', ');
    
    try {
      // Try using the modern clipboard API first
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(emailList);
        setCopySuccess(true);
        setTimeout(() => setCopySuccess(false), 2000);
      } else {
        // Fallback for older browsers or non-secure contexts
        const textArea = document.createElement('textarea');
        textArea.value = emailList;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        textArea.style.top = '-999999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        
        try {
          const successful = document.execCommand('copy');
          if (successful) {
            setCopySuccess(true);
            setTimeout(() => setCopySuccess(false), 2000);
          } else {
            console.error('Failed to copy emails using fallback method');
          }
        } finally {
          document.body.removeChild(textArea);
        }
      }
    } catch (err) {
      console.error('Failed to copy emails:', err);
      // Optionally show user-friendly error message
      alert('Failed to copy emails to clipboard. Please copy manually.');
    }
  };

  const handleInvitesSent = async () => {
    setUpdateStatus('loading');
    try {
      const response = await fetch(getApiUrl('/invites'), {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          emails: approvedInvites.map(invite => invite.email),
          status: 'sent'
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to update invite statuses');
      }

      setUpdateStatus('success');
      // Update the status in our local state
      setInvites(prevInvites => 
        prevInvites.map(invite => 
          approvedInvites.some(approved => approved.email === invite.email)
            ? { ...invite, status: 'sent' }
            : invite
        )
      );
      // Transition to the denied list
      setCurrentStep('mark-denied');
      setUpdateStatus('idle');
    } catch (err) {
      console.error('Failed to update invite statuses:', err);
      setUpdateStatus('error');
    }
  };

  const handleDeniedConfirmed = async () => {
    setUpdateStatus('loading');
    try {
      const response = await fetch(getApiUrl('/invites'), {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          emails: deniedInvites.map(invite => invite.email),
          status: 'denied'
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to update invite statuses');
      }

      setUpdateStatus('success');
      setCurrentStep('complete');
    } catch (err) {
      console.error('Failed to update invite statuses:', err);
      setUpdateStatus('error');
    }
  };

  const handleExit = async () => {
    // Reset all approved and denied invites back to pending
    setInvites(invites.map(invite => ({
      ...invite,
      status: invite.status === 'approved' || invite.status === 'denied' ? 'pending' : invite.status
    })));

    // Reload pending invites from backend
    try {
      const response = await fetch(getApiUrl('/invites'));
      if (!response.ok) {
        throw new Error('Failed to fetch invites');
      }
      const data = await response.json();
      setInvites(data.map((invite: Invite) => ({ ...invite, status: 'pending' })));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }

    // Reset workflow state
    setCurrentStep('screening');
    setActiveTab('pending');
  };

  const filteredInvites = invites.filter(invite => 
    activeTab === 'denied' 
      ? invite.status === 'denied'
      : activeTab === 'approved'
        ? invite.status === 'approved'
        : invite.status === activeTab
  );
  const approvedInvites = invites.filter(invite => invite.status === 'approved');
  const deniedInvites = invites.filter(invite => invite.status === 'denied');

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (currentStep === 'send-invites') {
    return (
      <div>
        <div className="mb-6">
          <button
            onClick={handleExit}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Exit
          </button>
        </div>
        <div className="mb-4">
          <h2 className="text-2xl font-bold mb-4">Send Invites</h2>
          <p className="text-gray-600 mb-4">
            Review the following {approvedInvites.length} approved invite(s):
          </p>
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
                <th className="px-6 py-3 border-b text-left">Actions</th>
              </tr>
            </thead>
            <tbody>
              {approvedInvites.map((invite, index) => (
                <tr key={index} className="hover:bg-gray-50">
                  <td className="px-6 py-4 border-b">{invite.name}</td>
                  <td className="px-6 py-4 border-b">{invite.role}</td>
                  <td className="px-6 py-4 border-b">{invite.email}</td>
                  <td className="px-6 py-4 border-b">{invite.company}</td>
                  <td className="px-6 py-4 border-b">{invite.yearsExperience}</td>
                  <td className="px-6 py-4 border-b">{invite.reasons}</td>
                  <td className="px-6 py-4 border-b">{invite.source}</td>
                  <td className="px-6 py-4 border-b">
                    <button
                      onClick={() => handleUndo(invite.email)}
                      className="px-3 py-1 bg-gray-500 text-white rounded hover:bg-gray-600"
                    >
                      Undo
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="mt-6 flex justify-end">
          <button
            onClick={() => {
              setUpdateStatus('idle');
              setCurrentStep('slack-preparation');
            }}
            disabled={approvedInvites.length === 0}
            className={`px-6 py-3 rounded text-white text-lg ${
              approvedInvites.length === 0
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-blue-500 hover:bg-blue-600'
            }`}
          >
            Confirm ({approvedInvites.length} approved)
          </button>
        </div>
      </div>
    );
  }

  if (currentStep === 'slack-preparation') {
    return (
      <div>
        <div className="mb-6">
          <button
            onClick={handleExit}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Exit
          </button>
        </div>
        <div className="mb-4">
          <h2 className="text-2xl font-bold mb-4">Slack Invites Preparation</h2>
          <p className="text-gray-600 mb-4">
            Copy the following email addresses to invite to Slack:
          </p>
        </div>
        <div className="bg-gray-100 p-4 rounded-lg mb-4">
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm text-gray-600">
              {approvedInvites.length} email(s)
            </span>
            <div className="flex items-center space-x-2">
              {copySuccess && (
                <span className="text-green-600 text-sm">
                  ✓ Copied to clipboard!
                </span>
              )}
              <button
                onClick={handleCopyEmails}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 flex items-center space-x-2"
              >
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
                  />
                </svg>
                <span>Copy Emails</span>
              </button>
            </div>
          </div>
          <div className="bg-white p-4 rounded border border-gray-300">
            {approvedInvites.map((invite, index) => (
              <span key={index}>
                {invite.email}
                {index < approvedInvites.length - 1 ? ', ' : ''}
              </span>
            ))}
          </div>
        </div>
        <div className="mt-6 flex justify-end">
          <button
            onClick={handleInvitesSent}
            disabled={updateStatus === 'loading'}
            className={`px-6 py-3 rounded text-white text-lg ${
              updateStatus === 'loading'
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-green-500 hover:bg-green-600'
            }`}
          >
            {updateStatus === 'loading' ? (
              'Updating...'
            ) : updateStatus === 'success' ? (
              '✓ Invites Updated!'
            ) : updateStatus === 'error' ? (
              'Error - Try Again'
            ) : (
              'Confirm Invites Sent'
            )}
          </button>
        </div>
      </div>
    );
  }

  if (currentStep === 'mark-denied') {
    return (
      <div>
        <div className="mb-6">
          <button
            onClick={handleExit}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Exit
          </button>
        </div>
        <div className="mb-4">
          <h2 className="text-2xl font-bold mb-4">Mark Denied Invites</h2>
          <p className="text-gray-600 mb-4">
            Review the following {deniedInvites.length} denied invite(s):
          </p>
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
              </tr>
            </thead>
            <tbody>
              {deniedInvites.map((invite, index) => (
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
        <div className="mt-6 flex justify-end">
          <button
            onClick={handleDeniedConfirmed}
            disabled={updateStatus === 'loading'}
            className={`px-6 py-3 rounded text-white text-lg ${
              updateStatus === 'loading'
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-red-500 hover:bg-red-600'
            }`}
          >
            {updateStatus === 'loading' ? (
              'Updating...'
            ) : updateStatus === 'success' ? (
              '✓ Denied!'
            ) : updateStatus === 'error' ? (
              'Error - Try Again'
            ) : (
              'Confirm'
            )}
          </button>
        </div>
      </div>
    );
  }

  if (currentStep === 'complete') {
    const sentCount = invites.filter(invite => invite.status === 'sent').length;
    const deniedCount = invites.filter(invite => invite.status === 'denied').length;
    
    return (
      <div>
        <div className="mb-6">
          <button
            onClick={handleExit}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Exit
          </button>
        </div>
        <div className="mb-4">
          <h2 className="text-2xl font-bold mb-4">Process Complete</h2>
          <div className="bg-green-100 p-4 rounded-lg mb-4">
            <h3 className="text-lg font-semibold mb-2">Summary</h3>
            <p className="mb-2">✓ {sentCount} invite(s) sent</p>
            <p>✓ {deniedCount} invite(s) denied</p>
          </div>
        </div>
      </div>
    );
  }

  if (invites.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-600 text-lg">No invite requests</p>
      </div>
    );
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
              <th className="px-6 py-3 border-b text-left">Actions</th>
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
                <td className="px-6 py-4 border-b">
                  {activeTab === 'pending' ? (
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
                  ) : (
                    <button
                      onClick={() => handleUndo(invite.email)}
                      className="px-3 py-1 bg-gray-500 text-white rounded hover:bg-gray-600"
                    >
                      Undo
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {activeTab === 'pending' && (
        <div className="mt-6 flex justify-end">
          <button
            onClick={() => setCurrentStep('send-invites')}
            disabled={approvedInvites.length === 0}
            className={`px-6 py-3 rounded text-white text-lg ${
              approvedInvites.length === 0
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-blue-500 hover:bg-blue-600'
            }`}
          >
            Next ({approvedInvites.length} approved)
          </button>
        </div>
      )}
    </div>
  );
}; 