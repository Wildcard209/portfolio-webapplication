import { notFound } from 'next/navigation';
import LoginFlow from '../components/Admin/LoginFlow';
import { ReactElement } from 'react';

interface AdminProps {
  searchParams: Promise<{ t: string }>;
}

export default async function AdminTokenValidation({
  searchParams,
}: AdminProps): Promise<ReactElement> {
  const validAdminToken = process.env.ADMIN_TOKEN;
  const params = await searchParams;

  if (!validAdminToken || params.t !== validAdminToken) {
    notFound();
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        backgroundColor: '#f8f9fa',
        paddingTop: '50px',
      }}
    >
      <LoginFlow adminToken={params.t} />
    </div>
  );
}
