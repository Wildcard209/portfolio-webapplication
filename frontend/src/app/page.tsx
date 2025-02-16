//import styles from "./page.module.css";
import HeroBanner from "@/app/components/HeroBanner";
import ProjectCard from "@/app/components/ProjectCard";
//import BlogCard from "@/app/components/BlogCard";
import Blog from "@/app/blog/page";

export default function Home() {
  return (
      <div>
          <HeroBanner />
          <ProjectCard />
          <hr className="divider" />
          <Blog />
      </div>
  );
}
