# Kubernetes Manifests

Plain Kubernetes manifests for the Food Supply Chain stack. They deploy
PostgreSQL, NATS (JetStream), the inventory and shipment services, and the API
gateway into the `foodsupplychain` namespace.

## Prerequisites

- A Kubernetes cluster (kind, minikube, or a managed cluster)
- The service images built and available to the cluster:

  ```bash
  make docker-build
  # For kind: kind load docker-image foodsupplychain-inventory:latest \
  #   foodsupplychain-shipment:latest foodsupplychain-gateway:latest
  ```

## Apply

Files are numbered for ordering — apply the whole directory:

```bash
kubectl apply -f deployments/k8s/
```

Or step by step:

```bash
kubectl apply -f deployments/k8s/00-namespace.yaml
kubectl apply -f deployments/k8s/01-config.yaml
kubectl apply -f deployments/k8s/02-postgres.yaml
kubectl apply -f deployments/k8s/03-nats.yaml
kubectl apply -f deployments/k8s/04-inventory.yaml
kubectl apply -f deployments/k8s/05-shipment.yaml
kubectl apply -f deployments/k8s/06-gateway.yaml
```

## Access

```bash
kubectl -n foodsupplychain get pods
kubectl -n foodsupplychain port-forward svc/api-gateway 3000:80
curl http://localhost:3000/health
```

## Notes

- `01-config.yaml` contains **development placeholder** secrets. Replace
  `JWT_SECRET` and `DB_PASSWORD` with values from a real secret manager
  (e.g. Sealed Secrets, External Secrets, or your cloud provider) before any
  non-local deployment.
- The service images bake `configs/config.yaml`; the ConfigMap/Secret only
  override the database, NATS, and JWT settings via environment variables.
