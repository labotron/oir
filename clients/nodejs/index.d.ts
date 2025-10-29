/**
 * TypeScript definitions for ODT Image Replacer Client
 */

export interface ImageSource {
  url?: string | null;
  base64?: string | null;
  buffer?: Buffer;
  filePath?: string;
}

export interface TemplateSource {
  url?: string | null;
  base64?: string | null;
  buffer?: Buffer;
  filePath?: string;
}

export interface ReplaceImagesParams {
  template: TemplateSource;
  data: {
    [tagName: string]: ImageSource;
  };
}

export interface ReplaceImagesOptions {
  returnBase64?: boolean;
}

export interface ClientOptions {
  timeout?: number;
}

export interface HealthCheckResponse {
  status: string;
  service: string;
}

export interface APIInfo {
  service: string;
  version: string;
  description: string;
  endpoints: {
    [key: string]: string;
  };
}

export class ODTImageReplacerClient {
  constructor(baseURL?: string, options?: ClientOptions);

  /**
   * Replace images in an ODT document
   */
  replaceImages(
    params: ReplaceImagesParams,
    options?: ReplaceImagesOptions
  ): Promise<Buffer | string>;

  /**
   * Replace images and save to file
   */
  replaceImagesAndSave(
    params: ReplaceImagesParams,
    outputPath: string
  ): Promise<void>;

  /**
   * Download the modified ODT directly from the API
   */
  downloadModifiedODT(
    params: ReplaceImagesParams,
    outputPath: string
  ): Promise<void>;

  /**
   * Check API health
   */
  healthCheck(): Promise<HealthCheckResponse>;

  /**
   * Get API information
   */
  getInfo(): Promise<APIInfo>;
}

export interface QuickReplaceOptions {
  apiUrl?: string;
}

/**
 * Convenience function for quick usage
 */
export function replaceImages(
  templatePath: string,
  images: { [tagName: string]: string },
  outputPath: string,
  options?: QuickReplaceOptions
): Promise<void>;
