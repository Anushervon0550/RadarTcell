import { Route, Routes, useLocation } from 'react-router-dom';
import { Suspense, lazy, useEffect } from 'react';
import { Layout } from './components/Layout';
import { ErrorBoundary } from './components/ErrorBoundary';
import { Spinner } from '@radartcell/ui';

const RadarPage = lazy(() =>
  import('./pages/radar/RadarPage').then((m) => ({ default: m.RadarPage })),
);
const CatalogPage = lazy(() =>
  import('./pages/CatalogPage').then((m) => ({ default: m.CatalogPage })),
);
const TechnologyPage = lazy(() =>
  import('./pages/TechnologyPage').then((m) => ({ default: m.TechnologyPage })),
);
const NotFoundPage = lazy(() =>
  import('./pages/NotFoundPage').then((m) => ({ default: m.NotFoundPage })),
);
const TrendsPage = lazy(() =>
  import('./pages/EntityPages').then((m) => ({ default: m.TrendsPage })),
);
const TagsPage = lazy(() =>
  import('./pages/EntityPages').then((m) => ({ default: m.TagsPage })),
);
const SDGsPage = lazy(() =>
  import('./pages/EntityPages').then((m) => ({ default: m.SDGsPage })),
);
const OrganizationsPage = lazy(() =>
  import('./pages/EntityPages').then((m) => ({ default: m.OrganizationsPage })),
);
const ByEntityPage = lazy(() =>
  import('./pages/ByEntityPage').then((m) => ({ default: m.ByEntityPage })),
);

function ScrollToTop() {
  const location = useLocation();
  useEffect(() => {
    if (location.pathname.startsWith('/technology/')) return;
    window.scrollTo({ top: 0, behavior: 'instant' as ScrollBehavior });
  }, [location.pathname]);
  return null;
}

function LoadingFallback() {
  return (
    <div className="flex h-[60vh] items-center justify-center">
      <Spinner />
    </div>
  );
}

export function App() {
  return (
    <ErrorBoundary>
      <ScrollToTop />
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<RadarPage />} />
            <Route path="/radar/:trend" element={<RadarPage />} />
            <Route path="/catalog" element={<CatalogPage />} />
            <Route path="/technology/:slug" element={<TechnologyPage />} />

            <Route path="/trends" element={<TrendsPage />} />
            <Route path="/tags" element={<TagsPage />} />
            <Route path="/sdgs" element={<SDGsPage />} />
            <Route path="/organizations" element={<OrganizationsPage />} />

            <Route path="/trend/:value" element={<ByEntityPage kind="trend" />} />
            <Route path="/tag/:value" element={<ByEntityPage kind="tag" />} />
            <Route path="/sdg/:value" element={<ByEntityPage kind="sdg" />} />
            <Route path="/organization/:value" element={<ByEntityPage kind="organization" />} />

            <Route path="*" element={<NotFoundPage />} />
          </Route>
        </Routes>
      </Suspense>
    </ErrorBoundary>
  );
}
