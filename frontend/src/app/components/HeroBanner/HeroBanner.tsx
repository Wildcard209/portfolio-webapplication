'use client';

import { useState, useRef } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import { useApi, useAdminApiFileUpload, getApiAssetUrl } from '@/lib/api/hooks/useApi';
import { FileValidator } from '@/lib/validation/fileValidator';
import styles from './HeroBanner.module.scss';

export default function HeroBanner() {
  const { isAuthenticated } = useAuth();

  const { data: assetInfo, refetch: refetchAssetInfo } = useApi<{ hero_banner_available: boolean }>(
    '/assets/info'
  );

  const {
    uploadFile,
    isLoading: uploading,
    error: uploadError,
  } = useAdminApiFileUpload('/assets/hero-banner');

  const [showUploadControls, setShowUploadControls] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string>('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const backgroundImage = assetInfo?.hero_banner_available
    ? getApiAssetUrl('/assets/hero-banner')
    : '';

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const validationResult = FileValidator.validateFile(file, {
      maxSize: 10 * 1024 * 1024,
      allowedTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
      maxNameLength: 255,
    });

    if (!validationResult.isValid) {
      alert(`File validation failed: ${validationResult.error}`);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      return;
    }

    if (!FileValidator.isImageFile(file)) {
      alert('Please select an image file');
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      return;
    }

    setSelectedFile(file);

    const preview = URL.createObjectURL(file);
    setPreviewUrl(preview);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      alert('No file selected');
      return;
    }

    try {
      const formData = new FormData();
      formData.append('file', selectedFile);

      const result = await uploadFile(formData);

      if (result) {
        await refetchAssetInfo();
        setShowUploadControls(false);
        setSelectedFile(null);
        setPreviewUrl('');
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
      setPreviewUrl('');
    }
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const currentBackgroundImage = previewUrl || backgroundImage;

  return (
    <div className={styles['hero-banner']}>
      <div
        className={`${styles['background-layer']}`}
        style={{
          backgroundImage: currentBackgroundImage ? `url(${currentBackgroundImage})` : undefined,
        }}
      ></div>

      <div className={`${styles['banner-text']}`}>
        <h1 className={`fancy-font ${styles['banner-text-large']}`}>Jessica Wylde</h1>
        <p>Software Engineer</p>
      </div>

      {/* Admin controls */}
      {isAuthenticated && !showUploadControls && (
        <button
          onClick={() => setShowUploadControls(true)}
          className={styles['admin-toggle-button']}
        >
          Change Hero Banner
        </button>
      )}

      {isAuthenticated && showUploadControls && (
        <div className={styles['admin-controls-panel']}>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            className={styles['file-input']}
          />

          {selectedFile && (
            <div className={styles['file-info']}>
              <p>Selected: {selectedFile.name}</p>
            </div>
          )}

          <div className={styles['button-group']}>
            <button
              onClick={handleUpload}
              disabled={!selectedFile || uploading}
              className={`${styles['upload-button']} ${
                selectedFile && !uploading ? styles['enabled'] : styles['disabled']
              }`}
            >
              {uploading ? 'Uploading...' : 'Upload'}
            </button>

            <button
              onClick={handleCancel}
              disabled={uploading}
              className={`${styles['cancel-button']} ${uploading ? styles['disabled'] : ''}`}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
