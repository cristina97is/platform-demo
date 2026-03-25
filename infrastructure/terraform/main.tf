resource "kubernetes_namespace" "event_service" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "event_service" {
  name      = var.helm_release_name
  namespace = kubernetes_namespace.event_service.metadata[0].name
  chart     = var.helm_chart_path

  depends_on = [
    kubernetes_namespace.event_service
  ]

  set {
    name  = "namespace"
    value = var.namespace
  }

  set {
    name  = "api.image.repository"
    value = var.image_repository
  }

  set {
    name  = "api.image.tag"
    value = var.image_tag
  }
}
