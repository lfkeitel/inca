#!/bin/bash

EMAIL_TO=""
EMAIL_FROM=""
MAILER="/usr/sbin/ssmtp"

# These paths must be absolute
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GIT_REPO="$SCRIPT_DIR"
CONFIGS_DIR="/opt/inca/latest"

INCA_SERVER="http://localhost:8080"
# github/gitlab/gitea style project URL
# https://gitlab.com/user/repo
GIT_HTTP_URL=""

[[ -z "$(which jq)" ]] && echo "jq command not found" && exit 1

DRY_RUN="f"
if [[ $1 == "-dr" ]]; then
    DRY_RUN="t"
    shift
fi

INCA_USERNAME="$1"
INCA_PASSWORD="$2"

isDryRun() {
    [[ $DRY_RUN == "t" ]]
}

echoError() {
    >&2 echo -e "$1"
}

mail_report() {
    if [ -z "$(which $MAILER 2> /dev/null)" ]; then
        echoError "Can't send mail, mailer not found"
        return 1
    fi

    if [ -z "$EMAIL_TO" ]; then
        echoError "No email addresses specified, can't send email"
        return 1
    fi

    gitCommitHash="$(git --no-pager log --pretty=oneline -n1 | cut -d' ' -f1)"

    TEMPMAIL="$(mktemp)"
    echo "To: $EMAIL_TO" >> "$TEMPMAIL"
    echo "From: $EMAIL_FROM" >> "$TEMPMAIL"
    echo "Subject: Configurations Changed" >> "$TEMPMAIL"
    echo >> "$TEMPMAIL"
    echo "Commit: $GIT_HTTP_URL/commit/$gitCommitHash" >> "$TEMPMAIL"
    echo >> "$TEMPMAIL"
    echo "The following configurations have been changed:" >> "$TEMPMAIL"
    echo >> "$TEMPMAIL"

    for ((i = 0; i < ${#changedConfigs[@]}; i++)); do
        config="${changedConfigs[$i]}"
        echo "- $config: $GIT_HTTP_URL/src/master/configs/$config/config.conf" >> "$TEMPMAIL"
    done

    $MAILER "$EMAIL_TO" < "$TEMPMAIL"
    rm "$TEMPMAIL"
}

getFileSize() {
    du --apparent-size --block-size=1 "$1" | awk '{print $1}'
}

copyConfigs() {
    cd "$GIT_REPO/configs"

    if [[ "$(ls -l $CONFIGS_DIR | wc -l)" -le 1 ]]; then
        echo "No files"
        return
    fi

    for hostdir in $CONFIGS_DIR/*; do
        if [ ! -d "$hostdir" ]; then
            continue
        fi

        srcconfig="$(ls -1 "$hostdir" | grep -v '_metadata' | sort -r | head -n1)"
        switchName="$(jq -r '.name' "$hostdir/_metadata.json")"
        if [ ! -d "$switchName" ]; then
            mkdir "$switchName"
        fi

        oldSha=""

        if [ -f "$switchName/config.conf" ]; then
            oldSha="$(sha256sum "$switchName/config.conf")"
        fi

        # Strip out unnecessary parts of the file
        case "$(jq -r '.dtype' "$hostdir/_metadata.json")" in
        *juniper*)
            cat "$hostdir/$srcconfig" | \
            sed -e '/^## Last changed/,/\[edit\]/!d' -e '// d' | \
            tr -d '\r' > "$switchName/config.conf"

            grep -P '^set ' "$hostdir/$srcconfig" > "$switchName/display-set.conf"

            # Prevent empty files from overwritting configs
            if [ $(getFileSize "$switchName/display-set.conf") -eq 0 ]; then
                git checkout -- "$switchName/display-set.conf"
            fi
            ;;
        *cisco*)
            cat "$hostdir/$srcconfig" | \
            sed -e '/^version/,/^end/!d;//d' | \
            grep -v 'ntp clock-period' | \
            tr -d '\r' > "$switchName/config.conf"
            ;;
        *)
            cp "$hostdir/$srcconfig" "$switchName/config.conf"
        esac

        # Prevent empty files from overwritting configs
        if [ $(getFileSize "$switchName/config.conf") -eq 0 ]; then
            git checkout -- "$switchName/config.conf"
        fi

        newSha="$(sha256sum "$switchName/config.conf")"
        if [ "$newSha" != "$oldSha" ]; then
            changedConfigs+=("$switchName")
        fi
    done

    cd ..
}

checkGit() {
    git update-index -q --refresh

    untracked="$(git ls-files --other --directory --exclude-standard --no-empty-directory)"
    if [ -n "$untracked" ]; then
        for file in $untracked; do
            git add "$file"
        done
    fi

    if ! git diff-index --quiet HEAD --; then
        if isDryRun; then
            echo "Dry run: Git commit"
        else
            commitAndPush
        fi
    fi
}

commitAndPush() {
    git commit -am "Updated config set"
    git push origin master
}

declare -a changedConfigs

cd "$GIT_REPO"
git pull
copyConfigs
checkGit

if [ ${#changedConfigs[@]} -gt 0 ]; then
    if isDryRun; then
        echo "Dry run: Mail report"
    else
        mail_report
    fi
fi
