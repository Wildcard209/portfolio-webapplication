'use client';

import HeroBanner from "@/app/components/HeroBanner/HeroBanner";
import Divider from "@/app/components/Divider/Divider";
import Blog from "@/app/blog/page";
import Projects from "@/app/projects/page";
import {useEffect, useState} from "react";

export default function Home() {
    const [data, setData] = useState<string>();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/api/test');
                const result = await response.json();
                setData(result);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    console.log(data);

    return (
      <div>
          <HeroBanner />
          <Projects />
          <Divider />
          <Blog />
      </div>
  );
}
