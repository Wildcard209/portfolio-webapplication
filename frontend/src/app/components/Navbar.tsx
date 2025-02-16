"use client";

import Link from 'next/link';
import Image from 'next/image';
import { usePathname } from 'next/navigation';

export default function Navbar() {
    const pathname = usePathname();

    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
            <div className="container">
                <Link href="/" className="navbar-brand fancy-font">
                    <Image
                        src="/assets/logos/LogoJWWhite.svg"
                        width={50}
                        height={50}
                        alt="JW"
                    />
                    <span className="logo-text">Jessica Wylde</span>
                </Link>
                <button
                    className="navbar-toggler"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target="#navbarSupportedContent"
                    aria-controls="navbarSupportedContent"
                    aria-expanded="false"
                    aria-label="Toggle navigation"
                >
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarSupportedContent">
                    <ul className="navbar-nav ms-auto mb-2 mb-lg-0">
                        <li className="nav-item">
                            <Link href="/" className={`nav-link ${pathname === '/' ? 'active' : ''}`}>
                                Home
                            </Link>
                        </li>
                        <li className="nav-item">
                            <div className="nav-divider"></div>
                        </li>
                        <li className="nav-item">
                            <Link href="/about-me" className={`nav-link ${pathname === '/about-me' ? 'active' : ''}`}>
                                About Me
                            </Link>
                        </li>
                        <li className="nav-item">
                            <div className="nav-divider"></div>
                        </li>
                        <li className="nav-item">
                            <Link href="/projects" className={`nav-link ${pathname === '/projects' ? 'active' : ''}`}>
                                Projects
                            </Link>
                        </li>
                        <li className="nav-item">
                            <div className="nav-divider"></div>
                        </li>
                        <li className="nav-item">
                            <Link href="/blog" className={`nav-link ${pathname === '/blog' ? 'active' : ''}`}>
                                Blog
                            </Link>
                        </li>
                        <li className="nav-item">
                            <div className="nav-divider"></div>
                        </li>
                        <li className="nav-item">
                            <Link href="/contact" className={`nav-link ${pathname === '/contact' ? 'active' : ''}`}>
                                Contact
                            </Link>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    );
}
