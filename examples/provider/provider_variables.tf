variable "loft_host" {
}

variable "loft_access_key" {
}

variable "loft_insecure" {
}

provider "loft" {
  host       = var.loft_host
  access_key = var.loft_access_key
  insecure   = var.loft_insecure
}
