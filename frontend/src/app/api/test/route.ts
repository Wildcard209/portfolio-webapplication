process.env.API_DOMAIN

import { NextResponse } from 'next/server';

// Define allowed origin
const ALLOWED_ORIGIN = process.env.API_DOMAIN ?? "";

export async function GET(request: Request) {
    const origin = request.headers.get('origin');

    if (origin !== ALLOWED_ORIGIN) {
        return new NextResponse(JSON.stringify({ error: 'Not allowed by CORS' }), {
            status: 403,
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': ALLOWED_ORIGIN,
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
