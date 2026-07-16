import { Link } from 'react-router-dom';
import { Button, EmptyState, PageHeader } from '@radartcell/ui';

export function NotFoundPage() {
  return (
    <>
      <PageHeader title="404" description="Страница не найдена" />
      <EmptyState
        title="Ничего нет по этому адресу"
        description="Проверьте URL или вернитесь на главную."
        action={
          <Link to="/">
            <Button variant="primary">На главную</Button>
          </Link>
        }
      />
    </>
  );
}
