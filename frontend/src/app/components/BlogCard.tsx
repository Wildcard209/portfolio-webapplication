'use client';

import React, { useEffect, useState, useRef } from "react";
import Image from "next/image";

type BlogCardProps = {
    name: string;
    description: string;
    lastUpdated: Date;
    imageId?: string;
    imageAlt?: string;
};

const imageLoader = ({ width }: { width: number }) => {
    return `https://picsum.photos/${width}/180`;
};

const BlogCard: React.FC<BlogCardProps> = ({ name, description, lastUpdated, imageId, imageAlt = "Blog image" }) => {
    const [relativeTime, setRelativeTime] = useState<string>("");
    const [isLargeScreen, setIsLargeScreen] = useState<boolean>(true);
    const [isVisible, setIsVisible] = useState<boolean>(false); // Track card visibility in viewport
    const cardRef = useRef<HTMLDivElement | null>(null); // Ref for the card

    const calculateRelativeTime = (date: Date) => {
        const now = new Date();
        const differenceInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

        if (differenceInSeconds < 60) return `${differenceInSeconds} seconds ago`;
        const differenceInMinutes = Math.floor(differenceInSeconds / 60);
        if (differenceInMinutes < 60) return `${differenceInMinutes} minutes ago`;
        const differenceInHours = Math.floor(differenceInMinutes / 60);
        if (differenceInHours < 24) return `${differenceInHours} hours ago`;
        const differenceInDays = Math.floor(differenceInHours / 24);
        if (differenceInDays < 30) return `${differenceInDays} days ago`;
        const differenceInMonths = Math.floor(differenceInDays / 30);
        if (differenceInMonths < 12) return `${differenceInMonths} months ago`;
        const differenceInYears = Math.floor(differenceInMonths / 12);
        return `${differenceInYears} years ago`;
    };

    useEffect(() => {
        const handleResize = () => {
            setIsLargeScreen(window.innerWidth >= 768);
        };
        handleResize();
        window.addEventListener("resize", handleResize);
        return () => window.removeEventListener("resize", handleResize);
    }, []);

    useEffect(() => {
        const updateRelativeTime = () => {
            setRelativeTime(calculateRelativeTime(lastUpdated));
        };

        if (isVisible) {
            updateRelativeTime();
            const intervalId = setInterval(updateRelativeTime, 1000);
            return () => clearInterval(intervalId);
        }
    }, [isVisible, lastUpdated]);

    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                setIsVisible(entry.isIntersecting);
            },
            { threshold: 0.1 }
        );

        if (cardRef.current) {
            observer.observe(cardRef.current);
        }

        return () => {
            if (cardRef.current) {
                observer.unobserve(cardRef.current);
            }
        };
    }, []);

    return (
        <div ref={cardRef} className="card card-as-button card-blog">
            <div className="row g-0 align-items-center">
                <div className="col-md-8">
                    <div className="card-body">
                        <h5 className="card-title">
                            <strong>{name}</strong>
                        </h5>
                        <p className="card-text">{description}</p>
                        <p className="card-text">
                            <small className="text-body-secondary">Last updated {relativeTime}</small>
                        </p>
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
                            style={{
                                objectFit: "cover"
                            }}
                        />
                    </div>
                )}
            </div>
        </div>
    );
};

export default BlogCard;