import BlogCard from "@/app/components/BlogCard/BlogCard";

export default function Blog() {
  return (
    <div className="container mt-4">
      <h1 className="fancy-font text-title text-center mb-4">Latest Blogs</h1>
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
