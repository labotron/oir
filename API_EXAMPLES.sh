#!/bin/bash

# ODT Image Replacer API - Example Commands
# ==========================================

echo "ODT Image Replacer API - Example Commands"
echo "=========================================="
echo ""

BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print command
print_command() {
    echo -e "${BLUE}➜${NC} $1"
    echo ""
}

# Function to print section
print_section() {
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════${NC}"
    echo ""
}

# 1. Health Check
print_section "1. Health Check"
print_command "curl -X GET $BASE_URL/health"
curl -X GET $BASE_URL/health
echo ""

# 2. Service Info
print_section "2. Service Information"
print_command "curl -X GET $BASE_URL/info"
curl -X GET $BASE_URL/info | jq .
echo ""

# 3. Replace Images with URL sources
print_section "3. Replace Images (URL Sources)"
print_command "curl -X POST $BASE_URL/api/replace -H 'Content-Type: application/json' -d '{...}'"

cat > /tmp/request_url.json <<'EOF'
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  },
  "data": {
    "image1": {
      "url": "https://picsum.photos/200/300",
      "base64": null
    }
  }
}
EOF

echo "Request JSON:"
cat /tmp/request_url.json | jq .
echo ""

# Uncomment to actually send the request
# curl -X POST $BASE_URL/api/replace \
#   -H "Content-Type: application/json" \
#   -d @/tmp/request_url.json | jq .

# 4. Replace Images with Base64 sources
print_section "4. Replace Images (Base64 Sources)"
print_command "Using base64-encoded template and images"

cat > /tmp/request_base64.json <<'EOF'
{
  "template": {
    "url": null,
    "base64": "UEsDBBQAAAAIAOB/Y1n5H7IWAAAA..."
  },
  "data": {
    "image1": {
      "url": null,
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAUA..."
    }
  }
}
EOF

echo "Request JSON:"
cat /tmp/request_base64.json | jq .
echo ""

# 5. Download ODT Directly
print_section "5. Replace and Download ODT File"
print_command "curl -X POST $BASE_URL/api/replace/download -d @request.json -o output.odt"

cat > /tmp/request_download.json <<'EOF'
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  },
  "data": {
    "image1": {
      "url": "https://picsum.photos/200/300",
      "base64": null
    },
    "image2": {
      "url": "https://picsum.photos/300/200",
      "base64": null
    }
  }
}
EOF

# Uncomment to download
# curl -X POST $BASE_URL/api/replace/download \
#   -H "Content-Type: application/json" \
#   -d @/tmp/request_download.json \
#   -o output.odt

# 6. Get JSON Response and Decode Base64
print_section "6. Get JSON Response and Decode Base64"
print_command "curl ... | jq -r '.output_base64' | base64 -d > output.odt"

# Uncomment to run
# curl -X POST $BASE_URL/api/replace \
#   -H "Content-Type: application/json" \
#   -d @/tmp/request_url.json | \
#   jq -r '.output_base64' | \
#   base64 -d > output.odt

# 7. Multiple Image Replacements
print_section "7. Replace Multiple Images"

cat > /tmp/request_multiple.json <<'EOF'
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  },
  "data": {
    "logo": {
      "url": "https://example.com/logo.png",
      "base64": null
    },
    "signature": {
      "url": "https://example.com/signature.png",
      "base64": null
    },
    "photo1": {
      "url": "https://example.com/photo1.jpg",
      "base64": null
    },
    "photo2": {
      "url": null,
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAUA..."
    }
  }
}
EOF

echo "Request JSON:"
cat /tmp/request_multiple.json | jq .
echo ""

# 8. Using with jq for pretty output
print_section "8. Pretty JSON Output with jq"
print_command "curl ... | jq '{success, message, replaced_tags}'"

# 9. Error handling example
print_section "9. Error Handling"
print_command "curl -X POST $BASE_URL/api/replace -d '{}'"

# Uncomment to test error
# curl -X POST $BASE_URL/api/replace \
#   -H "Content-Type: application/json" \
#   -d '{}' | jq .

# 10. Complete workflow example
print_section "10. Complete Workflow Example"

echo "Step 1: Create your request JSON"
echo "Step 2: Send request to API"
echo "Step 3: Receive and save result"
echo ""

cat > /tmp/workflow.sh <<'WORKFLOW'
#!/bin/bash

# Create request
cat > request.json <<'EOF'
{
  "template": {
    "url": "https://example.com/report.odt",
    "base64": null
  },
  "data": {
    "employee_photo": {
      "url": "https://example.com/john_doe.jpg",
      "base64": null
    },
    "company_logo": {
      "url": "https://example.com/logo.png",
      "base64": null
    }
  }
}
EOF

# Send request and download result
curl -X POST http://localhost:8080/api/replace/download \
  -H "Content-Type: application/json" \
  -d @request.json \
  -o report_final.odt

echo "✓ Report generated: report_final.odt"
WORKFLOW

chmod +x /tmp/workflow.sh
echo "Example workflow script created at: /tmp/workflow.sh"
cat /tmp/workflow.sh
echo ""

# 11. Testing with local files
print_section "11. Using Local Template and Images"

echo "Convert local template to base64:"
print_command "base64 -i template.odt | tr -d '\n'"
echo ""

echo "Convert local image to base64:"
print_command "base64 -i photo.png | tr -d '\n'"
echo ""

# 12. Batch processing example
print_section "12. Batch Processing Multiple Documents"

cat > /tmp/batch_process.sh <<'BATCH'
#!/bin/bash

# Process multiple employees
EMPLOYEES=("john" "jane" "bob" "alice")

for emp in "${EMPLOYEES[@]}"; do
  echo "Processing $emp..."

  cat > request_$emp.json <<EOF
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  },
  "data": {
    "employee_photo": {
      "url": "https://example.com/photos/${emp}.jpg",
      "base64": null
    }
  }
}
EOF

  curl -X POST http://localhost:8080/api/replace/download \
    -H "Content-Type: application/json" \
    -d @request_$emp.json \
    -o report_${emp}.odt

  echo "✓ Generated: report_${emp}.odt"
done

echo "✓ Batch processing complete!"
BATCH

chmod +x /tmp/batch_process.sh
echo "Batch processing script created at: /tmp/batch_process.sh"
cat /tmp/batch_process.sh
echo ""

# Summary
print_section "Summary"
echo "Example files created in /tmp/:"
echo "  - request_url.json"
echo "  - request_base64.json"
echo "  - request_download.json"
echo "  - request_multiple.json"
echo "  - workflow.sh"
echo "  - batch_process.sh"
echo ""
echo "To start the server:"
echo "  ./odt-api"
echo ""
echo "Then uncomment the curl commands in this script to test!"
echo ""
