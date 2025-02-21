'use client';

import { useState } from 'react';

const TwoFactorAuthentication = () => {
    const [code, setCode] = useState('');
    const [is2FAVerified, setIs2FAVerified] = useState(false);

    const handle2FA = async (e: React.FormEvent) => {
        e.preventDefault();

        if (code === '111111') {
            setIs2FAVerified(true);
        } else {
            alert('Invalid 2FA code. Try again!');
        }
    };

    if (is2FAVerified) {
        return (
            <div>
                <h2>Access Granted</h2>
                <p>You are now authorized to access the Admin Panel!</p>
            </div>
        );
    }

    return (
        <form onSubmit={handle2FA} style={{ textAlign: 'center' }}>
            <h2>Two-Factor Authentication</h2>
            <div style={{ marginBottom: '10px' }}>
                <label>Enter 2FA Code:</label>
                <input
                    type="text"
                    value={code}
                    maxLength={6}
                    onChange={(e) => setCode(e.target.value)}
                    required
                />
            </div>
            <button type="submit">Verify</button>
        </form>
    );
};

export default TwoFactorAuthentication;