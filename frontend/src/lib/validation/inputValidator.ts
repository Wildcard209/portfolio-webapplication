// Input validation utilities for client-side validation
export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

export interface ValidationOptions {
  minLength?: number;
  maxLength?: number;
  required?: boolean;
  pattern?: RegExp;
  customValidator?: (value: string) => ValidationResult;
}

export class InputValidator {
  static validateString(
    value: string,
    fieldName: string,
    options: ValidationOptions = {}
  ): ValidationResult {
    const { minLength = 0, maxLength = 1000, required = false, pattern, customValidator } = options;

    if (required && (!value || value.trim().length === 0)) {
      return {
        isValid: false,
        error: `${fieldName} is required`,
      };
    }

    if (!required && (!value || value.trim().length === 0)) {
      return {
        isValid: true,
      };
    }

    if (value.length < minLength) {
      return {
        isValid: false,
        error: `${fieldName} must be at least ${minLength} characters long`,
      };
    }

    if (value.length > maxLength) {
      return {
        isValid: false,
        error: `${fieldName} must not exceed ${maxLength} characters`,
      };
    }

    if (pattern && !pattern.test(value)) {
      return {
        isValid: false,
        error: `${fieldName} format is invalid`,
      };
    }

    const dangerousPatterns = this.getDangerousPatterns();
    for (const patternInfo of dangerousPatterns) {
      if (patternInfo.pattern.test(value.toLowerCase())) {
        return {
          isValid: false,
          error: `${fieldName} contains potentially dangerous content`,
        };
      }
    }

    if (customValidator) {
      const customResult = customValidator(value);
      if (!customResult.isValid) {
        return customResult;
      }
    }

    return {
      isValid: true,
    };
  }

  static validateUsername(username: string): ValidationResult {
    if (!username || username.trim().length === 0) {
      return {
        isValid: false,
        error: 'Username is required',
      };
    }

    if (username.length < 3) {
      return {
        isValid: false,
        error: 'Username must be at least 3 characters long',
      };
    }

    if (username.length > 50) {
      return {
        isValid: false,
        error: 'Username must not exceed 50 characters',
      };
    }

    const usernamePattern = /^[a-zA-Z0-9_-]+$/;
    if (!usernamePattern.test(username)) {
      return {
        isValid: false,
        error: 'Username can only contain letters, numbers, underscores, and hyphens',
      };
    }

    if (!/^[a-zA-Z0-9]/.test(username)) {
      return {
        isValid: false,
        error: 'Username must start with a letter or number',
      };
    }

    return {
      isValid: true,
    };
  }

  static validatePassword(password: string): ValidationResult {
    if (!password) {
      return {
        isValid: false,
        error: 'Password is required',
      };
    }

    if (password.length < 8) {
      return {
        isValid: false,
        error: 'Password must be at least 8 characters long',
      };
    }

    if (password.length > 128) {
      return {
        isValid: false,
        error: 'Password must not exceed 128 characters',
      };
    }

    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasDigit = /[0-9]/.test(password);
    const hasSpecial = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password);

    const missing = [];
    if (!hasUpper) missing.push('uppercase letter');
    if (!hasLower) missing.push('lowercase letter');
    if (!hasDigit) missing.push('digit');
    if (!hasSpecial) missing.push('special character');

    if (missing.length > 0) {
      return {
        isValid: false,
        error: `Password must contain at least one ${missing.join(', ')}`,
      };
    }

    return {
      isValid: true,
    };
  }

  static validateEmail(email: string): ValidationResult {
    if (!email || email.trim().length === 0) {
      return {
        isValid: false,
        error: 'Email is required',
      };
    }

    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailPattern.test(email)) {
      return {
        isValid: false,
        error: 'Please enter a valid email address',
      };
    }

    if (email.length > 254) {
      return {
        isValid: false,
        error: 'Email address is too long',
      };
    }

    return {
      isValid: true,
    };
  }

  static sanitizeInput(input: string): string {
    if (!input) return '';

    let sanitized = input.replace(/[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]/g, '');

    sanitized = sanitized
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#x27;');

    return sanitized.trim();
  }

  static validateFields(
    fields: Array<{
      value: string;
      name: string;
      options?: ValidationOptions;
    }>
  ): ValidationResult {
    for (const field of fields) {
      const result = this.validateString(field.value, field.name, field.options);
      if (!result.isValid) {
        return result;
      }
    }

    return {
      isValid: true,
    };
  }

  static getPasswordStrength(password: string): number {
    if (!password) return 0;

    let score = 0;

    if (password.length >= 8) score++;
    if (password.length >= 12) score++;

    if (/[a-z]/.test(password)) score++;
    if (/[A-Z]/.test(password)) score++;
    if (/[0-9]/.test(password)) score++;
    if (/[^a-zA-Z0-9]/.test(password)) score++;

    return Math.min(score, 4);
  }

  private static getDangerousPatterns(): Array<{ pattern: RegExp; description: string }> {
    return [
      { pattern: /union\s+select/i, description: 'SQL injection' },
      { pattern: /drop\s+table/i, description: 'SQL injection' },
      { pattern: /delete\s+from/i, description: 'SQL injection' },
      { pattern: /insert\s+into/i, description: 'SQL injection' },
      { pattern: /update\s+set/i, description: 'SQL injection' },
      { pattern: /<script/i, description: 'XSS' },
      { pattern: /javascript:/i, description: 'XSS' },
      { pattern: /vbscript:/i, description: 'XSS' },
      { pattern: /onload\s*=/i, description: 'XSS' },
      { pattern: /onerror\s*=/i, description: 'XSS' },
      { pattern: /onclick\s*=/i, description: 'XSS' },
      { pattern: /eval\s*\(/i, description: 'Code injection' },
      { pattern: /expression\s*\(/i, description: 'CSS injection' },
    ];
  }
}
