import { notFound } from "next/navigation";

export default async function AdminPanel({ searchParams }: { searchParams: { t: string } }) {
    const validAdminToken = process.env.ADMIN_TOKEN;

    const params = await searchParams;
    const { t } = params;

    if (!validAdminToken || t !== validAdminToken) {
        notFound();
    }

    return (
        <div style={{ textAlign: "center", marginTop: "50px" }}>
            <h1>Admin Panel</h1>
            <p>You successfully accessed the admin panel!</p>
        </div>
    );
}
