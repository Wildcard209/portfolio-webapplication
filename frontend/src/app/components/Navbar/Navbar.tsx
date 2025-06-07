"use client";

import styles from "./Navbar.module.scss";
import Link from 'next/link';
import Image from 'next/image';
import { usePathname } from 'next/navigation';
import { useState } from 'react';

export default function Navbar() {
    const pathname = usePathname();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const toggleMenu = () => {
        setIsMenuOpen(!isMenuOpen);
    };

    return (
        <nav className={styles.navbar}>
            <div className={styles.container}>
                <Link href="/" className={`${styles.brand} fancy-font`}>
                    <Image
                        src="/assets/logos/LogoJWWhite.svg"
                        width={50}
                        height={50}
                        alt="JW"
                    />
                    <span>Jessica Wylde</span>
                </Link>
                <button
                    className={styles["nav-button"]}
                    type="button"
                    onClick={toggleMenu}
                    aria-expanded={isMenuOpen}
                    aria-label="Toggle navigation"
                >
                    <div className={`${styles["menu-icon"]} ${isMenuOpen ? styles.open : ''}`}>
                        <span></span>
                        <span></span>
                        <span></span>
                    </div>
                </button>
                <div className={`${styles.collapse} ${isMenuOpen ? styles.show : ''}`}>
                    <ul className={styles.nav}>
                        <li className={styles["nav-item"]}>
                            <Link href="/" className={`${styles["nav-link"]} ${pathname === '/' ? styles.active : ''}`}>
                                Home
                            </Link>
                        </li>
                        <li className={styles["nav-item"]}>
                            <div className={styles["nav-divider"]}></div>
                        </li>
                        <li className={styles["nav-item"]}>
                            <Link href="/about-me" className={`${styles["nav-link"]} ${pathname === '/about-me' ? styles.active : ''}`}>
                                About Me
                            </Link>
                        </li>
                        <li className={styles["nav-item"]}>
                            <div className={styles["nav-divider"]}></div>
                        </li>
                        <li className={styles["nav-item"]}>
                            <Link href="/projects" className={`${styles["nav-link"]} ${pathname === '/projects' ? styles.active : ''}`}>
                                Projects
                            </Link>
                        </li>
                        <li className={styles["nav-item"]}>
                            <div className={styles["nav-divider"]}></div>
                        </li>
                        <li className={styles["nav-item"]}>
                            <Link href="/blog" className={`${styles["nav-link"]} ${pathname === '/blog' ? styles.active : ''}`}>
                                Blog
                            </Link>
                        </li>
                        <li className={styles["nav-item"]}>
                            <div className={styles["nav-divider"]}></div>
                        </li>
                        <li className={styles["nav-item"]}>
                            <Link href="/contact" className={`${styles["nav-link"]} ${pathname === '/contact' ? styles.active : ''}`}>
                                Contact
                            </Link>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    );
}
