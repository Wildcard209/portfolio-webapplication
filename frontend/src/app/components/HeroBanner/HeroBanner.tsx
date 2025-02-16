import styles from "./HeroBanner.module.css"

export default function HeroBanner() {
    return (
        <div className="col-lg-12">
            <div className={`${styles["hero-banner"]}`}>
                <div className={`${styles["background-layer"]}`}></div>
                <div className={`${styles["banner-text"]}`}>
                    <h1 className={`fancy-font ${styles["banner-text-large"]}`}>Jessica Wylde</h1>
                    <p>Software Engineer</p>
                </div>
            </div>
        </div>
    );
}