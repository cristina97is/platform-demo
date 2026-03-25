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
