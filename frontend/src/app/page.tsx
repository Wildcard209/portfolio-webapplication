import HeroBanner from "@/app/components/HeroBanner/HeroBanner";
import ProjectCard from "@/app/components/ProjectCard";
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
