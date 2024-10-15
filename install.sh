#!/bin/sh

set -e

main() {
  response=$(curl -s "https://api.github.com/repos/valyentdev/cli/releases/latest")
  latest_version=$(echo "$response" | grep -m 1 '"name":' | awk -F'"' '{print $4}')
  os=$(uname -s)
  arch=$(uname -m)
  version=${1:-$latest_version}

  citadel_dir="${CITADEL_DIR:-$HOME/.citadel}"

	bin_dir="$citadel_dir/bin"
	tmp_dir="$citadel_dir/tmp"
	exe="$bin_dir/citadel"

	mkdir -p "$bin_dir"
	mkdir -p "$tmp_dir"
  
  download_url="https://github.com/valyentdev/cli/releases/download/$version/cli_${os}_$arch.tar.gz"
	echo "Downloading $download_url..."
  curl -q --fail --location --progress-bar --output "$tmp_dir/citadel.tar.gz" "$download_url"
	tar -C "$tmp_dir" -xzf "$tmp_dir/citadel.tar.gz"
	chmod +x "$tmp_dir/citadel"
	mv "$tmp_dir/citadel" "$exe"
	rm "$tmp_dir/citadel.tar.gz"

	echo "Software Citadel CLI was installed successfully to $exe."

	if command -v citadel >/dev/null; then
		echo "Run \`citadel auth login\` to get started."
	else
		case $SHELL in
		/bin/zsh) shell_profile=".zshrc" ;;
		*) shell_profile=".bash_profile" ;;
		esac

    echo "\n# Software Citadel CLI" >> "$HOME/$shell_profile"
    echo "export CITADEL_INSTALL=\"$citadel_dir\"" >> "$HOME/$shell_profile"
    echo "export PATH=\"\$CITADEL_INSTALL/bin:\$PATH\"" >> "$HOME/$shell_profile"

    echo "Open a new terminal or run 'source $HOME/$shell_profile' to start using Software Citadel CLI"

		echo "Then, run \`citadel auth login\` to get started."
	fi
}

main "$1"