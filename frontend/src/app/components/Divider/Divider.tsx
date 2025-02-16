import styles from "./Divider.module.css";

export default function Divider() {
    return (
        <div>
            <hr className={`${styles["divider"]}`} />
        </div>
    )}