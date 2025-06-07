import ProjectContainer from "@/app/components/ProjectContainer/ProjectContainer";
import type { ProjectData } from "@/app/components/ProjectContainer/ProjectContainer";

const projectsData: ProjectData[] = [
    {
        title: "Dynamic Blog Title",
        description: "This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!This is a blog description. It contains insightful content!",
        lastUpdated: new Date(new Date().getTime() - 3600000),
        imageId: "1",
        imageAlt: "one"
    },
    {
        title: "Dynamic Blog Title 2",
        description: "This is a blog description. It contains insightful content!",
        lastUpdated: new Date(new Date().getTime() - 7200000),
        imageId: "2",
        imageAlt: "two"
    },
    {
        title: "Dynamic Blog Title 3",
        description: "This is a blog description. It contains insightful content!",
        lastUpdated: new Date(new Date().getTime() - 10800000),
        imageId: "3",
        imageAlt: "three"
    },
    {
        title: "Dynamic Blog Title 4",
        description: "This is a blog description. It contains insightful content!",
        lastUpdated: new Date(new Date().getTime() - 14400000),
        imageId: "4",
        imageAlt: "four"
    }
];

export default function Projects() {
    return (
        <div className="container mt-4">
            <h1 className="fancy-font text-title text-center mb-4">Most Recent Projects</h1>
            <div className="d-flex justify-content-center">
                <div className="col-md-10">
                    <ProjectContainer projects={projectsData} />
                </div>
            </div>
        </div>
    );
}
