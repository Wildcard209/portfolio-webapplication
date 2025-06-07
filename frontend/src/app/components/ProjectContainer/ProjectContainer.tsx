'use client';

import ProjectCard from "../ProjectCard/ProjectCard";
import styles from "./ProjectContainer.module.scss";

export type ProjectData = {
    title: string;
    description: string;
    lastUpdated: Date;
    imageId: string;
    imageAlt?: string;
};

type ProjectContainerProps = {
    projects: ProjectData[];
};

const ProjectContainer: React.FC<ProjectContainerProps> = ({ projects }) => {
    return (
        <div className={styles["projects-grid"]}>
            {projects.map((project, index) => (
                <ProjectCard
                    key={`${project.imageId}-${index}`}
                    title={project.title}
                    description={project.description}
                    lastUpdated={project.lastUpdated}
                    imageId={project.imageId}
                    imageAlt={project.imageAlt}
                />
            ))}
        </div>
    );
};

export default ProjectContainer;
