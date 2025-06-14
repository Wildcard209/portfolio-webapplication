"use client";

import { useState, useEffect, useRef } from "react";
import { useAuth } from "@/lib/hooks/useAuth";
import { EnhancedApiHandler } from "@/lib/api/enhancedApiHandler";
import styles from "./HeroBanner.module.scss";

export default function HeroBanner() {
  const { isAuthenticated } = useAuth();
  const [backgroundImage, setBackgroundImage] = useState<string>("");
  const [showUploadControls, setShowUploadControls] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string>("");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    loadHeroBanner();
  }, []);

  const loadHeroBanner = async () => {
    try {
      const infoResponse = await EnhancedApiHandler.get<any>('/assets/info');
      
      if (infoResponse.data?.hero_banner_available) {
        const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
        setBackgroundImage(`${apiUrl}/assets/hero-banner?t=${Date.now()}`);
      }
    } catch (error) {
      console.log('No hero banner found or error loading:', error);
    }
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      if (!file.type.startsWith('image/')) {
        alert('Please select an image file');
        return;
      }

      if (file.size > 10 * 1024 * 1024) {
        alert('File size must be less than 10MB');
        return;
      }

      setSelectedFile(file);
      
      const preview = URL.createObjectURL(file);
      setPreviewUrl(preview);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    setUploading(true);
    try {
      const adminToken = localStorage.getItem('admin_token');
      if (!adminToken) {
        alert('Admin token not found');
        return;
      }

      const formData = new FormData();
      formData.append('file', selectedFile);

      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const token = localStorage.getItem('auth_token');
      
      const response = await fetch(`${apiUrl}/${adminToken}/admin/assets/hero-banner`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.ok) {
        await loadHeroBanner();
        setShowUploadControls(false);
        setSelectedFile(null);
        setPreviewUrl("");
        alert('Hero banner updated successfully!');
      } else {
        const error = await response.json();
        alert(`Upload failed: ${error.message || error.error}`);
      }
    } catch (error) {
      console.error('Upload error:', error);
      alert('Upload failed. Please try again.');
    } finally {
      setUploading(false);
    }
  };

  const handleCancel = () => {
    setShowUploadControls(false);
    setSelectedFile(null);
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
      setPreviewUrl("");
    }
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const currentBackgroundImage = previewUrl || backgroundImage;

  return (
    <div className={styles["hero-banner"]}>
      <div 
        className={`${styles["background-layer"]}`}
        style={{
          backgroundImage: currentBackgroundImage ? `url(${currentBackgroundImage})` : undefined,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat'
        }}
      ></div>
      
      <div className={`${styles["banner-text"]}`}>
        <h1 className={`fancy-font ${styles["banner-text-large"]}`}>
          Jessica Wylde
        </h1>
        <p>Software Engineer</p>
      </div>

      {/* Admin controls */}
      {isAuthenticated && !showUploadControls && (
        <button
          onClick={() => setShowUploadControls(true)}
          style={{
            position: 'absolute',
            top: '20px',
            right: '20px',
            padding: '10px 15px',
            backgroundColor: 'rgba(0, 123, 255, 0.8)',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '14px',
            zIndex: 10
          }}
        >
          Change Hero Banner
        </button>
      )}

      {/* Upload controls */}
      {isAuthenticated && showUploadControls && (
        <div style={{
          position: 'absolute',
          top: '20px',
          right: '20px',
          backgroundColor: 'rgba(255, 255, 255, 0.95)',
          padding: '15px',
          borderRadius: '8px',
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
          zIndex: 10,
          minWidth: '250px'
        }}>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            style={{ marginBottom: '10px', width: '100%' }}
          />
          
          {selectedFile && (
            <div style={{ marginBottom: '10px' }}>
              <p style={{ margin: '0 0 5px 0', fontSize: '12px', color: '#666' }}>
                Selected: {selectedFile.name}
              </p>
            </div>
          )}

          <div style={{ display: 'flex', gap: '10px' }}>
            <button
              onClick={handleUpload}
              disabled={!selectedFile || uploading}
              style={{
                flex: 1,
                padding: '8px 12px',
                backgroundColor: selectedFile && !uploading ? '#28a745' : '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: selectedFile && !uploading ? 'pointer' : 'not-allowed',
                fontSize: '12px'
              }}
            >
              {uploading ? 'Uploading...' : 'Upload'}
            </button>
            
            <button
              onClick={handleCancel}
              disabled={uploading}
              style={{
                flex: 1,
                padding: '8px 12px',
                backgroundColor: '#dc3545',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: uploading ? 'not-allowed' : 'pointer',
                fontSize: '12px'
              }}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
