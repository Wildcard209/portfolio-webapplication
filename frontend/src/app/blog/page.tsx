import BlogCard from "@/app/components/BlogCard/BlogCard";

export default function Blog() {
    return (
        <div>
            <h1 className="fancy-font text-title">Latest Blogs</h1>
            <div className="row justify-content-center align-items-center">
                <div className="col-md-8"></div>
                    <BlogCard
                        name="Dynamic Blog Title"
                        description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
                        lastUpdated={new Date(new Date().getTime() - 3600000)}
                        imageId="1"
                        imageAlt="one"

                    />

                    <BlogCard
                        name="Dynamic Blog Title 2"
                        description="This is a blog description. It contains insightful content!"
                        lastUpdated={new Date(new Date().getTime() - 1500000)}
                        imageId="2"
                        imageAlt="two"
                    />

                    <BlogCard
                        name="Dynamic Blog Title 3"
                        description="This is a blog description. It contains insightful content!"
                        lastUpdated={new Date(new Date().getTime())}
                    />
            </div>
        </div>

    );
}
