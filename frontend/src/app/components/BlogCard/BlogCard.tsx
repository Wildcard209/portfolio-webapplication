'use client';

import styles from "./BlogCard.module.css"
import React, { useEffect, useState } from "react";
import Image from "next/image";
import { useRelativeTime } from "../../hooks/useRelativeTime";
import { useVisibilityObserver } from "../../hooks/useVisibilityObserver";

type BlogCardProps = {
    title: string;
    description: string;
    lastUpdated: Date;
    imageId?: string;
    imageAlt?: string;
};

const imageLoader = ({ width }: { width: number }) => {
    return `https://picsum.photos/${width}/180`;
};

const BlogCard: React.FC<BlogCardProps> = ({ title, description, lastUpdated, imageId, imageAlt = "Blog image" }) => {
    const [isLargeScreen, setIsLargeScreen] = useState<boolean>(true);
    const relativeTime = useRelativeTime(lastUpdated);
    const [isVisible, cardRef] = useVisibilityObserver<HTMLDivElement>();

    useEffect(() => {
        const handleResize = () => {
            setIsLargeScreen(window.innerWidth >= 768);
        };
        handleResize();
        window.addEventListener("resize", handleResize);
        return () => window.removeEventListener("resize", handleResize);
    }, []);

    return (
        <div ref={cardRef} className={`card card-as-button ${styles["card-blog"]}`}>
            <div className="row g-0 align-items-center">
                <div className="col-md-8">
                    <div className="card-body">
                        <h5 className="card-title">
                            <strong>{title}</strong>
                        </h5>
                        <p className="card-text">{description}</p>
                        {isVisible && (
                            <p className="card-text">
                                <small className="text-body-secondary">Last updated {relativeTime}</small>
                            </p>
                        )}
                    </div>
                </div>

                {imageId && isLargeScreen && (
                    <div className="col-md-4 d-flex justify-content-end">
                        <Image
                            src={imageId}
                            loader={imageLoader}
                            alt={imageAlt}
                            width={180}
                            height={180}
                            className={`${styles["card-image-blog"]}`}
                        />
                    </div>
                )}
            </div>
        </div>
    );
};

export default BlogCard;