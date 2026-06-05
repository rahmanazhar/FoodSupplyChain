#!/usr/bin/env bash
#
# run.sh — launch the full Food Supply Chain stack locally.
#
# Brings up Postgres + NATS (Docker), then builds and runs the inventory,
# shipment and gateway services. Everything is torn down on Ctrl-C.
#
# Usage:
#   ./scripts/run.sh                start the whole stack (Ctrl-C to stop)
#   ./scripts/run.sh --seed         also insert a demo product + location
#   ./scripts/run.sh --no-docker    assume Postgres/NATS are already running
#   ./scripts/run.sh --help         show this help
#
set -euo pipefail

# Resolve the repo root so the script works from any directory.
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

COMPOSE_FILE="deployments/docker/docker-compose.yml"
COMPOSE=(docker compose -f "$COMPOSE_FILE")
LOG_DIR="$ROOT/logs"
BIN_DIR="$ROOT/bin"

INVENTORY_PORT=8080
SHIPMENT_PORT=8081
GATEWAY_PORT=3000

SEED=false
USE_DOCKER=true
for arg in "$@"; do
  case "$arg" in
    --seed)      SEED=true ;;
    --no-docker) USE_DOCKER=false ;;
    -h|--help)   sed -n '2,13p' "$0" | sed 's/^#\s\{0,1\}//'; exit 0 ;;
    *) echo "unknown option: $arg (try --help)" >&2; exit 1 ;;
  esac
done

PIDS=()
CLEANED=false

cleanup() {
  $CLEANED && return 0
  CLEANED=true
  echo ""
  echo "==> Stopping services..."
  for pid in "${PIDS[@]:-}"; do
    [ -n "$pid" ] && kill "$pid" 2>/dev/null || true
  done
  if $USE_DOCKER; then
    echo "==> Stopping Docker (Postgres + NATS)..."
    "${COMPOSE[@]}" down
  fi
  echo "==> Done."
}
trap cleanup EXIT INT TERM

mkdir -p "$LOG_DIR" "$BIN_DIR"

# 1. Infrastructure ----------------------------------------------------------
if $USE_DOCKER; then
  echo "==> Starting Postgres + NATS..."
  "${COMPOSE[@]}" up -d postgres nats
  echo "==> Waiting for Postgres..."
  for i in $(seq 1 60); do
    if "${COMPOSE[@]}" exec -T postgres pg_isready -U supplychain >/dev/null 2>&1; then
      echo "    Postgres ready."
      break
    fi
    sleep 1
    if [ "$i" -eq 60 ]; then
      echo "Postgres did not become ready in time" >&2
      exit 1
    fi
  done
fi

# 2. Build -------------------------------------------------------------------
echo "==> Building services..."
go build -o "$BIN_DIR/inventory" ./cmd/inventory
go build -o "$BIN_DIR/shipment"  ./cmd/shipment
go build -o "$BIN_DIR/gateway"   ./cmd/gateway

# 3. Launch ------------------------------------------------------------------
echo "==> Launching services..."
"$BIN_DIR/inventory" >"$LOG_DIR/inventory.log" 2>&1 &
PIDS+=($!)
SERVER_PORT=$SHIPMENT_PORT "$BIN_DIR/shipment" >"$LOG_DIR/shipment.log" 2>&1 &
PIDS+=($!)
"$BIN_DIR/gateway" >"$LOG_DIR/gateway.log" 2>&1 &
PIDS+=($!)

# 4. Wait for readiness ------------------------------------------------------
wait_health() {
  local name=$1 url=$2
  for _ in $(seq 1 30); do
    if curl -fsS -o /dev/null "$url" 2>/dev/null; then
      echo "    $name ready"
      return 0
    fi
    sleep 0.5
  done
  echo "    $name did NOT start — see $LOG_DIR/$name.log" >&2
  return 1
}
echo "==> Waiting for services..."
wait_health inventory "http://localhost:$INVENTORY_PORT/health"
wait_health shipment  "http://localhost:$SHIPMENT_PORT/health"
wait_health gateway   "http://localhost:$GATEWAY_PORT/health"

# 5. Optional demo seed ------------------------------------------------------
if $SEED && $USE_DOCKER; then
  echo "==> Seeding demo product + location..."
  "${COMPOSE[@]}" exec -T postgres psql -U supplychain -d supplychain -q -c \
    "INSERT INTO products (id,name,sku,category,unit_price,created_at,updated_at) VALUES ('prod-1','Apples','SKU-APPLE','fruits',1.99,now(),now()) ON CONFLICT DO NOTHING;
     INSERT INTO locations (id,name,type,created_at,updated_at) VALUES ('loc-1','Main Warehouse','warehouse',now(),now()) ON CONFLICT DO NOTHING;" \
    && echo "    seeded prod-1 / loc-1"
fi

cat <<EOF

============================================================
  Food Supply Chain is running.

  Gateway    http://localhost:$GATEWAY_PORT
  Inventory  http://localhost:$INVENTORY_PORT          (no auth)
  Shipment   http://localhost:$SHIPMENT_PORT/api/v1     (JWT required)

  Mint a token:
    go run ./cmd/token -secret "your-secret-key-here" -role admin

  Logs:  $LOG_DIR/{inventory,shipment,gateway}.log

  Press Ctrl-C to stop everything.
============================================================
EOF

# Stay in the foreground; exit (and clean up) if any service dies.
wait
