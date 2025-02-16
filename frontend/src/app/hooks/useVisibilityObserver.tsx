import { useState, useEffect, useRef } from "react";

export function useVisibilityObserver<T extends HTMLElement>(): [boolean, React.RefObject<T | null>] {
    const [isVisible, setIsVisible] = useState<boolean>(false);
    const ref = useRef<T | null>(null);

    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => setIsVisible(entry.isIntersecting),
            { threshold: 0.1 }
        );

        if (ref.current) {
            observer.observe(ref.current);
        }

        return () => {
            if (ref.current) {
                observer.unobserve(ref.current);
            }
        };
    }, []);

    return [isVisible, ref];
}