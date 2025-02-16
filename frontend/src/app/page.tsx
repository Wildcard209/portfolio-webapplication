import HeroBanner from "@/app/components/HeroBanner/HeroBanner";
import ProjectCard from "@/app/components/ProjectCard";
import Divider from "@/app/components/Divider/Divider";
import Blog from "@/app/blog/page";

export default function Home() {
  return (
      <div>
          <HeroBanner />
          <ProjectCard />
          <Divider />
          <Blog />
      </div>
  );
}
