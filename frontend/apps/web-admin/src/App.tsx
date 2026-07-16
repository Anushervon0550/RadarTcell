import { Route, Routes } from 'react-router-dom';
import { Suspense, lazy } from 'react';
import { Spinner } from '@radartcell/ui';
import { AdminLayout } from './components/AdminLayout';
import { ProtectedRoute } from './auth/ProtectedRoute';
import { ResourceListPage } from './crud/ResourceListPage';
import {
  metricsResource,
  organizationsResource,
  sdgsResource,
  tagsResource,
  technologiesResource,
  trendsResource,
} from './crud/resources';

const LoginPage = lazy(() =>
  import('./pages/LoginPage').then((m) => ({ default: m.LoginPage })),
);
const DashboardPage = lazy(() =>
  import('./pages/DashboardPage').then((m) => ({ default: m.DashboardPage })),
);
const UsersPage = lazy(() =>
  import('./pages/UsersPage').then((m) => ({ default: m.UsersPage })),
);

function LoadingFallback() {
  return (
    <div className="flex h-[60vh] items-center justify-center">
      <Spinner />
    </div>
  );
}

export function App() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          element={
            <ProtectedRoute>
              <AdminLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<DashboardPage />} />
          <Route path="/technologies" element={<ResourceListPage resource={technologiesResource} />} />
          <Route path="/trends" element={<ResourceListPage resource={trendsResource} />} />
          <Route path="/tags" element={<ResourceListPage resource={tagsResource} />} />
          <Route path="/sdgs" element={<ResourceListPage resource={sdgsResource} />} />
          <Route path="/organizations" element={<ResourceListPage resource={organizationsResource} />} />
          <Route path="/metrics" element={<ResourceListPage resource={metricsResource} />} />
          <Route path="/users" element={<UsersPage />} />
        </Route>
      </Routes>
    </Suspense>
  );
}
