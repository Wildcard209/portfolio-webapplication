import styles from "./Footer.module.css"

export default function Footer() {
    return (
            <div className={`${styles["footer"]}`}>
                <p>© {new Date().getFullYear()} Jessica Wylde. All Rights Reserved.</p>
            </div>
    );
}
