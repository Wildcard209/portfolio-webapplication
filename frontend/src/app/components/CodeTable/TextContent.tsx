// components/TextContent.js
export default function TextContent({ text, alignment } : { text: string, alignment: string }) {
    return (
        <div className={`text-shrink ${alignment === 'left' ? 'text-left' : 'text-right'}`}>
            <h5 className="text-shrink">{text}</h5>
        </div>
    );
}