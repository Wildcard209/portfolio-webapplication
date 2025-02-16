import Image from 'next/image';

export default function ImageContent({ images, alignment }: { images: any[], alignment: string }) {
    return (
        <div className={`image-center ${alignment === 'left' ? 'image-left' : 'image-right'}`}>
            {images.map((img, index) => (
                <Image
                    key={index}
                    src={img.src}
                    className="image-shrink"
                    width={500}
                    height={500}
                    alt={img.alt}
                />
            ))}
        </div>
    );
}