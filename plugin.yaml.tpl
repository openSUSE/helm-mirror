name: "mirror"
# plugin
version: "foo"
usage: "Mirror Helm Charts from a repository into a local folder."
description: "Mirror Helm Charts from a repository into a local folder."
useTunnel: true
command: "$HELM_PLUGIN_DIR/bin/mirror"
hooks:
  install: "$HELM_PLUGIN_DIR/scripts/install-binary.sh"
  update: "$HELM_PLUGIN_DIR/scripts/install-binary.sh"
