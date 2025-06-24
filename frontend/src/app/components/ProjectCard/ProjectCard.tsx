'use client';

import Image from 'next/image';
import { useRelativeTime } from '../../hooks/useRelativeTime';
import { useVisibilityObserver } from '../../hooks/useVisibilityObserver';
import styles from './Project.module.scss';
import React from 'react';

type ProjectCardProps = {
  title: string;
  description: string;
  lastUpdated: Date;
  imageId: string;
  imageAlt?: string;
};

const imageLoader = ({ width }: { width: number }) => {
  return `https://picsum.photos/${width}/180`;
};

const ProjectCard: React.FC<ProjectCardProps> = ({
  title,
  description,
  lastUpdated,
  imageId,
  imageAlt = 'Project Image',
}) => {
  const relativeTime = useRelativeTime(lastUpdated);
  const [isVisible, cardRef] = useVisibilityObserver<HTMLDivElement>();

  return (
    <a href="#" className={styles['card-project']}>
      <Image
        src={imageId}
        loader={imageLoader}
        alt={imageAlt}
        width={280}
        height={180}
        className={styles['card-image-project']}
      />
      <div ref={cardRef} className={styles['card-body-project']}>
        <h5 className={styles['card-title']}>{title}</h5>
        <p className={styles['card-text']}>{description}</p>
        {isVisible && (
          <p className={styles['text-secondary']}>
            <small>Uploaded {relativeTime}</small>
          </p>
        )}
      </div>
    </a>
  );
};

export default ProjectCard;
