terraform {
  backend "s3" {
    bucket = "leone-terraform-states"
    key    = "expense-system.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}
