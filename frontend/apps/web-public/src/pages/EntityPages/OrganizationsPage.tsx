import { Link } from 'react-router-dom';
import { Card, EmptyState, PageHeader } from '@radartcell/ui';
import { useOrganizations } from '@/api/queries';
import { OrgLogo } from '@/components/OrgLogo';
import { LoadingGrid } from './_LoadingGrid';

export function OrganizationsPage() {
  const q = useOrganizations();
  return (
    <>
      <PageHeader
        title="Организации"
        description="Компании и институты, стоящие за технологиями."
      />
      {q.isLoading ? (
        <LoadingGrid height={140} />
      ) : (q.data ?? []).length === 0 ? (
        <EmptyState title="Пусто" />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {(q.data ?? []).map((o) => (
            <Link key={o.id} to={`/organization/${encodeURIComponent(o.slug)}`}>
              <Card interactive>
                <div className="flex gap-3">
                  <OrgLogo
                    name={o.name}
                    logoUrl={o.logo_url}
                    className="h-14 w-14 rounded-xl text-lg"
                  />
                  <div className="min-w-0">
                    <h3 className="mb-1 text-base font-semibold">{o.name}</h3>
                    <div className="truncate text-sm text-ink-muted">
                      {o.headquarters || o.website || ''}
                    </div>
                    <div className="mt-1 text-xs text-ink-muted-2">
                      Технологий: {o.technologies_count}
                    </div>
                  </div>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </>
  );
}
