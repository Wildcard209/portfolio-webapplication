// pages/about-me.js
import ImageContent from "../components/CodeTable/ImageContent";
import TextContent from "../components/CodeTable/TextContent";

export default function AboutMePage() {
    return (
        <div>
            <h1 className="fancy-font text-title">About Me</h1>
            <table className="table equal-width-table">
                <tbody>
                <tr>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/csharp.svg', alt: 'C#' }]}
                            alignment="right"
                        />
                    </td>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <TextContent
                        text="Lorem ipsum dolor sit amet vero ipsum vero..."
                        alignment="right"
                        />

                    </td>
                    <td>
                        <ImageContent
                            images={[
                                { src: '/assets/logos/mysql.svg', alt: 'MySql' },
                                { src: '/assets/logos/sql-server.svg', alt: 'SQL'}]}
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/git.svg', alt: 'Git' }]}
                            alignment="right"
                        />
                    </td>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="right"
                        />

                    </td>
                    <td>
                        <ImageContent
                            images={[
                                { src: '/assets/logos/html.svg', alt: 'HTML' },
                                { src: '/assets/logos/css.svg', alt: 'CSS'}]}
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/javascript.svg', alt: 'Javascript' },
                                { src: '/assets/logos/typescript.svg', alt: 'Typescript'}
                            ]}
                            alignment="right"
                        />
                    </td>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="right"
                        />
                    </td>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/c++.svg', alt: 'C++' }]}
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/python.svg', alt: 'Python' }]}
                            alignment="right"
                        />
                    </td>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="right"
                        />
                    </td>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/opengl.svg', alt: 'OpenGL' }]}
                            alignment="left"
                        />
                    </td>
                </tr>
                <tr>
                    <td>
                        <ImageContent
                            images={[{ src: '/assets/logos/rust.svg', alt: 'Rust' }]}
                            alignment="right"
                        />
                    </td>
                    <td>
                        <TextContent
                            text="Lorem ipsum dolor sit amet vero ipsum vero..."
                            alignment="left"
                        />
                    </td>
                </tr>
                </tbody>
            </table>
        </div>
    );
}