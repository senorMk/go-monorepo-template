#!/usr/bin/env bash
set -euo pipefail

# -------------------------------------------------------
# Fullstack Monorepo Template — init script
# Replaces all placeholder tokens with real project values
# -------------------------------------------------------

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo ""
echo -e "${BOLD}Fullstack Monorepo Template — Project Setup${NC}"
echo "--------------------------------------------"
echo ""

# -- Prompt --

read -rp "$(echo -e "${CYAN}App name${NC} (kebab-case, e.g. my-app): ")" APP_NAME
if [[ -z "$APP_NAME" ]]; then echo -e "${RED}App name is required.${NC}"; exit 1; fi

read -rp "$(echo -e "${CYAN}Display name${NC} (e.g. My App): ")" APP_DISPLAY_NAME
if [[ -z "$APP_DISPLAY_NAME" ]]; then echo -e "${RED}Display name is required.${NC}"; exit 1; fi

read -rp "$(echo -e "${CYAN}Domain${NC} (e.g. myapp.com): ")" APP_DOMAIN
if [[ -z "$APP_DOMAIN" ]]; then echo -e "${RED}Domain is required.${NC}"; exit 1; fi

read -rp "$(echo -e "${CYAN}GitHub username${NC} (e.g. johndoe): ")" GITHUB_USERNAME
if [[ -z "$GITHUB_USERNAME" ]]; then echo -e "${RED}GitHub username is required.${NC}"; exit 1; fi

# -- Derived values --
# snake_case from kebab-case  (my-app -> my_app)
APP_SNAKE="${APP_NAME//-/_}"
# PascalCase from kebab-case  (my-app -> MyApp)
APP_PASCAL="$(echo "$APP_NAME" | sed 's/-\([a-z]\)/\U\1/g; s/^\([a-z]\)/\U\1/')"

echo ""
echo -e "${BOLD}Applying:${NC}"
echo "  APP_NAME         = $APP_NAME"
echo "  APP_DISPLAY_NAME = $APP_DISPLAY_NAME"
echo "  APP_DOMAIN       = $APP_DOMAIN"
echo "  APP_SNAKE        = $APP_SNAKE"
echo "  APP_PASCAL       = $APP_PASCAL"
echo "  GITHUB_USERNAME  = $GITHUB_USERNAME"
echo ""

# -- Find all text files to process (skip .git and node_modules) --
FILES=$(find "$ROOT" \
  -not -path '*/.git/*' \
  -not -path '*/node_modules/*' \
  -not -path '*/build/*' \
  -not -path '*/.dart_tool/*' \
  -not -name '*.png' \
  -not -name '*.jpg' \
  -not -name '*.jpeg' \
  -not -name '*.ico' \
  -not -name '*.ttf' \
  -not -name '*.otf' \
  -not -name '*.woff' \
  -not -name '*.woff2' \
  -not -name '*.pbxproj' \
  -not -name 'init.sh' \
  -type f)

replace_in_files() {
  local from="$1"
  local to="$2"
  echo "$FILES" | while read -r file; do
    if grep -qF "$from" "$file" 2>/dev/null; then
      sed -i '' "s|${from}|${to}|g" "$file"
    fi
  done
}

replace_in_files "APP_DISPLAY_NAME" "$APP_DISPLAY_NAME"
replace_in_files "APP_SNAKE"        "$APP_SNAKE"
replace_in_files "APP_PASCAL"       "$APP_PASCAL"
replace_in_files "APP_NAME"         "$APP_NAME"
replace_in_files "APP_DOMAIN"       "$APP_DOMAIN"
replace_in_files "GITHUB_USERNAME"  "$GITHUB_USERNAME"

# -- Rename nginx config file --
if [ -f "$ROOT/deploy/nginx/app.conf" ]; then
  mv "$ROOT/deploy/nginx/app.conf" "$ROOT/deploy/nginx/${APP_NAME}.conf"
  # Update the vps compose mount path too
  sed -i '' "s|nginx/app.conf|nginx/${APP_NAME}.conf|g" "$ROOT/deploy/docker-compose.vps.yml"
fi

echo -e "${GREEN}Done!${NC} Project initialised as '${APP_DISPLAY_NAME}'."
echo ""
echo "Next steps:"
echo "  1. cp apps/api/.env.example apps/api/.env   # fill in secrets"
echo "  2. npm install                               # install web dependencies"
echo "  3. docker compose up -d                      # start Postgres"
echo "  4. npm run api                               # start API (auto-runs migrations on boot)"
echo "  5. npm run web                               # start web"
echo "  6. cd apps/mobile && flutter pub get         # install Flutter deps"
echo ""
