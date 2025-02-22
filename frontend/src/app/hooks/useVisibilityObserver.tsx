import React, { useState, useEffect, useRef } from "react";

export function useVisibilityObserver<T extends HTMLElement>(): [boolean, React.RefObject<T | null>] {
    const [isVisible, setIsVisible] = useState<boolean>(false);
    const ref = useRef<T | null>(null);

    useEffect(() => {
        const currentRef = ref.current;
        const observer = new IntersectionObserver(
            ([entry]) => setIsVisible(entry.isIntersecting),
            { threshold: 0.1 }
        );

        if (currentRef) {
            observer.observe(currentRef);
        }

        return () => {
            if (currentRef) {
                observer.unobserve(currentRef);
            }
        };
    }, []);

    return [isVisible, ref];
}