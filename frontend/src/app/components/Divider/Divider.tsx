import styles from "./Divider.module.scss";

export default function Divider() {
    return (
        <div>
            <hr className={`${styles["divider"]}`} />
        </div>
    )}