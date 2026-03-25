output "namespace" {
  value = kubernetes_namespace.event_service.metadata[0].name
}

output "helm_release_name" {
  value = helm_release.event_service.name
}
