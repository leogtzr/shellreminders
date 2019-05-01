#!/bin/bash

# TODO: remove after testing ... 
set -o xtrace
readonly work_dir=$(dirname "$(readlink --canonicalize-existing "${0}")")
grep --fixed-strings '#shellreminers-v.10' ~/.bashrc || {
    echo 'grep --extended-regexp --quiet "root" <<< "$(whoami)" || "/home/leo/bin/shellreminders" #shellreminers-v.10' >> ~/.bashrc
}

exit 0