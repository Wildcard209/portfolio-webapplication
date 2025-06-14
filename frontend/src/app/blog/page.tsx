"use client";

import BlogCard from "@/app/components/BlogCard/BlogCard";
import { useAuth } from "@/lib/hooks/useAuth";

export default function Blog() {
  const { isAuthenticated, loading } = useAuth();

  const handleAddBlog = () => {
    // TODO: Implement add blog functionality
    alert("Add Blog functionality will be implemented here!");
  };

  return (
    <div className="container mt-4">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h1 className="fancy-font text-title text-center flex-grow-1">Latest Blogs</h1>
        
        {/* Add Blog button - only show if authenticated */}
        {!loading && isAuthenticated && (
          <button
            onClick={handleAddBlog}
            className="btn btn-primary ms-3"
            style={{
              backgroundColor: '#007bff',
              borderColor: '#007bff',
              padding: '10px 20px',
              fontSize: '14px',
              fontWeight: '500',
              borderRadius: '6px',
              minWidth: '120px'
            }}
          >
            Add Blog
          </button>
        )}
      </div>

      <div className="row justify-content-center">
        <div className="col-md-10">
          <div className="d-flex flex-column gap-4">
            <BlogCard
              title="Dynamic Blog Title"
              description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
              lastUpdated={new Date(new Date().getTime() - 3600000)}
              imageId="1"
              imageAlt="one"
            />

            <BlogCard
              title="Dynamic Blog Title 2"
              description="This is a blog description. It contains insightful content!"
              lastUpdated={new Date(new Date().getTime() - 1500000)}
              imageId="2"
              imageAlt="two"
            />

            <BlogCard
              title="Dynamic Blog Title 3"
              description="This is a blog description. It contains insightful content!"
              lastUpdated={new Date(new Date().getTime())}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
