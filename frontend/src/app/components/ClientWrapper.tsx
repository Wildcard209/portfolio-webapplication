"use client";

if (typeof window !== "undefined") {
    // @ts-expect-error file dose exist at location
    import("bootstrap/dist/js/bootstrap.bundle.min");
}

const ClientWrapper: React.FC<{
    children: React.ReactNode;
}> = ({ children }) => {
    return <>{children}</>;
};

export default ClientWrapper;

