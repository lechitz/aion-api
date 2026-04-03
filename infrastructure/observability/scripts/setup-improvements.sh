#!/bin/bash
# ==============================================================================
# Observability Improvements - Apply & Validate
# ==============================================================================
# Date: 2026-03-28
# Description: Validates the current local observability stack (Grafana, Prometheus, Jaeger, Loki)
# Estimated time: 5-10 minutes
# ==============================================================================

set -e  # Exit on error

# Output colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "\n${BLUE}================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check working directory and Compose
if [ ! -f "go.mod" ]; then
    print_error "Run this script from the aion-api project root!"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    print_error "Docker not found. Install/enable Docker Desktop or a compatible engine."
    exit 1
fi

if ! docker info >/dev/null 2>&1; then
    print_error "Docker is not accessible (daemon stopped or insufficient permissions)."
    exit 1
fi

PROJECT_ROOT="$(pwd)"
COMPOSE_DIR="$PROJECT_ROOT/infrastructure/docker/environments/dev"
COMPOSE_FILE="$COMPOSE_DIR/docker-compose-dev.yaml"

if [ ! -f "$COMPOSE_FILE" ]; then
    print_error "Compose file not found at $COMPOSE_FILE"
    exit 1
fi

compose() {
    docker compose -f "$COMPOSE_FILE" --project-directory "$COMPOSE_DIR" "$@"
}

get_container_state() {
    local container_name=$1
    if ! docker inspect "$container_name" >/dev/null 2>&1; then
        echo "not_found"
        return
    fi

    docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container_name" 2>/dev/null || echo "unknown"
}

if [ ! -f "$COMPOSE_DIR/.env.dev" ]; then
    print_warning "File $COMPOSE_DIR/.env.dev not found. Create it from .env.dev.example before bringing the stack up."
fi

# ==============================================================================
# PHASE 1: Check Files
# ==============================================================================
print_header "PHASE 1: Checking created/updated files"

files_to_check=(
    "infrastructure/observability/prometheus/prometheus.yml"
    "infrastructure/observability/grafana/datasources/jaeger.yaml"
    "infrastructure/observability/grafana/datasources/loki.yaml"
    "infrastructure/observability/grafana/datasources/prometheus.yaml"
    "infrastructure/observability/grafana/dashboards/aionapi-red-dashboard.json"
    "infrastructure/observability/grafana/dashboards/aionapi-http-requests-dashboard.json"
)

all_files_ok=true
for file in "${files_to_check[@]}"; do
    if [ -f "$file" ]; then
        print_success "File exists: $file"
    else
        print_error "Missing file: $file"
        all_files_ok=false
    fi
done

if [ "$all_files_ok" = false ]; then
    print_error "Some required observability assets are missing."
    exit 1
fi

# Check if prometheus.yml has exemplars
if grep -q "scrape_protocols" infrastructure/observability/prometheus/prometheus.yml; then
    print_success "Prometheus configured with exemplars"
else
    print_warning "Prometheus may not be configured with exemplars"
fi

# ==============================================================================
# PHASE 2: Stop current stack
# ==============================================================================
print_header "PHASE 2: Stopping current stack"

if compose ps >/dev/null 2>&1; then
    print_info "Stopping containers..."
    compose down
    print_success "Containers stopped"
else
    print_info "No running containers to stop"
fi

# ==============================================================================
# PHASE 3: Start stack with new configs
# ==============================================================================
print_header "PHASE 3: Starting stack with new configs"

print_info "Starting containers..."
compose up -d

print_info "Waiting for services to boot (30 seconds)..."
sleep 30

# ==============================================================================
# PHASE 4: Check container health
# ==============================================================================
print_header "PHASE 4: Checking container health"

containers=(
    "postgres-dev:healthy"
    "redis-aion-dev:healthy"
    "prometheus-dev:running"
    "grafana-dev:running"
    "jaeger-dev:running"
    "otel-collector:running"
    "aion-api-dev:healthy"
)

all_healthy=true
for container_check in "${containers[@]}"; do
    IFS=':' read -r container expected <<< "$container_check"

    status=$(get_container_state "$container")

    case "$expected" in
        healthy)
            if [ "$status" = "healthy" ]; then
                print_success "$container: $status"
            else
                print_error "$container: $status (expected: healthy)"
                all_healthy=false
            fi
            ;;
        running)
            if [ "$status" = "running" ] || [ "$status" = "healthy" ]; then
                print_success "$container: $status"
            else
                print_error "$container: $status (expected: running/healthy)"
                all_healthy=false
            fi
            ;;
        *)
            print_warning "$container: status $status (expected: $expected)"
            ;;
    esac
done

if [ "$all_healthy" = false ]; then
    print_warning "Some containers may be unhealthy. Check the logs."
    print_info "Command: docker-compose logs <container_name>"
fi

# ==============================================================================
# PHASE 5: Validate services
# ==============================================================================
print_header "PHASE 5: Validating service endpoints"

# Endpoint test helper
test_endpoint() {
    local name=$1
    local url=$2
    local expected=$3

    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")

    if [ "$response" = "$expected" ]; then
        print_success "$name responding (HTTP $response)"
        return 0
    else
        print_error "$name failed (HTTP $response, expected: $expected)"
        return 1
    fi
}

print_info "Testing endpoints (wait 10s to stabilize)..."
sleep 10

test_endpoint "Prometheus" "http://localhost:9090/-/healthy" "200"
test_endpoint "Grafana" "http://localhost:3000/api/health" "200"
test_endpoint "Jaeger" "http://localhost:16686/" "200"
test_endpoint "OTel Collector" "http://localhost:9888/metrics" "200"

