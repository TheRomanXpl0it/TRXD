import * as React from "react";

interface ErrorBoundaryProps {
  fallback: React.ReactNode;
  children?: React.ReactNode;
  onError?(error: Error, info: React.ErrorInfo): void; // optional callback
}

interface ErrorBoundaryState {
  hasError: boolean;
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false };

  static getDerivedStateFromError(_error: Error) {
    // Update state so the next render shows the fallback UI.
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // Log the error + component stack; DO NOT call React internals.
    if (this.props.onError) this.props.onError(error, info);

    // Keep console noise minimal in production if you want
    // eslint-disable-next-line no-console
    console.error("ErrorBoundary caught an error:", error);
    // eslint-disable-next-line no-console
    console.error(info.componentStack);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children ?? null;
  }
}

export { ErrorBoundary };
