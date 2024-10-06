locals {
  function_name = "bucket-checker"
}

provider "google" {
  region = "eu-west-1"
}

resource "google_project_service" "photo_upload_api" {
  service = "run.googleapis.com"
}

resource "google_cloud_v2_service" "photo_upload_cloud_run_service" {
  name = "photo-upload-cloud-run-service"
  location = "eu-west-1"

  labels = {
    service = "photo-upload" // For billing purposes
  }

  template {
    spec {
      containers {
        image = "docker.io/docker_username/interview_devops:latest"
      }
    }
  }
}

resource "google_cloud_run_service_iam_member" "photo_upload_invoker" {
  location = google_cloud_run_service.photo_upload_cloud_run_service.location
  project = google_cloud_run_service.photo_upload_cloud_run_service.project
  service = google_cloud_run_service.photo_upload_cloud_run_service.name
  role = "roles/run.invoker"
  member = "allUsers"
}

output "service_url" {
  value = google_cloud_run_service.photo_upload_cloud_run_service.status[0].url
}