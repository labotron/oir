const { ODTImageReplacerClient, replaceImages } = require('./index');

// Example 1: Simple usage with the convenience function
async function example1() {
  console.log('Example 1: Simple usage with local files');
  console.log('==========================================\n');

  try {
    await replaceImages(
      './template.odt',           // Template path
      {
        image1: './photo1.png',   // Image for tag "image1"
        image2: './photo2.jpg'    // Image for tag "image2"
      },
      './output.odt'              // Output path
    );

    console.log('✓ Successfully replaced images and saved to output.odt\n');
  } catch (error) {
    console.error('✗ Error:', error.message, '\n');
  }
}

// Example 2: Using the client class with URLs
async function example2() {
  console.log('Example 2: Using URLs for template and images');
  console.log('==============================================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    const buffer = await client.replaceImages({
      template: {
        url: 'https://example.com/template.odt'
      },
      data: {
        logo: {
          url: 'https://example.com/logo.png'
        },
        signature: {
          url: 'https://example.com/signature.png'
        }
      }
    });

    console.log(`✓ Received ${buffer.length} bytes`);
    console.log('✓ You can now save it: fs.writeFileSync("output.odt", buffer)\n');
  } catch (error) {
    console.error('✗ Error:', error.message, '\n');
  }
}

// Example 3: Mixed sources (URL, base64, file paths)
async function example3() {
  console.log('Example 3: Mixed image sources');
  console.log('===============================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    await client.replaceImagesAndSave({
      template: {
        filePath: './template.odt'  // Local template
      },
      data: {
        logo: {
          url: 'https://example.com/logo.png'  // From URL
        },
        photo: {
          filePath: './photo.jpg'  // From local file
        },
        signature: {
          base64: 'iVBORw0KGgo...'  // From base64
        }
      }
    }, './output-mixed.odt');

    console.log('✓ Successfully saved to output-mixed.odt\n');
  } catch (error) {
    console.error('✗ Error:', error.message, '\n');
  }
}

// Example 4: Direct download endpoint
async function example4() {
  console.log('Example 4: Using download endpoint');
  console.log('===================================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    await client.downloadModifiedODT({
      template: {
        filePath: './template.odt'
      },
      data: {
        image1: {
          filePath: './photo.png'
        }
      }
    }, './output-download.odt');

    console.log('✓ Successfully downloaded to output-download.odt\n');
  } catch (error) {
    console.error('✗ Error:', error.message, '\n');
  }
}

// Example 5: Get base64 output instead of Buffer
async function example5() {
  console.log('Example 5: Get base64 output');
  console.log('=============================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    const base64Data = await client.replaceImages({
      template: {
        filePath: './template.odt'
      },
      data: {
        image1: {
          filePath: './photo.png'
        }
      }
    }, { returnBase64: true });

    console.log(`✓ Received base64 string (${base64Data.length} chars)`);
    console.log(`✓ First 50 chars: ${base64Data.substring(0, 50)}...\n`);
  } catch (error) {
    console.error('✗ Error:', error.message, '\n');
  }
}

// Example 6: Health check and API info
async function example6() {
  console.log('Example 6: Health check and API info');
  console.log('=====================================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    // Health check
    const health = await client.healthCheck();
    console.log('Health check:', JSON.stringify(health, null, 2));

    // API info
    const info = await client.getInfo();
    console.log('\nAPI info:', JSON.stringify(info, null, 2));
    console.log();
  } catch (error) {
    console.error('✗ Error:', error.message);
    console.error('  Make sure the API server is running: ./odt-api\n');
  }
}

// Example 7: Error handling
async function example7() {
  console.log('Example 7: Error handling');
  console.log('=========================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');

  try {
    // This will fail because template is missing
    await client.replaceImages({
      template: {},
      data: {
        image1: { filePath: './photo.png' }
      }
    });
  } catch (error) {
    console.log('✓ Caught expected error:', error.message);
    console.log('  (This demonstrates proper error handling)\n');
  }
}

// Example 8: Using Buffer for template and images
async function example8() {
  console.log('Example 8: Using Buffer directly');
  console.log('=================================\n');

  const client = new ODTImageReplacerClient('http://localhost:8080');
  const fs = require('fs');

  try {
    // Read files into Buffers
    const templateBuffer = fs.readFileSync('./template.odt');
    const imageBuffer = fs.readFileSync('./photo.png');

    console.log(`Template Buffer: ${templateBuffer.length} bytes`);
    console.log(`Image Buffer: ${imageBuffer.length} bytes\n`);

    // Use Buffers directly - they will be auto-converted to base64
    const outputBuffer = await client.replaceImages({
      template: {
        buffer: templateBuffer  // Buffer will be converted to base64
      },
      data: {
        image1: {
          buffer: imageBuffer  // Buffer will be converted to base64
        }
      }
    });

    console.log(`✓ Generated output: ${outputBuffer.length} bytes`);
    console.log('✓ Buffers were automatically converted to base64\n');
  } catch (error) {
    console.error('✗ Error:', error.message);
    console.log('  (This is expected if files don\'t exist)\n');
  }
}

// Example 9: Custom API URL and timeout
async function example9() {
  console.log('Example 9: Custom configuration');
  console.log('================================\n');

  const client = new ODTImageReplacerClient('http://production-server:3000', {
    timeout: 60000  // 60 second timeout for large files
  });

  console.log('✓ Client configured with custom URL and timeout');
  console.log('  Base URL:', client.baseURL);
  console.log('  Timeout:', client.timeout, 'ms\n');
}

// Run all examples
async function runAllExamples() {
  console.log('╔════════════════════════════════════════════════╗');
  console.log('║  ODT Image Replacer - Node.js Client Examples ║');
  console.log('╚════════════════════════════════════════════════╝\n');

  // Note: Most examples will fail without actual files or running API
  // They are here to demonstrate usage patterns

  await example6();  // This one works if API is running
  await example7();  // This one demonstrates error handling
  await example8();  // This one shows Buffer usage
  await example9();  // This one just shows configuration

  console.log('════════════════════════════════════════════════');
  console.log('For working examples, ensure:');
  console.log('  1. API server is running: ./odt-api');
  console.log('  2. You have template.odt and image files');
  console.log('  3. Run: npm install (to install axios)');
  console.log('════════════════════════════════════════════════\n');
}

// Run if called directly
if (require.main === module) {
  runAllExamples().catch(console.error);
}

module.exports = {
  example1,
  example2,
  example3,
  example4,
  example5,
  example6,
  example7,
  example8,
  example9
};
