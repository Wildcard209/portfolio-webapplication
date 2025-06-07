import { NextResponse } from 'next/server';

// Define allowed origin
const ALLOWED_ORIGIN = process.env.API_DOMAIN ?? "http://localhost:3000";

export async function GET(request: Request) {
    const origin = request.headers.get('origin') || "http://localhost:3000";
    
    // Allow requests from specified origin
    if (origin !== ALLOWED_ORIGIN && origin !== "http://localhost:3000") {
        return new NextResponse(JSON.stringify({ error: 'Not allowed by CORS' }), {
            status: 403,
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': origin,
            },
        });
    }

    try {
        const response = await fetch('http://localhost/api/hello');
        const data = await response.json();

        return new NextResponse(JSON.stringify(data), {
            status: 200,
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': ALLOWED_ORIGIN,
            },
        });
    } catch (error) {
        console.error('Error fetching external API:', error);

        return new NextResponse(JSON.stringify({ error: 'Failed to fetch external API' }), {
            status: 500,
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': ALLOWED_ORIGIN,
            },
        });
    }
}
