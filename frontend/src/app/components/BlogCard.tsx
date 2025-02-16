'use client'

import React, { useEffect, useState } from "react";

type BlogCardProps = {
    name: string;
    description: string;
    lastUpdated: Date;
};

const BlogCard: React.FC<BlogCardProps> = ({ name, description, lastUpdated }) => {
    const [relativeTime, setRelativeTime] = useState<string>("");

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
        const updateRelativeTime = () => {
            setRelativeTime(calculateRelativeTime(lastUpdated));
        };

        updateRelativeTime();
        const intervalId = setInterval(updateRelativeTime, 1000);

        return () => clearInterval(intervalId);
    }, [lastUpdated]);

    return (
        <div className="card card-as-button card-blog">
            <div className="card-blog-body">
                <h5 className="card-title">
                    <strong>{name}</strong>
                </h5>
                <p className="card-text card-text-blog">{description}</p>
                <p className="card-text">
                    <small className="text-body-secondary">Last updated {relativeTime}</small>
                </p>
            </div>
        </div>
    );
};

export default BlogCard;