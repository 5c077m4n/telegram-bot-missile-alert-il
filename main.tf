variable "PROJECT_ID" {
  type = string
}
variable "REGION" {
  type    = string
  default = "me-west1"
}
variable "TELEGRAM_BOT_API_TOKEN" {
  type      = string
  sensitive = true
}
variable "MONGO_URI" {
  type      = string
  sensitive = true
  default   = "mongodb://mongodb:27017"
}
variable "MONGO_DB" {
  type      = string
  sensitive = true
  default   = "missile_alert"
}

provider "google" {
  project = var.PROJECT_ID
  region  = var.REGION
}

resource "google_service_account" "instance_service_account" {
  account_id   = "vm-service-account"
  display_name = "Service Account for VM"
}

resource "google_secret_manager_secret" "TELEGRAM_BOT_API_TOKEN" {
  secret_id = "TELEGRAM_BOT_API_TOKEN"
  replication {
    auto {}
  }
}
resource "google_secret_manager_secret_version" "telegram_bot_api_token_data" {
  secret      = google_secret_manager_secret.TELEGRAM_BOT_API_TOKEN.id
  secret_data = var.TELEGRAM_BOT_API_TOKEN
}
resource "google_secret_manager_secret_iam_member" "telegram_bot_api_token_access" {
  secret_id = google_secret_manager_secret.TELEGRAM_BOT_API_TOKEN.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.instance_service_account.email}"
}

resource "google_secret_manager_secret" "MONGO_URI" {
  secret_id = "MONGO_URI"
  replication {
    auto {}
  }
}
resource "google_secret_manager_secret_version" "mongo_uri_data" {
  secret      = google_secret_manager_secret.MONGO_URI.id
  secret_data = var.MONGO_URI
}
resource "google_secret_manager_secret_iam_member" "mongo_uri_access" {
  secret_id = google_secret_manager_secret.MONGO_URI.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.instance_service_account.email}"
}

resource "google_secret_manager_secret" "MONGO_DB" {
  secret_id = "MONGO_DB"
  replication {
    auto {}
  }
}
resource "google_secret_manager_secret_version" "mongo_db_data" {
  secret      = google_secret_manager_secret.MONGO_DB.id
  secret_data = var.MONGO_DB
}
resource "google_secret_manager_secret_iam_member" "mongo_db_access" {
  secret_id = google_secret_manager_secret.MONGO_DB.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.instance_service_account.email}"
}

resource "google_compute_instance" "vm_instance" {
  name         = "telegram-bot-missile-alert-il"
  machine_type = "e2-medium"
  zone         = "${var.REGION}-c"

  tags = ["telegram-bot"]

  boot_disk {
    initialize_params { image = "ubuntu-2404-lts-amd64" }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  service_account {
    email  = google_service_account.instance_service_account.email
    scopes = ["cloud-platform"]
  }

  metadata_startup_script = <<-EOT
    #! /bin/bash
    set -ueo pipefail

    sudo apt update
    sudo apt install --yes git neovim curl ufw fail2ban ca-certificates curl gnupg tar
    sudo ufw default deny incoming
    sudo ufw allow OpenSSH
    sudo ufw --force enable
    sudo systemctl enable fail2ban --now
    export EDITOR=nvim

    curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
    sudo sh /tmp/get-docker.sh
    rm -f /tmp/get-docker.sh

    mkdir --parent ~/workspace/
    cd ~/workspace/
    git clone https://github.com/5c077m4n/telegram-bot-missile-alert-il.git
    cd telegram-bot-missile-alert-il/

    curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | sudo -E bash
    sudo apt install task

    echo "ENV=prod" >>.env
    echo "TELEGRAM_BOT_API_TOKEN=$(gcloud secrets versions access latest --secret="TELEGRAM_BOT_API_TOKEN")" >>.env
    echo "MONGO_URI=$(gcloud secrets versions access latest --secret="MONGO_URI")" >>.env
    echo "MONGO_DB=$(gcloud secrets versions access latest --secret="MONGO_DB")" >>.env

    task compose
  EOT
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh"
  network = google_compute_instance.vm_instance.network_interface[0].network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = google_compute_instance.vm_instance.tags
}
resource "google_compute_firewall" "deny_all" {
  name      = "deny-all-ingress"
  network   = google_compute_instance.vm_instance.network_interface[0].network
  priority  = 65534
  direction = "INGRESS"

  deny { protocol = "all" }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = google_compute_instance.vm_instance.tags
}
