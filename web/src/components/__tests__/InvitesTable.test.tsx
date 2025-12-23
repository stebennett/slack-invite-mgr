import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { InvitesTable } from '../InvitesTable';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock clipboard API
Object.defineProperty(navigator, 'clipboard', {
  value: {
    writeText: vi.fn().mockResolvedValue(undefined),
  },
  writable: true,
});

// Mock window.isSecureContext
Object.defineProperty(window, 'isSecureContext', {
  value: true,
  writable: true,
});

// Mock document.execCommand
document.execCommand = vi.fn();

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

    // Wait for and check if invite appears in approved list
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Switch back to pending
    const pendingTab = screen.getByText('Pending');
    fireEvent.click(pendingTab);

    // Wait for pending list to update
    await waitFor(() => {
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    });

    // Deny second invite
    const denyButton = screen.getAllByText('Deny')[0];
    fireEvent.click(denyButton);

    // Switch to denied tab
    const deniedTab = screen.getByText('Denied');
    fireEvent.click(deniedTab);

    // Wait for and check if invite appears in denied list
    await waitFor(() => {
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    });
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

  it('handles the complete workflow', async () => {
    // Mock API responses for the complete workflow
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockInvites),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ status: 'success' }),
      })
      .mockResolvedValueOnce({
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

    // Click Next button to go to Send Invites
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Wait for Send Invites screen
    await waitFor(() => {
      expect(screen.getByText('Send Invites')).toBeInTheDocument();
    });

    // Click Confirm button to go to Slack Preparation
    const confirmButton = screen.getByText(/Confirm/);
    fireEvent.click(confirmButton);

    // Wait for Slack Preparation screen
    await waitFor(() => {
      expect(screen.getByText('Slack Invites Preparation')).toBeInTheDocument();
    });

    // Click Copy Emails button
    const copyButton = screen.getByText('Copy Emails');
    fireEvent.click(copyButton);

    // Check if clipboard API was called
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('john@example.com');

    // Click Confirm Invites Sent button
    const sentButton = screen.getByText('Confirm Invites Sent');
    fireEvent.click(sentButton);

    // Wait for Mark Denied screen
    await waitFor(() => {
      expect(screen.getByText('Mark Denied Invites')).toBeInTheDocument();
    });

    // Click Confirm button
    const deniedButton = screen.getByText('Confirm');
    fireEvent.click(deniedButton);

    // Wait for Complete screen
    await waitFor(() => {
      expect(screen.getByText('Process Complete')).toBeInTheDocument();
    });

    // Check summary
    expect(screen.getByText('✓ 1 invite(s) sent')).toBeInTheDocument();
    expect(screen.getByText('✓ 0 invite(s) denied')).toBeInTheDocument();
  });

  it('handles exit from workflow', async () => {
    // Mock the initial fetch and the reload fetch
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockInvites),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockInvites),
      });

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Deny second invite
    const denyButton = screen.getAllByText('Deny')[0];
    fireEvent.click(denyButton);

    // Click Next button to go to Send Invites
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Wait for Send Invites screen
    await waitFor(() => {
      expect(screen.getByText('Send Invites')).toBeInTheDocument();
    });

    // Click Exit button
    const exitButton = screen.getByText('Exit');
    fireEvent.click(exitButton);

    // Wait for screening screen and verify state
    await waitFor(() => {
      // Check that we're back on the pending tab
      const pendingTab = screen.getByText('Pending');
      expect(pendingTab).toHaveClass('bg-blue-500');

      // Check that both invites are back in pending state
      expect(screen.getByText('John Doe')).toBeInTheDocument();
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();

      // Verify that the backend was called again to reload invites
      expect(mockFetch).toHaveBeenCalledTimes(2);
    });

    // Switch to approved tab and verify it's empty
    const approvedTab = screen.getByText('Approved');
    fireEvent.click(approvedTab);
    expect(screen.queryByText('John Doe')).not.toBeInTheDocument();

    // Switch to denied tab and verify it's empty
    const deniedTab = screen.getByText('Denied');
    fireEvent.click(deniedTab);
    expect(screen.queryByText('Jane Smith')).not.toBeInTheDocument();
  });

  it('handles error states', async () => {
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockInvites),
      })
      .mockRejectedValueOnce(new Error('Failed to update'));

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button to go to Send Invites
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Wait for Send Invites screen
    await waitFor(() => {
      expect(screen.getByText('Send Invites')).toBeInTheDocument();
    });

    // Click Confirm button to go to Slack Preparation
    const confirmButton = screen.getByText(/Confirm/);
    fireEvent.click(confirmButton);

    // Wait for Slack Preparation screen
    await waitFor(() => {
      expect(screen.getByText('Slack Invites Preparation')).toBeInTheDocument();
    });

    // Click Confirm Invites Sent button
    const sentButton = screen.getByText('Confirm Invites Sent');
    fireEvent.click(sentButton);

    // Wait for error state
    await waitFor(() => {
      expect(screen.getByText('Error - Try Again')).toBeInTheDocument();
    });
  });

  it('displays message when there are no invites', async () => {
    // Mock empty invites response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    render(<InvitesTable />);
    
    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
    });

    // Check that the "No invite requests" message is displayed
    expect(screen.getByText('No invite requests')).toBeInTheDocument();
    
    // Verify that no table is rendered
    expect(screen.queryByRole('table')).not.toBeInTheDocument();
  });

  it('handles API error when fetching invites', async () => {
    // Mock failed API response
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
    });

    render(<InvitesTable />);
    
    // Wait for error state
    await waitFor(() => {
      expect(screen.getByText('Error: Failed to fetch invites')).toBeInTheDocument();
    });
    
    // Verify that no table is rendered
    expect(screen.queryByRole('table')).not.toBeInTheDocument();
  });

  it('handles clipboard copy with fallback method', async () => {
    // Mock the scenario where modern clipboard API is not available
    const originalClipboard = navigator.clipboard;
    const originalIsSecureContext = window.isSecureContext;
    
    // Mock document.execCommand
    const execCommandMock = vi.fn().mockReturnValue(true);
    document.execCommand = execCommandMock;
    
    // Test when clipboard API is not available
    Object.defineProperty(navigator, 'clipboard', {
      value: undefined,
      writable: true,
    });

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button to go to Send Invites
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Click Confirm button to go to Slack Preparation
    const confirmButton = screen.getByText(/Confirm/);
    fireEvent.click(confirmButton);

    // Wait for Slack Preparation screen
    await waitFor(() => {
      expect(screen.getByText('Slack Invites Preparation')).toBeInTheDocument();
    });

    // Click Copy Emails button
    const copyButton = screen.getByText('Copy Emails');
    fireEvent.click(copyButton);

    // Check if fallback method was used
    expect(execCommandMock).toHaveBeenCalledWith('copy');
    
    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText('✓ Copied to clipboard!')).toBeInTheDocument();
    });

    // Restore original clipboard
    Object.defineProperty(navigator, 'clipboard', {
      value: originalClipboard,
      writable: true,
    });
    Object.defineProperty(window, 'isSecureContext', {
      value: originalIsSecureContext,
      writable: true,
    });
  });

  it('handles clipboard copy error gracefully', async () => {
    // Mock clipboard API to throw an error
    const alertMock = vi.spyOn(window, 'alert').mockImplementation(() => {});
    navigator.clipboard.writeText = vi.fn().mockRejectedValue(new Error('Permission denied'));
    
    // Also make isSecureContext true to ensure we try the modern API
    Object.defineProperty(window, 'isSecureContext', {
      value: true,
      writable: true,
    });

    render(<InvitesTable />);
    
    // Wait for invites to load
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // Approve first invite
    const approveButton = screen.getAllByText('Approve')[0];
    fireEvent.click(approveButton);

    // Click Next button to go to Send Invites
    const nextButton = screen.getByText(/Next/);
    fireEvent.click(nextButton);

    // Click Confirm button to go to Slack Preparation
    const confirmButton = screen.getByText(/Confirm/);
    fireEvent.click(confirmButton);

    // Wait for Slack Preparation screen
    await waitFor(() => {
      expect(screen.getByText('Slack Invites Preparation')).toBeInTheDocument();
    });

    // Click Copy Emails button
    const copyButton = screen.getByText('Copy Emails');
    fireEvent.click(copyButton);

    // Wait for the alert to be called
    await waitFor(() => {
      expect(alertMock).toHaveBeenCalledWith('Failed to copy emails to clipboard. Please copy manually.');
    });

    alertMock.mockRestore();
  });
}); 