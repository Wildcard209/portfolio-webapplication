'use client'

import Image from "next/image";

const imageLoader = ({ width }: { width: number }) => {
    return `https://picsum.photos/${width}/180`
}

export default function ProjectCard() {
    return (
        <div>
            <h1 className="fancy-font text-title">Most Recent Projects</h1>
            <div className="row justify-content-center align-items-center">
                {["Project 1", "Project 2", "Project 3", "Project 4"].map((project, index) => (
                    <a href="#" key={index} className="card card-as-button" style={{ width: '280px' }}>
                        <Image
                            src="https://picsum.photos/280/180"
                            loader={imageLoader}
                            alt="Card image"
                            width={280}
                            height={180}
                            />
                        <div className="card-body card-body-project d-flex flex-column">
                            <h5 className="card-title"><strong>{project}</strong></h5>
                            <p className="card-text">Description for {project}</p>
                        </div>
                    </a>
                ))}
            </div>
        </div>
    );
}