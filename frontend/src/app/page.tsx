import HeroBanner from "@/app/components/HeroBanner/HeroBanner";
import Divider from "@/app/components/Divider/Divider";
import Blog from "@/app/blog/page";
import Projects from "@/app/projects/page";

export default function Home() {
  return (
      <div>
          <HeroBanner />
          <Projects />
          <Divider />
          <Blog />
      </div>
  );
}
