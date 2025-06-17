"use client";

import { useState, useRef } from "react";
import { useAuth } from "@/lib/hooks/useAuth";
import { useApi, useAdminApiFileUpload } from "@/lib/api/hooks/useApi";
import { ApiHandler } from "@/lib/api/apiHandler";
import styles from "./HeroBanner.module.scss";

export default function HeroBanner() {
  const { isAuthenticated } = useAuth();
  
  const { 
    data: assetInfo, 
    refetch: refetchAssetInfo 
  } = useApi<{ hero_banner_available: boolean }>('/assets/info');

  const { 
    uploadFile, 
    isLoading: uploading, 
    error: uploadError 
  } = useAdminApiFileUpload('/assets/hero-banner');

  const [showUploadControls, setShowUploadControls] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Get the background image URL using ApiHandler when hero banner is available
  const backgroundImage = assetInfo?.hero_banner_available 
    ? ApiHandler.getAssetUrl('/assets/hero-banner')
    : "";

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
    if (!selectedFile) {
      alert('No file selected');
      return;
    }

    try {
      // Use the admin file upload hook
      const formData = new FormData();
      formData.append('file', selectedFile);

      const result = await uploadFile(formData);

      if (result) {
        // Refetch the asset info to update the hero banner
        await refetchAssetInfo();
        setShowUploadControls(false);
        setSelectedFile(null);
        setPreviewUrl("");
        alert('Hero banner updated successfully!');
      } else if (uploadError) {
        alert(`Upload failed: ${uploadError}`);
      }
    } catch (error) {
      console.error('Upload error:', error);
      alert('Upload failed. Please try again.');
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
