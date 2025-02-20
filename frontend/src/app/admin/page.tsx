import { notFound } from 'next/navigation';
import LoginFlow from '../components/Admin/LoginFlow';

export default async function AdminTokenValidation({searchParams}: { searchParams: { t: string } }) {
    const validAdminToken = process.env.ADMIN_TOKEN;
    const token = await searchParams;

    if (!validAdminToken || token.t !== validAdminToken) {
        notFound();
    }

    return (
        <div style={{textAlign: 'center', marginTop: '50px'}}>
            <h1>Welcome to Admin Panel</h1>
            <LoginFlow/>
        </div>
    );
}

