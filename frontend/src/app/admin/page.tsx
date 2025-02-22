import { notFound } from 'next/navigation';
import LoginFlow from '../components/Admin/LoginFlow';

interface AdminProps {
    searchParams: Promise<{ t: string }>;
}

export default async function AdminTokenValidation({ searchParams,}: AdminProps) {
    const validAdminToken = process.env.ADMIN_TOKEN;
    const params = await searchParams;

    if (!validAdminToken || params.t !== validAdminToken) {
        notFound();
    }

    return (
        <div style={{textAlign: 'center', marginTop: '50px'}}>
            <LoginFlow/>
        </div>
    );
}