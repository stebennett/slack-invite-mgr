type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
  level: LogLevel;
  message: string;
  context?: Record<string, unknown>;
}

const getApiUrl = (): string => {
  return window.APP_CONFIG?.API_URL || 'http://localhost:8080';
};

const sendToBackend = async (entry: LogEntry): Promise<void> => {
  try {
    const apiUrl = getApiUrl();
    await fetch(`${apiUrl}/api/logs`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(entry),
    });
  } catch {
    // Silently fail - we don't want logging failures to break the app
  }
};

const log = (level: LogLevel, message: string, context?: Record<string, unknown>): void => {
  const entry: LogEntry = { level, message, context };

  // Always output to console for development
  const consoleOutput = JSON.stringify({ time: new Date().toISOString(), ...entry });
  switch (level) {
    case 'error':
      console.error(consoleOutput);
      break;
    case 'warn':
      console.warn(consoleOutput);
      break;
    case 'debug':
      console.debug(consoleOutput);
      break;
    default:
      console.log(consoleOutput);
  }

  // Send error and warn levels to backend for Loki ingestion
  if (level === 'error' || level === 'warn') {
    sendToBackend(entry);
  }
};

export const logger = {
  debug: (message: string, context?: Record<string, unknown>) => log('debug', message, context),
  info: (message: string, context?: Record<string, unknown>) => log('info', message, context),
  warn: (message: string, context?: Record<string, unknown>) => log('warn', message, context),
  error: (message: string, context?: Record<string, unknown>) => log('error', message, context),
};
