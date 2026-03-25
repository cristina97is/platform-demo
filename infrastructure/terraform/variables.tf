variable "namespace" {
  description = "Kubernetes namespace for event-service"
  type        = string
  default     = "event-service"
}

variable "helm_release_name" {
  description = "Helm release name"
  type        = string
  default     = "event-service"
}

variable "helm_chart_path" {
  description = "Path to Helm chart"
  type        = string
  default     = "../../event-service-chart"
}

variable "image_repository" {
  description = "Container image repository"
  type        = string
  default     = "ghcr.io/cristina97is/event-service"
}

variable "image_tag" {
  description = "Container image tag. latest is fallback only."
  type        = string
  default     = "latest"
}
