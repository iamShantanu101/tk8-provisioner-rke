output "ssh_key_path" {
  value = ["${aws_key_pair.rke-node-key.id}"]
}

output "ssh_username" {
  value = "ubuntu"
}

output "addresses" {
  value = ["${aws_instance.rke-node.*.public_dns}"]
}
