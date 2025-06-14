"use client";

import HeroBanner from "@/app/components/HeroBanner/HeroBanner";
import Divider from "@/app/components/Divider/Divider";
import Blog from "@/app/blog/page";
import Projects from "@/app/projects/page";
import { useApi } from "@/lib/api/hooks/useApi";

export default function Home() {
  const { data, error, isLoading } = useApi<{ message?: string }>('/test', {
    revalidateTime: 60 
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      {data && (
        <div style={{ padding: '20px', margin: '20px', backgroundColor: '#f5f5f5', borderRadius: '8px' }}>
          <h2>API Response:</h2>
          <pre>{JSON.stringify(data, null, 2)}</pre>
        </div>
      )}
      <HeroBanner />
      <Projects />
      <Divider />
      <Blog />
    </div>
  );
}
