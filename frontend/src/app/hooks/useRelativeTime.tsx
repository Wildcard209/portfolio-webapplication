import { useState, useEffect } from "react";

export function useRelativeTime(lastUpdated: Date): string {
    const [relativeTime, setRelativeTime] = useState<string>("");

    const calculateRelativeTime = (date: Date) => {
        const now = new Date();
        const differenceInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

        // Handle future dates
        if (differenceInSeconds < 0) return "in the future";

        const formatTimeUnit = (value: number, unit: string): string =>
            `${value} ${unit}${value !== 1 ? "s" : ""} ago`;

        if (differenceInSeconds < 60) return formatTimeUnit(differenceInSeconds, "second");
        const differenceInMinutes = Math.floor(differenceInSeconds / 60);
        if (differenceInMinutes < 60) return formatTimeUnit(differenceInMinutes, "minute");
        const differenceInHours = Math.floor(differenceInMinutes / 60);
        if (differenceInHours < 24) return formatTimeUnit(differenceInHours, "hour");
        const differenceInDays = Math.floor(differenceInHours / 24);
        if (differenceInDays < 30) return formatTimeUnit(differenceInDays, "day");
        const differenceInMonths = Math.floor(differenceInDays / 30);
        if (differenceInMonths < 12) return formatTimeUnit(differenceInMonths, "month");
        const differenceInYears = Math.floor(differenceInMonths / 12);
        return formatTimeUnit(differenceInYears, "year");
    };

    useEffect(() => {
        const updateRelativeTime = () => {
            setRelativeTime(calculateRelativeTime(lastUpdated));
        };

        updateRelativeTime();
        const intervalId = setInterval(updateRelativeTime, 1000);
        return () => clearInterval(intervalId);
    }, [lastUpdated]);

    return relativeTime;
}