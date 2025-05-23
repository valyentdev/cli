#!/bin/sh
set -e
main() {
    response=$(curl -s "https://api.github.com/repos/valyentdev/cli/releases/latest")
    latest_version=$(echo "$response" | grep -m 1 '"name":' | awk -F'"' '{print $4}')
    latest_version_without_v=${latest_version#v}
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    # Map architecture names
    arch=$(uname -m)
    case $arch in
        x86_64)  arch="amd64" ;;
        aarch64) arch="arm64" ;;
        armv7l)  arch="arm" ;;
    esac
    
    version=${1:-$latest_version}
    valyent_dir="${VALYENT_DIR:-$HOME/.valyent}"
    bin_dir="$valyent_dir/bin"
    tmp_dir="$valyent_dir/tmp"
    exe="$bin_dir/valyent"
    
    mkdir -p "$bin_dir"
    mkdir -p "$tmp_dir"
    
    # Remove existing binary if it exists
    if [ -f "$exe" ]; then
        echo "Removing existing Valyent binary at $exe..."
        rm -f "$exe"
    fi

    download_url="https://github.com/valyentdev/cli/releases/download/$version/valyent_${latest_version_without_v}_${os}_$arch.tar.gz"
    echo "Downloading $download_url..."
    curl -q --fail --location --progress-bar --output "$tmp_dir/valyent.tar.gz" "$download_url"
    
    tar -C "$tmp_dir" -xzf "$tmp_dir/valyent.tar.gz"
    chmod +x "$tmp_dir/valyent"
    mv "$tmp_dir/valyent" "$exe"
    rm "$tmp_dir/valyent.tar.gz"
    
    echo "Valyent was installed successfully to $exe."
    
    if command -v valyent >/dev/null; then
        echo "Run \`valyent auth login\` to get started."
    else
        case $SHELL in
        /bin/zsh) shell_profile=".zshrc" ;;
        *) shell_profile=".bash_profile" ;;
        esac
        echo "\n# Valyent CLI" >> "$HOME/$shell_profile"
        echo "export VALYENT_INSTALL=\"$valyent_dir\"" >> "$HOME/$shell_profile"
        echo "export PATH=\"\$VALYENT_INSTALL/bin:\$PATH\"" >> "$HOME/$shell_profile"
        echo "Open a new terminal or run 'source $HOME/$shell_profile' to start using Valyent CLI"
        echo "Then, run \`valyent auth login\` to get started."
    fi
}

main "$1"