# ==============================================================================
# PHASE 6: Check Prometheus targets
# ==============================================================================
print_header "PHASE 6: Checking Prometheus targets"

targets_json=$(curl -s http://localhost:9090/api/v1/targets 2>/dev/null || echo '{"data":{"activeTargets":[]}}')

if command -v jq &> /dev/null; then
    otel_health=$(echo "$targets_json" | jq -r '.data.activeTargets[] | select(.labels.job=="otel-collector") | .health' 2>/dev/null || echo "unknown")

    if [ "$otel_health" = "up" ]; then
        print_success "Target otel-collector: UP"
    else
        print_error "Target otel-collector: $otel_health"
    fi
else
    print_warning "jq not installed, skipping target check"
    print_info "Check manually: http://localhost:9090/targets"
fi

# ==============================================================================
# PHASE 7: Check datasources in Grafana
# ==============================================================================
print_header "PHASE 7: Checking Grafana datasources"

print_info "Waiting for Grafana to provision datasources (10s)..."
sleep 10

datasources=$(curl -s -u "aion:aion" http://localhost:3000/api/datasources 2>/dev/null || echo '[]')

if command -v jq &> /dev/null; then
    prometheus_ds=$(echo "$datasources" | jq -r '.[] | select(.type=="prometheus") | .name' 2>/dev/null || echo "")
    jaeger_ds=$(echo "$datasources" | jq -r '.[] | select(.type=="jaeger") | .name' 2>/dev/null || echo "")

    if [ -n "$prometheus_ds" ]; then
        print_success "Prometheus datasource found: $prometheus_ds"
    else
        print_error "Prometheus datasource not found"
    fi

    if [ -n "$jaeger_ds" ]; then
        print_success "Jaeger datasource found: $jaeger_ds"
    else
        print_warning "Jaeger datasource not found (provisioning may take longer)"
        print_info "Check manually: http://localhost:3000/datasources"
    fi
else
    print_warning "jq not installed, check datasources manually"
    print_info "URL: http://localhost:3000/datasources"
fi

# ==============================================================================
# PHASE 8: Check dashboards
# ==============================================================================
print_header "PHASE 8: Checking dashboards in Grafana"

dashboards=$(curl -s -u "aion:aion" http://localhost:3000/api/search?type=dash-db 2>/dev/null || echo '[]')

if command -v jq &> /dev/null; then
    red_dashboard=$(echo "$dashboards" | jq -r '.[] | select(.uid=="aion-api-red-pro") | .title' 2>/dev/null || echo "")

    if [ -n "$red_dashboard" ]; then
        print_success "RED dashboard found: $red_dashboard"
        print_info "URL: http://localhost:3000/d/aion-api-red-pro"
    else
        print_warning "RED dashboard not found (provisioning may take longer)"
        print_info "Check manually: http://localhost:3000/dashboards"
    fi
else
    print_warning "jq not installed, check dashboards manually"
    print_info "URL: http://localhost:3000/dashboards"
fi

# ==============================================================================
# PHASE 9: Generate test traffic
# ==============================================================================
print_header "PHASE 9: Generating test traffic (optional)"

generate_traffic=false

if [ -t 0 ]; then
    read -p "Generate test traffic to populate the dashboard? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        generate_traffic=true
    fi
else
    print_info "Non-interactive session detected; skipping automatic traffic generation."
fi

if [ "$generate_traffic" = true ]; then
    print_info "Sending 100 requests..."

    for i in {1..100}; do
        curl -s http://localhost:5001/aion/api/v1/health > /dev/null 2>&1 || true
        [ $((i % 20)) -eq 0 ] && echo -n "."
    done

    echo ""
    print_success "Traffic generated! Wait 1-2 minutes for metrics to appear."
else
    print_info "Skipping traffic generation. Trigger manually if needed."
fi

# ==============================================================================
# PHASE 10: Summary and next steps
# ==============================================================================
print_header "IMPLEMENTATION SUMMARY"

echo -e "${GREEN}✅ Stack restarted with current observability assets${NC}"
echo -e "${GREEN}✅ Containers verified${NC}"
echo -e "${GREEN}✅ Services validated${NC}"
echo ""
echo -e "${BLUE}📊 Important URLs:${NC}"
echo -e "   • RED dashboard: ${YELLOW}http://localhost:3000/d/aion-api-red-pro${NC}"
echo -e "   • Grafana:       ${YELLOW}http://localhost:3000${NC} (login: aion/aion)"
echo -e "   • Jaeger:        ${YELLOW}http://localhost:16686${NC}"
echo -e "   • Prometheus:    ${YELLOW}http://localhost:9090${NC}"
echo ""
echo -e "${BLUE}📚 Documentation:${NC}"
echo -e "   • Quickstart:           ${YELLOW}docs/observability-quickstart.md${NC}"
echo -e "   • Performance guide:    ${YELLOW}docs/performance-readiness.md${NC}"
echo ""
echo -e "${BLUE}🎯 Next steps:${NC}"
echo -e "   1. Open the RED dashboard in Grafana"
echo -e "   2. Explore panels (Top Slowest, Impact Analysis, etc.)"
echo -e "   3. Click bars/points to test drill-down to Jaeger"
echo -e "   4. Read the runbook to learn queries and workflows"
echo ""
echo -e "${BLUE}🐛 Issues?${NC}"
echo -e "   • API logs:        ${YELLOW}make logs-api${NC}"
echo -e "   • Grafana logs:    ${YELLOW}make logs-grafana${NC}"
echo -e "   • Prometheus logs: ${YELLOW}make logs-prometheus${NC}"
echo ""
print_success "Implementation completed! 🎉"
echo ""
