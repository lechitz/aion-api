# Ollama GPU Acceleration (Optional)

Enable this only if your host has NVIDIA GPU support and you need faster local inference.

## Prerequisites

- NVIDIA GPU with updated driver
- NVIDIA Container Toolkit installed

## Install Toolkit (Ubuntu)

```bash
distribution=$(. /etc/os-release; echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/libnvidia-container/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/libnvidia-container.list \
  | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker
```

## Validate GPU from Docker

```bash
docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi
```

## Enable GPU in Dev Compose

Edit `infrastructure/docker/environments/dev/docker-compose-dev.yaml` and add the GPU reservation block under `ollama`.

```yaml
ollama:
  image: ollama/ollama:latest
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: 1
            capabilities: [gpu]
```

## Restart and Verify

```bash
make dev-down
make dev
docker exec aion-dev-ollama nvidia-smi
docker logs aion-dev-ollama | grep -i cuda
```

If CUDA is detected, GPU acceleration is active.
