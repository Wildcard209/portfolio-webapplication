import ProjectCard from "@/app/components/ProjectCard/ProjectCard";

export default function Projects() {
    return (
        <div>
            <h1 className="fancy-font text-title">Most Recent Projects</h1>
            <div className="row justify-content-center align-items-center">
            <ProjectCard
                title="Dynamic Blog Title"
                description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
                lastUpdated={new Date(new Date().getTime() - 3600000)}
                imageId="1"
                imageAlt="one"
            />
                <ProjectCard
                    title="Dynamic Blog Title"
                    description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
                    lastUpdated={new Date(new Date().getTime() - 3600000)}
                    imageId="1"
                    imageAlt="one"
                />
                <ProjectCard
                    title="Dynamic Blog Title"
                    description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
                    lastUpdated={new Date(new Date().getTime() - 3600000)}
                    imageId="1"
                    imageAlt="one"
                />
                <ProjectCard
                    title="Dynamic Blog Title"
                    description="This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!"
                    lastUpdated={new Date(new Date().getTime() - 3600000)}
                    imageId="1"
                    imageAlt="one"
                />
            </div>
        </div>
    );
}
