import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { InvitesTable } from '../InvitesTable';

// Mock fetch
const mockFetch = jest.fn();
global.fetch = mockFetch;

// Mock clipboard API
Object.defineProperty(navigator, 'clipboard', {
  value: {
    writeText: jest.fn().mockResolvedValue(undefined),
  },
  writable: true,
});

describe('InvitesTable', () => {
  const mockInvites = [
    {
      name: 'John Doe',
      role: 'Developer',
      email: 'john@example.com',
      company: 'Company A',
      yearsExperience: '5',
      reasons: 'Great experience',
      source: 'LinkedIn',
    },
    {
      name: 'Jane Smith',
      role: 'Designer',
      email: 'jane@example.com',
      company: 'Company B',
      yearsExperience: '3',
      reasons: 'Strong portfolio',
      source: 'Referral',
    },
  ];

  beforeEach(() => {
    mockFetch.mockClear();
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockInvites),
    });
  });

  it('loads and displays invites', async () => {
    render(<InvitesTable />);
    
    // Check loading state
    expect(screen.getByText('Loading...')).toBeInTheDocument();
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    });
  });

  it('handles approve and deny actions', async () => {
    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Switch to approved tab
    const approvedTab = screen.getByText('Approved');
    fireEvent.click(approvedTab);

    // Check if invite appears in approved list
    expect(screen.getByText('John Doe')).toBeInTheDocument();

    // Switch back to pending
    const pendingTab = screen.getByText('Pending');
    fireEvent.click(pendingTab);

    // Deny second invite
    const denyButton = screen.getAllByText('Deny')[1];
    fireEvent.click(denyButton);

    // Switch to denied tab
    const deniedTab = screen.getByText('Denied');
    fireEvent.click(deniedTab);

    // Check if invite appears in denied list
    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
  });

  it('handles undo action', async () => {
    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Switch to approved tab
    const approvedTab = screen.getByText('Approved');
    fireEvent.click(approvedTab);

    // Undo the approval
    const undoButton = screen.getByText('Undo');
    fireEvent.click(undoButton);

    // Switch back to pending
    const pendingTab = screen.getByText('Pending');
    fireEvent.click(pendingTab);

    // Check if invite is back in pending list
    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });

  it('handles marking invites as sent', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockInvites),
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ status: 'success' }),
    });

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Click Invites Sent button
    const sentButton = screen.getByText('Invites Sent');
    fireEvent.click(sentButton);

    // Check loading state
    expect(screen.getByText('Updating...')).toBeInTheDocument();

    // Wait for success state
    await waitFor(() => {
      expect(screen.getByText('✓ Invites Updated!')).toBeInTheDocument();
    });
  });

  it('handles marking invites as rejected', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockInvites),
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ status: 'success' }),
    });

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Deny first invite
    const denyButton = screen.getAllByText('Deny')[0];
    fireEvent.click(denyButton);

    // Switch to denied tab
    const deniedTab = screen.getByText('Denied');
    fireEvent.click(deniedTab);

    // Click Reject All button
    const rejectButton = screen.getByText('Reject All');
    fireEvent.click(rejectButton);

    // Check loading state
    expect(screen.getByText('Updating...')).toBeInTheDocument();

    // Wait for success state
    await waitFor(() => {
      expect(screen.getByText('✓ Rejected!')).toBeInTheDocument();
    });
  });

  it('handles error states', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockInvites),
    }).mockRejectedValueOnce(new Error('Failed to update'));

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Click Invites Sent button
    const sentButton = screen.getByText('Invites Sent');
    fireEvent.click(sentButton);

    // Wait for error state
    await waitFor(() => {
      expect(screen.getByText('Error - Try Again')).toBeInTheDocument();
    });
  });

  it('handles copying emails to clipboard', async () => {
    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Click Copy Emails button
    const copyButton = screen.getByText('Copy Emails');
    fireEvent.click(copyButton);

    // Check if clipboard API was called
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('john@example.com');

    // Check success message
    expect(screen.getByText('✓ Copied to clipboard!')).toBeInTheDocument();
  });
}); 