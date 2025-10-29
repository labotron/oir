const axios = require('axios');
const fs = require('fs');

/**
 * ODT Image Replacer Client
 * A simple Node.js client for the ODT Image Replacer API
 */
class ODTImageReplacerClient {
  /**
   * Create a new client instance
   * @param {string} baseURL - The base URL of the API (e.g., 'http://localhost:8080')
   * @param {Object} options - Additional options
   * @param {number} options.timeout - Request timeout in milliseconds (default: 30000)
   */
  constructor(baseURL = 'http://localhost:8080', options = {}) {
    this.baseURL = baseURL;
    this.timeout = options.timeout || 30000;
    this.axios = axios.create({
      baseURL: this.baseURL,
      timeout: this.timeout,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  /**
   * Replace images in an ODT document
   * @param {Object} params - Parameters for image replacement
   * @param {Object} params.template - Template source
   * @param {string} params.template.url - URL to ODT template (optional)
   * @param {string} params.template.base64 - Base64-encoded ODT template (optional)
   * @param {Buffer} params.template.buffer - Buffer containing ODT template (optional)
   * @param {string} params.template.filePath - Local file path to ODT template (optional)
   * @param {Object} params.data - Image data mapping (tag name -> image source)
   * @param {Object} params.data[tagName].url - URL to image (optional)
   * @param {Object} params.data[tagName].base64 - Base64-encoded image (optional)
   * @param {Object} params.data[tagName].buffer - Buffer containing image data (optional)
   * @param {Object} params.data[tagName].filePath - Local file path to image (optional)
   * @param {Object} options - Additional options
   * @param {boolean} options.returnBase64 - Return base64 string instead of Buffer (default: false)
   * @returns {Promise<Buffer|string>} - The modified ODT file as Buffer or base64 string
   */
  async replaceImages(params, options = {}) {
    // Build the request payload
    const payload = await this._buildPayload(params);

    try {
      // Send request to API
      const response = await this.axios.post('/api/replace', payload);

      // Check if successful
      if (!response.data.success) {
        throw new Error(response.data.error || 'Image replacement failed');
      }

      // Return base64 or Buffer
      const base64Data = response.data.output_base64;
      if (options.returnBase64) {
        return base64Data;
      }

      // Convert base64 to Buffer
      return Buffer.from(base64Data, 'base64');
    } catch (error) {
      if (error.response && error.response.data) {
        throw new Error(error.response.data.error || error.message);
      }
      throw error;
    }
  }

  /**
   * Replace images and save to file
   * @param {Object} params - Parameters for image replacement (same as replaceImages)
   * @param {string} outputPath - Path where to save the output ODT file
   * @returns {Promise<void>}
   */
  async replaceImagesAndSave(params, outputPath) {
    const buffer = await this.replaceImages(params);
    await fs.promises.writeFile(outputPath, buffer);
  }

  /**
   * Download the modified ODT directly from the API
   * @param {Object} params - Parameters for image replacement
   * @param {string} outputPath - Path where to save the output ODT file
   * @returns {Promise<void>}
   */
  async downloadModifiedODT(params, outputPath) {
    const payload = await this._buildPayload(params);

    try {
      const response = await this.axios.post('/api/replace/download', payload, {
        responseType: 'arraybuffer'
      });

      await fs.promises.writeFile(outputPath, response.data);
    } catch (error) {
      if (error.response && error.response.data) {
        // Try to parse error message from arraybuffer
        const errorText = Buffer.from(error.response.data).toString();
        throw new Error(errorText || error.message);
      }
      throw error;
    }
  }

  /**
   * Check API health
   * @returns {Promise<Object>} - Health check response
   */
  async healthCheck() {
    const response = await this.axios.get('/health');
    return response.data;
  }

  /**
   * Get API information
   * @returns {Promise<Object>} - API information
   */
  async getInfo() {
    const response = await this.axios.get('/info');
    return response.data;
  }

  /**
   * Build the API payload from params
   * @private
   */
  async _buildPayload(params) {
    const payload = {
      template: {},
      data: {}
    };

    // Handle template
    if (params.template.buffer) {
      // Convert Buffer to base64
      payload.template.base64 = params.template.buffer.toString('base64');
      payload.template.url = null;
    } else if (params.template.filePath) {
      const fileData = await fs.promises.readFile(params.template.filePath);
      payload.template.base64 = fileData.toString('base64');
      payload.template.url = null;
    } else if (params.template.base64) {
      payload.template.base64 = params.template.base64;
      payload.template.url = null;
    } else if (params.template.url) {
      payload.template.url = params.template.url;
      payload.template.base64 = null;
    } else {
      throw new Error('Template source must be provided (url, base64, buffer, or filePath)');
    }

    // Handle images
    for (const [tag, source] of Object.entries(params.data)) {
      if (source.buffer) {
        // Convert Buffer to base64
        payload.data[tag] = {
          base64: source.buffer.toString('base64'),
          url: null
        };
      } else if (source.filePath) {
        const fileData = await fs.promises.readFile(source.filePath);
        payload.data[tag] = {
          base64: fileData.toString('base64'),
          url: null
        };
      } else if (source.base64) {
        payload.data[tag] = {
          base64: source.base64,
          url: null
        };
      } else if (source.url) {
        payload.data[tag] = {
          url: source.url,
          base64: null
        };
      } else {
        throw new Error(`Image source for tag '${tag}' must be provided (url, base64, buffer, or filePath)`);
      }
    }

    return payload;
  }
}

/**
 * Convenience function for quick usage
 * @param {string} templatePath - Path to ODT template file
 * @param {Object} images - Image mapping (tag name -> image file path)
 * @param {string} outputPath - Path where to save the output ODT file
 * @param {Object} options - Additional options
 * @param {string} options.apiUrl - API base URL (default: 'http://localhost:8080')
 * @returns {Promise<void>}
 */
async function replaceImages(templatePath, images, outputPath, options = {}) {
  const client = new ODTImageReplacerClient(options.apiUrl || 'http://localhost:8080');

  // Convert simple image paths to proper format
  const data = {};
  for (const [tag, imagePath] of Object.entries(images)) {
    data[tag] = { filePath: imagePath };
  }

  await client.replaceImagesAndSave({
    template: { filePath: templatePath },
    data: data
  }, outputPath);
}

module.exports = {
  ODTImageReplacerClient,
  replaceImages
};
