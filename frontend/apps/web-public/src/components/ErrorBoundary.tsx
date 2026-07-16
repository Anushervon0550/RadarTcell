import { Component, type ErrorInfo, type ReactNode } from 'react';
import { Button, Card, EmptyState } from '@radartcell/ui';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    // eslint-disable-next-line no-console
    console.error('UI ErrorBoundary caught error', error, info);
  }

  reset = () => this.setState({ error: null });

  render() {
    const { error } = this.state;
    if (!error) return this.props.children;
    if (this.props.fallback) return this.props.fallback;
    return (
      <Card className="mx-auto my-16 max-w-lg">
        <EmptyState
          title="Что-то пошло не так"
          description={error.message}
          action={
            <Button variant="primary" onClick={this.reset}>
              Попробовать снова
            </Button>
          }
        />
      </Card>
    );
  }
}
