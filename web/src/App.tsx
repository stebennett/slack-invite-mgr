import { ErrorBoundary } from './components/ErrorBoundary';
import { InvitesTable } from './components/InvitesTable';

export const App = () => {
  return (
    <ErrorBoundary>
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Outstanding Invites</h1>
        <InvitesTable />
      </div>
    </ErrorBoundary>
  );
};
