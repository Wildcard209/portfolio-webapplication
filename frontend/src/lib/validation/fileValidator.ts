export interface FileValidationResult {
  isValid: boolean;
  error?: string;
}

export interface FileValidationOptions {
  maxSize?: number;
  allowedTypes?: string[];
  maxNameLength?: number;
}

export class FileValidator {
  private static readonly DEFAULT_MAX_SIZE = 10 * 1024 * 1024; // 10MB
  private static readonly DEFAULT_ALLOWED_TYPES = [
    'image/jpeg',
    'image/jpg', 
    'image/png',
    'image/gif',
    'image/webp'
  ];
  private static readonly DEFAULT_MAX_NAME_LENGTH = 255;

  static validateFile(
    file: File, 
    options: FileValidationOptions = {}
  ): FileValidationResult {
    const {
      maxSize = this.DEFAULT_MAX_SIZE,
      allowedTypes = this.DEFAULT_ALLOWED_TYPES,
      maxNameLength = this.DEFAULT_MAX_NAME_LENGTH
    } = options;

    if (!file) {
      return {
        isValid: false,
        error: 'No file selected'
      };
    }

    if (file.size === 0) {
      return {
        isValid: false,
        error: 'File is empty'
      };
    }

    if (file.size > maxSize) {
      return {
        isValid: false,
        error: `File size (${this.formatFileSize(file.size)}) exceeds maximum allowed size (${this.formatFileSize(maxSize)})`
      };
    }

    if (!allowedTypes.includes(file.type)) {
      return {
        isValid: false,
        error: `File type "${file.type}" is not allowed. Allowed types: ${allowedTypes.join(', ')}`
      };
    }

    const filenameValidation = this.validateFilename(file.name, maxNameLength);
    if (!filenameValidation.isValid) {
      return filenameValidation;
    }

    return {
      isValid: true
    };
  }

  static validateFilename(filename: string, maxLength: number): FileValidationResult {
    if (!filename || filename.length === 0) {
      return {
        isValid: false,
        error: 'Filename is empty'
      };
    }

    if (filename.length > maxLength) {
      return {
        isValid: false,
        error: `Filename too long (max ${maxLength} characters)`
      };
    }

    const dangerousChars = ['..', '/', '\\', ':', '*', '?', '"', '<', '>', '|', '\0'];
    for (const char of dangerousChars) {
      if (filename.includes(char)) {
        return {
          isValid: false,
          error: `Filename contains dangerous character: "${char}"`
        };
      }
    }

    if (!filename.includes('.')) {
      return {
        isValid: false,
        error: 'Filename must have an extension'
      };
    }

    const validFilenamePattern = /^[a-zA-Z0-9._-]+$/;
    if (!validFilenamePattern.test(filename)) {
      return {
        isValid: false,
        error: 'Filename contains invalid characters (only letters, numbers, dots, underscores, and hyphens allowed)'
      };
    }

    return {
      isValid: true
    };
  }

  static validateFiles(
    files: FileList | File[], 
    options: FileValidationOptions = {}
  ): FileValidationResult {
    const fileArray = Array.from(files);
    
    if (fileArray.length === 0) {
      return {
        isValid: false,
        error: 'No files selected'
      };
    }

    for (let i = 0; i < fileArray.length; i++) {
      const result = this.validateFile(fileArray[i], options);
      if (!result.isValid) {
        return {
          isValid: false,
          error: `File ${i + 1}: ${result.error}`
        };
      }
    }

    return {
      isValid: true
    };
  }

  static formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  static isImageFile(file: File): boolean {
    return file.type.startsWith('image/');
  }

  static getFileExtension(filename: string): string {
    return filename.split('.').pop()?.toLowerCase() || '';
  }

  static sanitizeFilename(filename: string): string {
    const sanitized = filename.replace(/[^a-zA-Z0-9._-]/g, '_');
    
    const cleaned = sanitized.replace(/_+/g, '_');
    
    return cleaned.replace(/^_+|_+$/g, '');
  }
}
